# syntax=docker/dockerfile:1.7
FROM ubuntu:24.04

ARG ENGINE=llama-cpp
ARG VENDOR=cpu
ARG BUILD_DATE
ARG VCS_REF
ARG LLAMA_CPP_REF=v0.4.1
ARG LLAMA_CPP_CMAKE_FLAGS=""

LABEL org.opencontainers.image.source="https://github.com/azrtydxb/kuvryn-ai-engine" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      kuvryn.engine="${ENGINE}" \
      kuvryn.vendor="${VENDOR}" \
      kuvryn.flavor="llama-cpp-cpu" \
      kuvryn.certification="cpu-only"

COPY ci/images/kuvryn-llama-cpp-serve ci/images/install-common-debian.sh ci/images/install-llama-cpp.sh /tmp/kuvryn-install/
# hadolint ignore=DL3008
RUN /tmp/kuvryn-install/install-common-debian.sh \
    && LLAMA_CPP_REF="${LLAMA_CPP_REF}" LLAMA_CPP_CMAKE_FLAGS="${LLAMA_CPP_CMAKE_FLAGS}" /tmp/kuvryn-install/install-llama-cpp.sh \
    && install -m 0755 /tmp/kuvryn-install/kuvryn-llama-cpp-serve /usr/local/bin/kuvryn-llama-cpp-serve \
    && rm -rf /tmp/kuvryn-install

COPY bin/kuvryn-runtime-wrapper /usr/local/bin/kuvryn-runtime-wrapper
# debt: image runs as root; revisit once non-root is verified against GPU device access on the NVIDIA, AMD and Intel CT runners.
# nosemgrep: dockerfile.security.missing-user-entrypoint.missing-user-entrypoint
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/kuvryn-runtime-wrapper"]
# nosemgrep: dockerfile.security.missing-user.missing-user
CMD ["kuvryn-llama-cpp-serve"]
