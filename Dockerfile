# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.10.2 as builder

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/c6h3un/echogo

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go install -installsuffix cgo github.com/c6h3un/echogo

# Run the outyet command by default when the container starts.
#ENTRYPOINT /go/bin/echogo

FROM alpine:latest
#RUN apk --no-cache add ca-certificates
RUN apk add --no-cache curl
RUN addgroup -g 1000 -S app && adduser -u 1000 -S -G app app 
USER 1000
WORKDIR /opt/bin/
COPY --from=builder /go/bin/echogo .
CMD ["/opt/bin/echogo"]
# Document that the service listens on port 8888.
EXPOSE 8888
