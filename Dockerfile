FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

ENV GOOS linux

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/wallets ./cmd/app/main.go

FROM alpine AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

WORKDIR /app

COPY ./config.env .

COPY --from=builder /app/wallets /app/wallets

EXPOSE 8080

ENTRYPOINT [ "/app/wallets" ]
