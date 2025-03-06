FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG AppVersion=$(cat VERSION)

RUN go build -ldflags="-X 'main.AppVersion=${AppVersion}'" -o necrest .

FROM debian:bookworm-slim AS runtime

COPY --from=builder /app/necrest /necrest
COPY --from=builder /app/VERSION /VERSION

COPY config/rest.yaml /etc/noerrorcode/rest.yaml

ENTRYPOINT ["/necrest", "serve", "--config", "/etc/noerrorcode/rest.yaml", "--log", "trace"]

ARG AppVersion
LABEL version="${AppVersion}"
