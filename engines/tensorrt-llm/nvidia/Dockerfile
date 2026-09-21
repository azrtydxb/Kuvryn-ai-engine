# syntax=docker/dockerfile:1.7
FROM nvcr.io/nvidia/tensorrt-llm/release:1.2.1

ARG ENGINE=tensorrt-llm
ARG VENDOR=nvidia
ARG BUILD_DATE
ARG VCS_REF
ENV LD_LIBRARY_PATH="/usr/local/tensorrt/lib:${LD_LIBRARY_PATH}"

LABEL org.opencontainers.image.source="https://github.com/azrtydxb/kuvryn-ai-engine" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      kuvryn.engine="${ENGINE}" \
      kuvryn.vendor="${VENDOR}" \
      kuvryn.flavor="tensorrt-llm-nvidia" \
      kuvryn.certification="hardware-certified"

COPY ci/images/kuvryn-trtllm-serve ci/images/install-common-debian.sh ci/images/verify-tensorrt-llm.sh /tmp/kuvryn-install/
# hadolint ignore=DL3008
RUN /tmp/kuvryn-install/install-common-debian.sh \
    && TENSORRT_LLM_SKIP_HELP=1 /tmp/kuvryn-install/verify-tensorrt-llm.sh \
    && install -m 0755 /tmp/kuvryn-install/kuvryn-trtllm-serve /usr/local/bin/kuvryn-trtllm-serve \
    && rm -rf /tmp/kuvryn-install

COPY bin/kuvryn-runtime-wrapper /usr/local/bin/kuvryn-runtime-wrapper
# debt: image runs as root; revisit once non-root is verified against GPU device access on the NVIDIA, AMD and Intel CT runners.
# nosemgrep: dockerfile.security.missing-user-entrypoint.missing-user-entrypoint
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/kuvryn-runtime-wrapper"]
# nosemgrep: dockerfile.security.missing-user.missing-user
CMD ["kuvryn-trtllm-serve"]
