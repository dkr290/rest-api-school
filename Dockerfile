FROM golang:1.25.4-alpine3.22 AS builder
ENV CGO_ENABLED=0
WORKDIR /build
COPY . .

# Build the service binary.
RUN go build -o rest-api -ldflags "-X main.build=main" ./cmd/api/main.go

FROM alpine:3.22.0
RUN addgroup -g 1000 -S restapi && \
  adduser -u 1000 -h /cmd -G restapi -S restapi

WORKDIR /home/restapi
# Copy from stage 0 builder only the binary files
COPY --from=builder --chown=restapi:restapi /build/rest-api .

RUN mkdir /docs && \
  chown -R restapi:restapi /docs

EXPOSE 8082
USER restapi
CMD [ "./rest-api" ]



