FROM golang:1.22.4 AS builder

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./
RUN go mod download

COPY . ./
COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp1 .

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

RUN mkdir -p /root/logs
RUN touch /root/logs/app.log

COPY --from=builder /app/myapp1 .
COPY --from=builder /app/.env .

EXPOSE 8082


CMD [ "./myapp1" ]