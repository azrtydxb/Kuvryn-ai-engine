# syntax=docker/dockerfile:1.7
FROM rocm/dev-ubuntu-24.04:10.0.0-full

ARG ENGINE=llama-cpp
ARG VENDOR=amd
ARG BUILD_DATE
ARG VCS_REF
ARG LLAMA_CPP_REF=v0.5.0
ARG LLAMA_CPP_GPU_TARGETS="gfx908;gfx90a;gfx942;gfx950;gfx1030;gfx1100;gfx1101;gfx1102;gfx1150;gfx1151;gfx1200;gfx1201"
ARG LLAMA_CPP_CMAKE_FLAGS="-DGGML_HIP=ON -DGGML_NATIVE=OFF"
# The ROCm 10 dev image does not register /opt/rocm/lib with the dynamic linker.
ENV LD_LIBRARY_PATH=/opt/rocm/lib

LABEL org.opencontainers.image.source="https://github.com/azrtydxb/kuvryn-ai-engine" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      kuvryn.engine="${ENGINE}" \
      kuvryn.vendor="${VENDOR}" \
      kuvryn.flavor="llama-cpp-amd" \
      kuvryn.certification="hardware-untested"

COPY ci/images/kuvryn-llama-cpp-serve ci/images/install-common-debian.sh ci/images/install-llama-cpp.sh /tmp/kuvryn-install/
# hadolint ignore=DL3008
RUN /tmp/kuvryn-install/install-common-debian.sh \
    && HIPCXX="$(hipconfig -l)/clang" HIP_PATH="$(hipconfig -R)" LLAMA_CPP_REF="${LLAMA_CPP_REF}" LLAMA_CPP_CMAKE_FLAGS="${LLAMA_CPP_CMAKE_FLAGS} -DAMDGPU_TARGETS=${LLAMA_CPP_GPU_TARGETS}" /tmp/kuvryn-install/install-llama-cpp.sh \
    && install -m 0755 /tmp/kuvryn-install/kuvryn-llama-cpp-serve /usr/local/bin/kuvryn-llama-cpp-serve \
    && rm -rf /tmp/kuvryn-install

COPY bin/kuvryn-runtime-wrapper /usr/local/bin/kuvryn-runtime-wrapper
# debt: image runs as root; revisit once non-root is verified against GPU device access on the NVIDIA, AMD and Intel CT runners.
# nosemgrep: dockerfile.security.missing-user-entrypoint.missing-user-entrypoint
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/kuvryn-runtime-wrapper"]
# nosemgrep: dockerfile.security.missing-user.missing-user
CMD ["kuvryn-llama-cpp-serve"]
