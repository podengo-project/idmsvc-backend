##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.21.11-2.1720624888@sha256:5c948cdfd0132e982426bc9d3a81eeae66871080ef274abdde1a4a8303509188 as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN git log -1
RUN make get-deps build


# https://catalog.redhat.com/software/containers/ubi9-minimal/61832888c0d15aff4912fe0d
FROM registry.access.redhat.com/ubi9-minimal@sha256:a7d837b00520a32502ada85ae339e33510cdfdbc8d2ddf460cc838e12ec5fa5a
LABEL idmsvc-backend=backend
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
RUN mkdir -p /opt/bin /opt/bin/scripts/db /opt/bin/configs
WORKDIR /opt/bin
RUN microdnf update -y && microdnf clean all
COPY --from=builder /go/src/app/bin/* ./
COPY scripts/db/migrations /opt/bin/scripts/db/migrations
COPY configs/config.example.yaml /opt/bin/configs/config.yaml
USER 1001
# TODO Accomodate to the service arguments
CMD ["/opt/bin/service"]
