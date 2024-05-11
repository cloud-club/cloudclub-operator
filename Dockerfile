FROM golang:1.21.3-alpine as base
# ARG GOPATH
RUN apk --no-cache add tzdata ca-certificates
WORKDIR /src/
RUN --mount=type=bind,source=go.mod,target=go.mod \
        --mount=type=bind,source=go.sum,target=go.sum \
        --mount=type=cache,target=/go/pkg/mod/ \
        go mod download -x

FROM base as stage
WORKDIR /src/
RUN --mount=type=cache,target=/go/pkg/mod/ \
        --mount=type=bind,source=.,target=. \
        CGO_ENABLED=0  go build -a -o /bin/manager main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=stage /bin/manager /manager
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Asia/Seoul
ENTRYPOINT ["/manager"]
