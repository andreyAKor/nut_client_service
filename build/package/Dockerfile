# Environment
FROM golang:1.17 as build-env

RUN mkdir -p /opt/nut_client_service
WORKDIR /opt/nut_client_service
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/nut_client_service ./cmd/nut_client_service/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/nut_client_service/configs/nut_client_service.yml /etc/nut_client_service.yml
COPY --from=build-env /opt/bin/nut_client_service /bin/nut_client_service

RUN mkdir -p /var/log/nut_client_service
ENV LOGGING_FILE=/var/log/nut_client_service/nut_client_service.log

ENTRYPOINT ["/bin/nut_client_service", "--config=/etc/nut_client_service.yml"]
