FROM golang AS build-env

ADD . /app

WORKDIR /app

RUN go build -o watcher ./cmd/lonelywatcher


# target image
FROM debian:10-slim

# Install curl and install/updates certificates
RUN apt-get update \
    && apt-get install -y -q --no-install-recommends \
    ca-certificates \
    curl \
    && apt-get clean

COPY --from=build-env /app/watcher /usr/bin/watcher

CMD ["watcher"]
