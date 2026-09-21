#!/usr/bin/env bash
set -euo pipefail

flavor=${1:?usage: engine-runtime-smoke-job.sh <flavor> <image>}
image=${2:?usage: engine-runtime-smoke-job.sh <flavor> <image>}
namespace=${KUVRYN_K8S_NAMESPACE:-default}
context=${KUVRYN_KUBECTL_CONTEXT:-}
job=${KUVRYN_K8S_JOB_NAME:-kuvryn-smoke-${flavor}}
timeout=${KUVRYN_K8S_SMOKE_TIMEOUT:-15m}
node_selector=${KUVRYN_K8S_NODE_SELECTOR:-}

kube() {
	if [[ -n "$context" ]]; then
		kubectl --context "$context" "$@"
	else
		kubectl "$@"
	fi
}

case "$flavor" in
llama-cpp-*)
	repo=${KUVRYN_LLAMA_GGUF_REPO:-aladar/tiny-random-LlamaForCausalLM-GGUF}
	file=${KUVRYN_LLAMA_GGUF_FILE:-tiny-random-LlamaForCausalLM.gguf}
	runtime_script=$(
		cat <<EOF
set -eux
mkdir -p /models
curl -fL --retry 3 --retry-delay 2 -o /models/${file} https://huggingface.co/${repo}/resolve/main/${file}
llama-server -m /models/${file} --host 127.0.0.1 --port 8000 > /tmp/kuvryn-engine.log 2>&1 &
pid=\$!
trap 'kill \$pid 2>/dev/null || true' EXIT
for i in \$(seq 1 90); do
  if curl -fsS http://127.0.0.1:8000/v1/models >/tmp/models.json; then break; fi
  if ! kill -0 \$pid 2>/dev/null; then cat /tmp/kuvryn-engine.log; exit 1; fi
  sleep 2
done
cat /tmp/models.json
curl -fsS -H 'Content-Type: application/json' -d '{"model":"${file}","messages":[{"role":"user","content":"Say ok"}],"max_tokens":8}' http://127.0.0.1:8000/v1/chat/completions | tee /tmp/chat.json
grep -q 'choices' /tmp/chat.json
EOF
	)
	;;
vllm-*)
	model=${KUVRYN_HF_TINY_MODEL:-HuggingFaceH4/tiny-random-LlamaForCausalLM}
	runtime_script=$(
		cat <<EOF
set -eux
vllm serve ${model} --host 127.0.0.1 --port 8000 > /tmp/kuvryn-engine.log 2>&1 &
pid=\$!
trap 'kill \$pid 2>/dev/null || true' EXIT
for i in \$(seq 1 180); do
  if curl -fsS http://127.0.0.1:8000/v1/models >/tmp/models.json; then break; fi
  if ! kill -0 \$pid 2>/dev/null; then cat /tmp/kuvryn-engine.log; exit 1; fi
  sleep 2
done
cat /tmp/models.json
curl -fsS -H 'Content-Type: application/json' -d '{"model":"${model}","messages":[{"role":"user","content":"Say ok"}],"max_tokens":8}' http://127.0.0.1:8000/v1/chat/completions | tee /tmp/chat.json
grep -q 'choices' /tmp/chat.json
EOF
	)
	;;
sglang-*)
	model=${KUVRYN_HF_TINY_MODEL:-HuggingFaceH4/tiny-random-LlamaForCausalLM}
	runtime_script=$(
		cat <<EOF
set -eux
sglang --model-path ${model} --host 127.0.0.1 --port 8000 > /tmp/kuvryn-engine.log 2>&1 &
pid=\$!
trap 'kill \$pid 2>/dev/null || true' EXIT
for i in \$(seq 1 180); do
  if curl -fsS http://127.0.0.1:8000/v1/models >/tmp/models.json; then break; fi
  if ! kill -0 \$pid 2>/dev/null; then cat /tmp/kuvryn-engine.log; exit 1; fi
  sleep 2
done
cat /tmp/models.json
curl -fsS -H 'Content-Type: application/json' -d '{"model":"${model}","messages":[{"role":"user","content":"Say ok"}],"max_tokens":8}' http://127.0.0.1:8000/v1/chat/completions | tee /tmp/chat.json
grep -q 'choices' /tmp/chat.json
EOF
	)
	;;
tensorrt-llm-*)
	if [[ -z "${KUVRYN_TRTLLM_SERVE_ARGS:-}" ]]; then
		echo "KUVRYN_TRTLLM_SERVE_ARGS is required for TensorRT-LLM Kubernetes smoke" >&2
		exit 2
	fi
	runtime_script=$(
		cat <<EOF
set -eux
trtllm-serve ${KUVRYN_TRTLLM_SERVE_ARGS} > /tmp/kuvryn-engine.log 2>&1 &
pid=\$!
trap 'kill \$pid 2>/dev/null || true' EXIT
for i in \$(seq 1 180); do
  if curl -fsS http://127.0.0.1:8000/v1/models >/tmp/models.json; then break; fi
  if ! kill -0 \$pid 2>/dev/null; then cat /tmp/kuvryn-engine.log; exit 1; fi
  sleep 2
done
cat /tmp/models.json
EOF
	)
	;;
*)
	echo "unknown flavor: $flavor" >&2
	exit 2
	;;
esac

selector_yaml=""
if [[ -n "$node_selector" ]]; then
	selector_yaml="      nodeSelector:"
	IFS=',' read -ra pairs <<<"$node_selector"
	for pair in "${pairs[@]}"; do
		key=${pair%%=*}
		value=${pair#*=}
		selector_yaml+=$'\n'"        ${key}: \"${value}\""
	done
fi

kube -n "$namespace" delete job "$job" --ignore-not-found --wait=true >/dev/null 2>&1 || true
kube -n "$namespace" apply -f - >/dev/null <<YAML
apiVersion: batch/v1
kind: Job
metadata:
  name: ${job}
spec:
  backoffLimit: 0
  template:
    spec:
      restartPolicy: Never
      tolerations:
        - operator: Exists
${selector_yaml}
      containers:
        - name: smoke
          image: ${image}
          imagePullPolicy: Always
          command: ["/bin/sh", "-lc"]
          args:
            - |
$(printf '%s\n' "$runtime_script" | sed 's/^/              /')
YAML

if ! kube -n "$namespace" wait --for=condition=complete "job/$job" --timeout="$timeout"; then
	kube -n "$namespace" logs "job/$job" --all-containers=true --tail=-1 >&2 || true
	exit 1
fi
kube -n "$namespace" logs "job/$job" --all-containers=true --tail=-1
