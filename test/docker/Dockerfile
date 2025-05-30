ARG GO_IMAGE=golang:1.23.9
FROM $GO_IMAGE

# Copy SDK code into the container
RUN mkdir -p /app/go-algorand-sdk
COPY . /app/go-algorand-sdk
WORKDIR /app/go-algorand-sdk

# Run integration tests
CMD ["/bin/bash", "-c", "make unit && make integration && make smoke-test-examples"]
