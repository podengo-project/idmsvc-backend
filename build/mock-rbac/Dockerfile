##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.23.9-1749052980@sha256:d4bcb27c15c76c727742a78173ba02a56d9acc695ecc2603a0dc0fc745a75685 as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make bin/mock-rbac

# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.6-1749489516@sha256:f172b3082a3d1bbe789a1057f03883c1113243564f01cd3020e27548b911d3f8
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
