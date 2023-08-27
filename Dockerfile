# Build
FROM golang:1.21-alpine3.17 AS build
WORKDIR /build

# Install dependencies
COPY go.* .
RUN go mod download

# Build the binary
ENV CGO_ENABLED=0
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app ./cmd/app/main.go

## Deploy
FROM scratch

# Copy our static executable
COPY --from=build /app /app
COPY ./swagger ./swagger

# Create environment
COPY env/app.env /

# Expose application port
EXPOSE ${APP_PORT}

# Run the binary
ENTRYPOINT ["/app", "--config=/app.env"]
