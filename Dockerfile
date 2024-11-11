# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.22 as builder

WORKDIR /go/app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.mod .
RUN go mod download

# Copy the rest
COPY . .
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todos -mod=readonly -v ./cmd/todos/server.go

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/app/todos ./
COPY ./views ./views
COPY ./public ./public

EXPOSE 3000

# Run the web service on container startup.
ENTRYPOINT ["/todos"]