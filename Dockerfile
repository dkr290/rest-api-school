FROM golang:1.25.4-alpine3.22 AS builder
ENV CGO_ENABLED=0
WORKDIR /build
COPY . .

# Build the service binary.
RUN go build -o rest-api -ldflags "-X main.build=main" ./cmd/api/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /home/rest-api
# Copy from stage 0 builder only the binary files
COPY --from=builder  /build/rest-api .


EXPOSE 8082
CMD [ "./rest-api" ]
