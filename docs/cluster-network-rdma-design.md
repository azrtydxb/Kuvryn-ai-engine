# Cluster network and RDMA discovery design

## Status

Kuvryn AI Engine does not yet have implemented cluster network discovery or RDMA setup code in this repository. The current repository has design documents only.

This document makes network/RDMA discovery a required part of the engine runtime foundation, because distributed inference must know whether nodes can actually communicate over the expected fabric before selecting sharding, tensor parallel, pipeline parallel, or collective settings.

## Goal

Provide Sparkrun-class cluster discovery, but generalized for NVIDIA, AMD, Intel, Kubernetes, and direct Docker/local execution.

The cluster builder must discover:

- management connectivity;
- accelerator inventory;
- high-speed fabric interfaces;
- RDMA/RoCE/InfiniBand capabilities;
- topology shape;
- per-host collective environment;
- whether distributed inference is safe to schedule.

## Non-goals

- Do not assume every cluster is DGX Spark or ConnectX-7.
- Do not assume RDMA is always present.
- Do not silently fall back from RDMA to slow Ethernet for a deployment that requires RDMA.
- Do not encode host-specific HCA/interface names in recipes.
- Do not mutate host networking during read-only discovery.

## Sparkrun lessons to carry forward

Sparkrun has several useful mechanisms we should keep conceptually:

| Sparkrun idea                                         | Kuvryn version                                                                                      |
| ----------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| Combined accelerator + IB probe in one SSH round trip | Combined host capability probe for accelerator, network, driver, fabric, and runtime prerequisites. |
| Sentinel-delimited probe output                       | Strict machine-readable probe sections with versioned schema.                                       |
| CX7/RoCE detection                                    | Vendor-neutral fabric detection with specialized ConnectX handling.                                 |
| Topology classification                               | Classify direct, switched, ring/mesh, Kubernetes fabric, and unknown.                               |
| NCCL env generation                                   | Collective env generation for NCCL, RCCL, HCCL, oneCCL, and UCX.                                    |
| RDMA perftest                                         | Real fabric validation, not just IP reachability.                                                   |
| Ring-specific overrides                               | Topology-specific collective tuning with explicit reasons.                                          |

## Architecture

```text
cluster builder
  -> host access resolver
  -> host probe runner
  -> network/fabric parser
  -> topology classifier
  -> RDMA validator
  -> collective env planner
  -> placement constraints
  -> runtime launch plan
```

## Host capability probe

Each host gets one combined probe. For SSH/direct hosts this is a remote script. For Kubernetes this is a privileged or device-aware diagnostic Job/DaemonSet.

Probe sections:

```text
KUVRYN_PROBE_ACCEL_START
...
KUVRYN_PROBE_ACCEL_END
KUVRYN_PROBE_NETWORK_START
...
KUVRYN_PROBE_NETWORK_END
KUVRYN_PROBE_FABRIC_START
...
KUVRYN_PROBE_FABRIC_END
KUVRYN_PROBE_RUNTIME_START
...
KUVRYN_PROBE_RUNTIME_END
```

The parser must reject malformed or missing section boundaries. Partial results are allowed only with explicit limitations attached to the host record.

## Data model

```yaml
host: node-a
management:
  reachable: true
  iface: eno1
  ip: 10.0.0.11
accelerators:
  - vendor: nvidia
    model: H100
    count: 8
    memory_gb: 80
    capabilities: [cuda, nvlink]
fabric:
  rdma_present: true
  interfaces:
    - name: rocep1s0f1
      ip: 192.168.10.11
      subnet: 192.168.10.0/24
      mtu: 9000
      hca: mlx5_0
      link_layer: ethernet
      gid_index: 3
      roce_version: v2
      speed_gbps: 200
  capabilities: [rdma, roce-v2]
collectives:
  supported: [nccl]
  preferred: nccl
limitations: []
```

## Network discovery

Discovery must collect, per host:

- hostname and stable node identity;
- management IP/interface;
- default route interface;
- interface names, MACs, MTU, state, addresses, subnet;
- RDMA HCA names and port state;
- GID table and RoCE v2 IPv4 index where applicable;
- InfiniBand/RoCE link layer;
- driver/tool availability: `nvidia-smi`, `rocm-smi`, `hl-smi`, `ibv_devinfo`, `show_gids`, `rdma`, `ethtool`;
- Kubernetes device-plugin resources when running in Kubernetes.

Discovery is read-only by default.

## Topology classification

The classifier should produce one of:

