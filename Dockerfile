FROM golang:1.24.1-bullseye AS builder
WORKDIR /base
COPY . .
RUN make clean-build

FROM debian:stable-20241202-slim
ARG APP_USER=app-user
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /base/bin/app /usr/bin
RUN groupadd $APP_USER && useradd -g $APP_USER $APP_USER
USER $APP_USER

ENTRYPOINT ["/usr/bin/app"]
