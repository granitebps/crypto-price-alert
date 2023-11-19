FROM golang:1.19-alpine as build

RUN apk update && apk add ca-certificates && apk add tzdata && apk add build-base

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .
RUN cp /build/.env .
RUN cp /build/alias.json .
RUN cp /build/alert.json .

FROM bash as runtime

RUN apk add curl

WORKDIR /app

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /dist/main .
COPY --from=build /build/.env .
COPY --from=build /build/alias.json .
COPY --from=build /build/alert.json .


ENV TZ Asia/Jakarta

CMD ["./main"]