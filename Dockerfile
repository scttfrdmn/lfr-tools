FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY lfr-tools .

ENTRYPOINT ["./lfr-tools"]