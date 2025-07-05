############################
# STEP 1 build executable binary
############################
FROM golang:1.24-alpine AS builder
# Set the GOPATH inside the container
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
WORKDIR /builder
# Copy Go module files and dowuload dependencies
COPY go.mod go.sum ./
RUN go mod download
# COPY the rest of the source code
COPY .. .
# Build the binary
RUN mkdir -p bin
RUN CGO_ENABLED=0 go build -o bin/service cmd/server/*.go
############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /builder/bin/service /usr/local/bin/appbin
COPY --from=builder /builder/database/migrations /usr/local/database/migrations

ENTRYPOINT ["/usr/local/bin/appbin"]