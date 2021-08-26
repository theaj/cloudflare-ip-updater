############################
## Prebuild
############################
FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates
RUN adduser -D -g '' monitoruser

WORKDIR /src

# this is to cache any go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/monitor main.go

RUN ls -lah /go/bin/monitor

############################
## Actual Image
############################
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/monitor /go/bin/monitor

USER monitoruser

ENTRYPOINT ["./go/bin/monitor"]