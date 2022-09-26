############# builder
FROM golang:1.18.5 AS builder

WORKDIR /go/src/github.com/Gerrit91/gardener-extension-registry-cache
COPY . .
RUN make install

############# gardener-extension-registry-cache
FROM gcr.io/distroless/static-debian11:nonroot AS gardener-extension-registry-cache
WORKDIR /

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-registry-cache /gardener-extension-registry-cache
ENTRYPOINT ["/gardener-extension-registry-cache"]
