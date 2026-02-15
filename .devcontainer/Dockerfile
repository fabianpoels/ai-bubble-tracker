FROM golang:1.26.0-bookworm

# Install additional tools
RUN apt-get update && apt-get -y install --no-install-recommends \
    postgresql-client \
    redis-tools \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Create workspace directory
RUN mkdir -p /workspace && chmod -R 777 /go

CMD ["bash"]