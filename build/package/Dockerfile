# Environment
FROM golang:1.17 as build-env

RUN mkdir -p /opt/nut_parser
WORKDIR /opt/nut_parser
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/nut_parser ./cmd/nut_parser/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/nut_parser/configs/nut_parser.yml /etc/nut_parser.yml
COPY --from=build-env /opt/bin/nut_parser /bin/nut_parser

RUN mkdir -p /var/log/nut_parser
ENV LOGGING_FILE=/var/log/nut_parser/nut_parser.log

ENTRYPOINT ["/bin/nut_parser", "--config=/etc/nut_parser.yml"]
