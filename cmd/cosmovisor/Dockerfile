FROM golang:1.20.2-bullseye AS build-env

RUN apt-get update -y
RUN apt-get install git -y

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.4.0

WORKDIR /root

COPY . .

RUN git checkout release/v1.0.x
RUN make build
RUN mv /root/build/cvnd /root/build/cvnd-v1

RUN git checkout release/v2.0.x
RUN make build

FROM debian:bullseye

RUN apt-get update -y
RUN apt-get install ca-certificates jq -y

WORKDIR /root

ENV DAEMON_NAME=cvnd
ENV DAEMON_HOME=/root/.cvnd
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=false
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_POLL_INTERVAL=1s
ENV DAEMON_RESTART_DELAY=1s
ENV UNSAFE_SKIP_BACKUP=true
ENV DAEMON_DATA_BACKUP_DIR=/root/.cvnd
ENV DAEMON_PREUPGRADE_MAX_RETRIES=3

COPY --from=build-env /go/bin/cosmovisor /usr/bin/cosmovisor

COPY --from=build-env /root/build/cvnd-v1 /usr/bin/cvnd

RUN cosmovisor init /usr/bin/cvnd

COPY --from=build-env /root/build/cvnd /usr/bin/cvnd

RUN mkdir -p /root/.cvnd/cosmovisor/upgrades/v2.0.0/bin/
RUN cp /usr/bin/cvnd /root/.cvnd/cosmovisor/upgrades/v2.0.0/bin/

EXPOSE 26656 26657 1317 9090 8545 8546

CMD ["cosmovisor", "run", "start"]

