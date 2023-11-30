# syntax=docker/dockerfile:1

FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /tlsproxy

# Expose (accept client connections)
EXPOSE 443
EXPOSE 9200

ENV TLSPROXY_TARGET "localhost"
ENV TLSPROXY_PORTS "443:80"
ENV TLSPROXY_CERT_FILE "server.crt"
ENV TLSPROXY_KEY_FILE "server.key"

# Run
CMD ["/tlsproxy"]