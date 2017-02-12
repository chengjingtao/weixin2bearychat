FROM debian:jessie
COPY /weixin-gate /bin/weixin-gate
RUN  chmod +x /bin/weixin-gate
ENTRYPOINT [ "weixin-gate" ]