| Topology            | Meaning                                                      | Scheduling effect                                                    |
| ------------------- | ------------------------------------------------------------ | -------------------------------------------------------------------- |
| `single-node`       | One host only.                                               | No inter-node collective required.                                   |
| `direct`            | Two hosts with direct high-speed links.                      | Distributed allowed if RDMA validation passes.                       |
| `ring`              | Three or more hosts connected in ring/mesh style.            | Requires topology-aware collective env.                              |
| `switch`            | Hosts share a switched RDMA fabric.                          | Preferred for multi-node distributed inference.                      |
| `kubernetes-fabric` | Fabric exposed through Kubernetes networking/device plugins. | Use Kubernetes-specific validation path.                             |
| `ethernet-only`     | No RDMA fabric found.                                        | Disallow RDMA-required plans; allow non-RDMA only if policy permits. |
| `unknown`           | Detection incomplete.                                        | Fail closed for distributed/RDMA-required workloads.                 |

## RDMA validation

IP reachability is not enough. Validation must distinguish a real high-speed RDMA path from a slow or misrouted interface.

Validation stages:

1. SSH/API reachability on management path.
2. RDMA device presence and active port state.
3. Fabric IP reachability between planned peers.
4. Real RDMA transfer test when tools are available.
5. Collective bootstrap test for distributed flavors.

Suggested tools:

- `ib_write_bw` / `ib_write_lat` for RDMA bandwidth/latency;
- `nccl-tests` for NVIDIA/NCCL;
- RCCL equivalent tests for AMD;
- HCCL/oneCCL tests for Intel where available;
- UCX checks for engines using UCX.

A deployment requiring RDMA must not proceed when validation is missing, failed, or stale.

## Collective environment planning

The planner derives environment variables from live discovery, not from recipes.

| Vendor | Collective        | Example env family                                                              |
| ------ | ----------------- | ------------------------------------------------------------------------------- |
| NVIDIA | NCCL              | `NCCL_IB_HCA`, `NCCL_SOCKET_IFNAME`, `NCCL_IB_GID_INDEX`, `UCX_NET_DEVICES`     |
| AMD    | RCCL              | `RCCL_*`, `NCCL_*` compatibility vars where ROCm honors them, `UCX_NET_DEVICES` |
| Intel  | HCCL/oneCCL       | HCCL/CCL fabric variables and socket interface pins                             |
| CPU    | none/UCX optional | Socket interface pins only when distributed CPU runtime exists                  |

Recipe env may override generated env only with an explicit unsafe/advanced flag. Otherwise, host-specific HCA and interface names are owned by the cluster planner.

## Kubernetes mode

Kubernetes discovery must support:

- node labels and allocatable GPU resources;
- device plugin resources, for example `nvidia.com/gpu`, AMD GPU resource keys, Intel/Gaudi resources;
- RDMA device plugin resources where installed;
- Multus or secondary network attachment definitions when used;
- pod-level visibility of `/dev/infiniband` or equivalent devices;
- a diagnostic DaemonSet/Job to validate fabric access from the same security context as runtime pods.

Kubernetes scheduling must not assume host-level RDMA visibility equals pod-level RDMA visibility.

## Cluster build phases

```text
1. register hosts or import Kubernetes cluster
2. read-only probe all candidate nodes
3. classify topology
4. validate RDMA/fabric where required
5. generate placement constraints
6. generate collective env plan
7. select engine image flavor
8. launch runtime
9. run readiness and inference smoke
```

Host mutation, such as writing netplan or installing packages, is a separate explicit setup operation. It must never happen during read-only discovery.

## Failure behavior

| Failure                                                | Behavior                                                                         |
| ------------------------------------------------------ | -------------------------------------------------------------------------------- |
| Missing RDMA on distributed/RDMA-required runtime      | Fail placement with actionable diagnostics.                                      |
| Mixed vendor nodes without explicit layout             | Fail placement.                                                                  |
| RDMA device present but transfer test fails            | Mark fabric unhealthy and block publish/deploy paths requiring it.               |
| Management network works but fabric IP does not        | Do not use RDMA path; fail RDMA-required workloads.                              |
| Kubernetes node has GPU but pod cannot see RDMA device | Fail Kubernetes RDMA validation.                                                 |
| Unknown topology                                       | Allow single-node only unless user explicitly accepts degraded distributed mode. |

## Acceptance criteria for implementation

- Unit tests cover parser behavior for complete, partial, and malformed probe output.
- Topology classifier tests cover single-node, direct, ring, switch, ethernet-only, and unknown.
- Env planner tests cover NCCL, RCCL, HCCL/oneCCL, and no-RDMA cases.
- Kubernetes validation tests distinguish node-visible from pod-visible RDMA.
- Distributed runtime planning refuses RDMA-required workloads when validation is absent or failed.
- No discovery path prints secrets, tokens, or full environment dumps.
