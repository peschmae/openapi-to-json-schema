FROM golang:alpine as build

WORKDIR /root
COPY . /root
RUN go build .

FROM alpine:latest

RUN addgroup -S openapi \
    && adduser -S openapi -G openapi \
    && mkdir /app \
    && chown -R openapi:openapi /app

COPY --from=build /root/openapi-to-json-schema /usr/bin/openapi-to-json-schema

USER openapi
WORKDIR /app

CMD ["openapi-to-json-schema"]
