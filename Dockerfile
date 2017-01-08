FROM debian:jessie
COPY bin/* /bin/
ENTRYPOINT [ "weixin-gate" ]