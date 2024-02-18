FROM golang:latest as builder

ADD . /usr/src/price-tracker
WORKDIR /usr/src/price-tracker
RUN --mount=type=cache,target=/go/pkg/mod \
      --mount=type=bind,source=go.mod,target=go.mod \
      --mount=type=bind,source=go.sum,target=go.sum \
      go mod download -x
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /usr/bin/price-tracker .

FROM debian:stable-slim
WORKDIR /opt/price-tracker
RUN apt-get update && apt-get install ca-certificates curl -y
COPY --from=builder /usr/bin/price-tracker /usr/bin/price-tracker

CMD ["/usr/bin/price-tracker"]
