FROM alpine:latest
WORKDIR /app
COPY ./cmd/server/bin/app-linux-arm /app/app
RUN chmod u+x app
ENTRYPOINT ["./app"]
