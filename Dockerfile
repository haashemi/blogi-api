FROM golang:1.23.0-alpine3.20

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./internal ./internal

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /blogi ./cmd/blogi

EXPOSE 8080

# Run
CMD ["/blogi"]
