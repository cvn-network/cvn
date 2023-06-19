FROM golang:1.20-bullseye as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/cvn-network/cvn

# Add source files
COPY . .

# build cvnd
RUN make build-linux

# Final image
FROM debian:bullseye as final

WORKDIR /

RUN apt-get update

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/cvn-network/cvn/build/cvnd /
COPY --from=build-env /go/src/github.com/cvn-network/cvn/networks/local/start-docker.sh /

EXPOSE 26656 26657 1317 9090 8545 8546

# Run cvnd by default, omit entrypoint to ease using container with cvnd
ENTRYPOINT ["/bin/bash", "-c"]