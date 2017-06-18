FROM alpine:latest

ARG BUILD_DATE
ARG VCS_REF

EXPOSE 8080

RUN apk add --update --no-cache ca-certificates && \
    rm -rvf /var/cache/apk/*

ADD bin/staticd-linux-amd64 /staticd

ENTRYPOINT ["/staticd"]
