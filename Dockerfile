FROM golang:1.24-alpine AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/kubebridge ./cmd/kubebridge

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /out/kubebridge /kubebridge

USER nonroot:nonroot

ENTRYPOINT ["/kubebridge"]
