# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8.3

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/c6h3un/echogo

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/c6h3un/echogo

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/echogo

# Document that the service listens on port 8888.
EXPOSE 8888
