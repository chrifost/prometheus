# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go
# Use 'onbuild' variant - this automatically copies source, buids and configures
# for startup
FROM golang

# Copy the local package files to the container's workspace
ADD . /go/src/p8s-test

# Get prometheus repo
RUN go get github.com/prometheus/client_golang/prometheus

# Build install our p8s test app
RUN go install p8s-test

# Document that the service listens on port 8080
EXPOSE 8080

# Run the p8s-test command by default when the container starts.
ENTRYPOINT /go/bin/p8s-test

# Start this continer with: go run -t -p 8080:8080 my-p8s-app
# Incrememet counter with: localhost:8080/counter
# Prometheus metrics at:   localhost:8080/metrics
