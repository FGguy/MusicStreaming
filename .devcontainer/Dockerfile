FROM golang:1.23-bookworm AS dev_containers_target_stage

ENV DEBIAN_FRONTEND=noninteractive

# Install additional tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install sqlc via direct download
RUN curl -L https://github.com/sqlc-dev/sqlc/releases/download/v1.27.0/sqlc_1.27.0_linux_amd64.tar.gz -o sqlc.tar.gz && \
    tar -xzf sqlc.tar.gz -C /usr/local/bin && \
    rm sqlc.tar.gz

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin

WORKDIR /workspace