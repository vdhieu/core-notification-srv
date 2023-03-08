ARG REGISTRY_REPOSITORY
FROM ${REGISTRY_REPOSITORY}/base-go-image:latest AS builder

COPY . .

COPY go.mod go.sum ./
RUN go env -w GOPRIVATE=github.com/Neutronpay/*
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-srv main.go

FROM alpine:3.15.0 as runner
RUN apk --no-cache add ca-certificates
WORKDIR /build
COPY --from=builder /build/notification-srv .
COPY config.json /build/

ENTRYPOINT [ "./notification-srv" ]
