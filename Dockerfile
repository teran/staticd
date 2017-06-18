FROM scratch

ARG BUILD_DATE
ARG VCS_REF

EXPOSE 8080

ADD bin/staticd-linux-amd64 /staticd

ENTRYPOINT ["/staticd"]
