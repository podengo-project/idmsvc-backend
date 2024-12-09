##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:9.5-1733160835@sha256:c9c0129bff8e94b9e23e9e4a1ebb5b932cccc24e064940a792bdc002da512d2a as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make bin/mock-rbac

# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.5-1733767867@sha256:f598528219a1be07cf520fbe82a2d2434dc9841e1f0a878382c8a13bf42cb486
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
