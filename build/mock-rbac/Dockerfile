##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.21.9-1.1715774364@sha256:cc5b93478d74c7a67d3add8be3b41a900f6c97b0cece878465d6aa84edb0e693 as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN git log -1
RUN make bin/mock-rbac

# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4-949.1714662671@sha256:cb6f20b9225031c8d704af95932c0331cfe0b229cd5b3dd805681d92a439186e
LABEL idmsvc-backend=backend
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
RUN mkdir -p /opt/bin /opt/bin/scripts/db /opt/bin/configs
WORKDIR /opt/bin
COPY --from=builder /go/src/app/bin/mock-rbac ./
USER 1001

ENV CLIENTS_RBAC_BASE_URL=http://0.0.0.0:8020/api/rbac/v1
# Valid values for APP_CLIENTS_RBAC_PROFILE are
# super-admin domain-admin domain-read-only domain-no-perms custom
ENV APP_CLIENTS_RBAC_PROFILE=domain-admin

# Command to execute by default
CMD ["/opt/bin/mock-rbac"]
