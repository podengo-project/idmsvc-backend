##
# This build the container image for the service
##
FROM registry.access.redhat.com/ubi9/go-toolset:1.19 as builder
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
WORKDIR /go/src/app
COPY . .
USER 0
RUN make get-deps build


FROM registry.access.redhat.com/ubi9/ubi:9.2-755
# https://developers.redhat.com/articles/2022/05/31/your-go-application-fips-compliant
ENV OPENSSL_FORCE_FIPS_MODE=1
RUN mkdir -p /opt/bin /opt/bin/scripts/db /opt/bin/configs
WORKDIR /opt/bin
COPY --from=builder /go/src/app/bin/* ./
COPY scripts/db/migrations /opt/bin/scripts/db/migrations
COPY configs/config.example.yaml /opt/bin/configs/config.yaml
USER 1001
# TODO Accomodate to the service arguments
CMD ["/opt/bin/service"]
