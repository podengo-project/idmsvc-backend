##
# This build the container image for the service
##

# https://catalog.redhat.com/software/containers/ubi9/go-toolset/61e5c00b4ec9945c18787690
FROM registry.access.redhat.com/ubi9/go-toolset:1.21.11-9.1726696856@sha256:2fead14be2734020f98c731fef30cc54e99fe889a53cac97e80f6f043778b567 as builder
LABEL idmsvc-backend=builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make get-deps build


# https://catalog.redhat.com/software/containers/ubi9-minimal/61832888c0d15aff4912fe0d
FROM registry.access.redhat.com/ubi9-minimal:9.4-1227.1726694542@sha256:c0e70387664f30cd9cf2795b547e4a9a51002c44a4a86aa9335ab030134bf392
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
