# local/fluentd/Dockerfile
FROM --platform=linux/arm64 ruby:3.2-alpine

# 1) Install OS tools + TLS headers
USER root
RUN apk add --no-cache \
      build-base \
      libffi-dev \
      openssl-dev \
      ca-certificates \
      curl \
      netcat-openbsd \
    && update-ca-certificates

# 2) Install Fluentd core + ES plugin
RUN gem install fluentd -v 1.17.0 \
    && gem install fluent-plugin-elasticsearch

# 3) Create fluent user + ensure /fluentd dirs are owned
RUN addgroup -g 1000 fluent \
    && adduser  -u 1000 -G fluent -s /bin/sh -D fluent \
    && mkdir -p /fluentd/etc /fluentd/log \
    && chown -R fluent:fluent /fluentd

# 4) Drop to unprivileged user
USER fluent
WORKDIR /home/fluent

# 5) Expose the TCP/UDP ports for Docker logging driver
EXPOSE 24224 24224/udp 5140

# 6) Embed a simple healthcheck
HEALTHCHECK --interval=15s --timeout=3s --retries=3 \
  CMD fluentd --dry-run -c /fluentd/etc/fluent.conf || exit 1

# 7) Launch Fluentd pointing at the mounted config
CMD ["fluentd", "-c", "/fluentd/etc/fluent.conf", "--no-supervisor", "-v"]