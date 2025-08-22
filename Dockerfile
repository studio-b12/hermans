FROM golang:1.25-alpine AS build
WORKDIR /build
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .
RUN go build -v -o hermans cmd/hermans/main.go

FROM alpine
WORKDIR /var/hermans
COPY --from=build /build/hermans /opt/hermans
COPY webapp/ /var/hermans/webapp
RUN mkdir -p /var/hermans/db
ENV HMS_BIND_ADDRESS="0.0.0.0:8080"
ENV HMS_DATABASE_DSN="/var/hermans/db/db.sqlite"
ENV HMS_CACHE_DIR="/var/hermans/cache"
ENV HMS_LOG_LEVEL="info"
EXPOSE 8080
ENTRYPOINT [ "/opt/hermans" ]
