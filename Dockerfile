# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS builder
ENV GOTOOLCHAIN=local
WORKDIR /src
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/convertix ./cmd/convertix


FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/convertix /convertix
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/convertix"]