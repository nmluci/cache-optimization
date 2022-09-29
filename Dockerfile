FROM golang:1.19 as build
WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

RUN go mod download
RUN go mod tidy

COPY . /app/

RUN CGO_ENABLED=0 go build -o /app/main

# Deploy

FROM alpine:3.16.0 as webservice
WORKDIR /app

EXPOSE 3000

RUN apk update
RUN apk add --no-cache tzdata

ENV TZ=Asia/Makassar
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime

COPY --from=build /app/main /app/main
COPY --from=build /app/conf /app/conf

CMD ["./main"]