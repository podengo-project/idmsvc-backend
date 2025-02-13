##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.22.9-1739184363@sha256:7d918e22a74a8151f697775cf4f4e3ede3766c32b79a4ac0fded15a8c3f76466 as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make get-deps build
RUN mkdir /licenses \
    && cp -vf LICENSE /licenses/LICENSE


# https://catalog.redhat.com/software/containers/ubi9/ubi-minimal/615bd9b4075b022acc111bf5
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.5-1739377040@sha256:98d90ac833df94ce6c07c7908ff29cfa3398a71b629e22d7af2d0df050ca7757
LABEL idmsvc-backend=backend
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
RUN mkdir -p /opt/bin /opt/bin/scripts/db /opt/bin/configs
WORKDIR /opt/bin
RUN microdnf update -y && microdnf clean all
COPY --from=builder /go/src/app/bin/* ./
COPY --from=builder /licenses /licenses
COPY scripts/db/migrations /opt/bin/scripts/db/migrations
COPY configs/config.example.yaml /opt/bin/configs/config.yaml
USER 1001
# TODO Accomodate to the service arguments
CMD ["/opt/bin/service"]
