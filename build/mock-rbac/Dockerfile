##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.24.4-1753853351@sha256:3ce6311380d5180599a3016031a9112542d43715244816d1d0eabc937952667b as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make bin/mock-rbac

# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.6-1758184547@sha256:7c5495d5fad59aaee12abc3cbbd2b283818ee1e814b00dbc7f25bf2d14fa4f0c
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
