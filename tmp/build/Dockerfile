FROM alpine:3.6

ADD tmp/_output/bin/ambassador-shim-operator-sdk /usr/local/bin/ambassador-shim-operator-sdk

RUN adduser -D ambassador-shim-operator-sdk
USER ambassador-shim-operator-sdk
