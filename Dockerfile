FROM golang:1-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Now copy the source code
COPY cmd ./cmd
COPY internal ./internal

# Build the binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o server cmd/main.go
# RUN CGO_ENABLED=0 go build -o server cmd/main.go 

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
    
COPY --from=build /app/server .

COPY --from=build /app/internal ./internal 

ENV PORT 8080

CMD ["./server"]