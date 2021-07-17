FROM alpine:3.14.0

WORKDIR /app

COPY ./serve /app

COPY ./static/* /app/static/

ENTRYPOINT ["/app/serve"]