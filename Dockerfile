FROM golang:1.12.6-alpine as builder

# Install curl, git, ssh, certificates, and timezones
RUN apk update && apk add --no-cache curl git ca-certificates openssh tzdata && update-ca-certificates

ENV GO111MODULE=on
WORKDIR /src/code-bot

COPY . .
RUN go build

FROM alpine

FROM alpine:3.9.2
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /src/code-bot/code-bot /usr/bin/code-bot

EXPOSE 8080

ENTRYPOINT [ "code-bot" ]

CMD [ "api" ]
