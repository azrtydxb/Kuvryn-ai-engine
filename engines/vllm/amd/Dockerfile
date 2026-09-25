# syntax=docker/dockerfile:1.7
FROM ubuntu:24.04

ARG ENGINE=vllm
ARG VENDOR=amd
ARG BUILD_DATE
ARG VCS_REF
ENV VIRTUAL_ENV=/opt/kuvryn/venv
ENV PATH="/opt/kuvryn/venv/bin:${PATH}"
ARG VLLM_VERSION=0.30.0
ARG RAY_VERSION=2.58.0
# ROCm 10 ships as Python wheels (TheRock); torch, triton and the SDK come from AMD's index.
ARG ROCM_VERSION=10.0.0
ARG TORCH_VERSION=2.13.0+rocm10.0.0
ARG TORCHVISION_VERSION=0.28.0+rocm10.0.0
ARG TRITON_VERSION=3.8.0+git4cff872c.rocm10.0.0
ARG AITER_REF=v0.1.22.post1
ARG PYTORCH_ROCM_ARCH="gfx90a;gfx942;gfx950;gfx1100;gfx1101;gfx1150;gfx1151;gfx1200;gfx1201"

# install-vllm.sh links the SDK wheel's devel tree to /opt/rocm so torch's cpp_extension,
# vLLM's CMake and AITER's JIT compiler find hipcc and headers where they expect them.
ENV ROCM_PATH=/opt/rocm \
    ROCM_HOME=/opt/rocm \
    HIP_DEVICE_LIB_PATH=/opt/kuvryn/venv/lib/python3.12/site-packages/_rocm_sdk_core/lib/llvm/amdgcn/bitcode \
    LD_LIBRARY_PATH=/opt/rocm/lib:/opt/kuvryn/venv/lib/python3.12/site-packages/_rocm_sdk_core/lib \
    PATH="/opt/rocm/bin:/opt/rocm/llvm/bin:${PATH}" \
    PYTORCH_ROCM_ARCH="${PYTORCH_ROCM_ARCH}" \
    HSA_NO_SCRATCH_RECLAIM=1

LABEL org.opencontainers.image.source="https://github.com/azrtydxb/kuvryn-ai-engine" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      kuvryn.engine="${ENGINE}" \
      kuvryn.vendor="${VENDOR}" \
      kuvryn.flavor="vllm-amd" \
      kuvryn.certification="hardware-untested"

COPY ci/images/kuvryn-vllm-serve ci/images/install-common-debian.sh ci/images/install-python-runtime.sh ci/images/install-vllm.sh /tmp/kuvryn-install/
# hadolint ignore=DL3008
RUN /tmp/kuvryn-install/install-common-debian.sh \
    && /tmp/kuvryn-install/install-python-runtime.sh \
    && KUVRYN_VENDOR="${VENDOR}" VLLM_VERSION="${VLLM_VERSION}" RAY_VERSION="${RAY_VERSION}" \
       ROCM_VERSION="${ROCM_VERSION}" TORCH_VERSION="${TORCH_VERSION}" TORCHVISION_VERSION="${TORCHVISION_VERSION}" \
       TRITON_VERSION="${TRITON_VERSION}" AITER_REF="${AITER_REF}" /tmp/kuvryn-install/install-vllm.sh \
    && install -m 0755 /tmp/kuvryn-install/kuvryn-vllm-serve /usr/local/bin/kuvryn-vllm-serve \
    && rm -rf /tmp/kuvryn-install

COPY bin/kuvryn-runtime-wrapper /usr/local/bin/kuvryn-runtime-wrapper
# debt: image runs as root; revisit once non-root is verified against GPU device access on the NVIDIA, AMD and Intel CT runners.
# nosemgrep: dockerfile.security.missing-user-entrypoint.missing-user-entrypoint
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/kuvryn-runtime-wrapper"]
# nosemgrep: dockerfile.security.missing-user.missing-user
CMD ["kuvryn-vllm-serve"]
