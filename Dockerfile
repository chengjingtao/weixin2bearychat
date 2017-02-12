FROM debian:jessie
RUN apt-get install -y ca-certificates
RUN mkdir -p /etc/weixin-gate
COPY /weixin-gate /bin/weixin-gate
COPY /etc/weixin-gate/tmpl /etc/weixin-gate/
RUN  chmod +x /bin/weixin-gate
ENTRYPOINT [ "weixin-gate" ]
