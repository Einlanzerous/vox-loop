FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /vox-loop ./cmd/vox-loop

FROM matrixdotorg/dendrite-monolith:latest

COPY --from=builder /vox-loop /usr/local/bin/vox-loop

ENTRYPOINT ["/usr/local/bin/vox-loop"]
