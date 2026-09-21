# syntax=docker/dockerfile:1.7
FROM intel/oneapi-runtime:2025.3.1-0-devel-ubuntu24.04

ARG ENGINE=sglang
ARG VENDOR=intel
ARG BUILD_DATE
ARG VCS_REF
ENV VIRTUAL_ENV=/opt/kuvryn/venv
ENV PATH="/opt/kuvryn/venv/bin:${PATH}"
ARG SGLANG_VERSION=0.5.19

LABEL org.opencontainers.image.source="https://github.com/azrtydxb/kuvryn-ai-engine" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      kuvryn.engine="${ENGINE}" \
      kuvryn.vendor="${VENDOR}" \
      kuvryn.flavor="sglang-intel" \
      kuvryn.certification="hardware-untested"

COPY ci/images/kuvryn-sglang-serve ci/images/install-common-debian.sh ci/images/install-python-runtime.sh ci/images/install-sglang.sh /tmp/kuvryn-install/
# hadolint ignore=DL3008
RUN /tmp/kuvryn-install/install-common-debian.sh \
    && /tmp/kuvryn-install/install-python-runtime.sh \
    && KUVRYN_VENDOR="${VENDOR}" SGLANG_VERSION="${SGLANG_VERSION}" /tmp/kuvryn-install/install-sglang.sh \
    && install -m 0755 /tmp/kuvryn-install/kuvryn-sglang-serve /usr/local/bin/kuvryn-sglang-serve \
    && rm -rf /tmp/kuvryn-install

COPY bin/kuvryn-runtime-wrapper /usr/local/bin/kuvryn-runtime-wrapper
# debt: image runs as root; revisit once non-root is verified against GPU device access on the NVIDIA, AMD and Intel CT runners.
# nosemgrep: dockerfile.security.missing-user-entrypoint.missing-user-entrypoint
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/kuvryn-runtime-wrapper"]
# nosemgrep: dockerfile.security.missing-user.missing-user
CMD ["kuvryn-sglang-serve"]
