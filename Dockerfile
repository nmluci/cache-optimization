FROM golang:1.17 as build
WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

RUN go mod download
RUN go mod tidy

COPY . /app/

RUN go build -o /app/main

# Deploy

FROM alpine:3.16.0
WORKDIR /app

EXPOSE 5000

RUN apk update
RUN apk add --no--cache tzdata
ENV cp /usr/share/zoneinfo/Asia/Makassar /etc/localtime
RUN echo "Asia/Makassar" > /etc/timezone

COPY --from=build /app/main /app/main

CMD ["./main"]