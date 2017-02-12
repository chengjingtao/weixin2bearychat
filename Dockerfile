FROM debian:jessie
ENV TZ=Asia/Shanghai
RUN apt-get install -y ca-certificates \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
	&& echo $TZ > /etc/timezone
RUN mkdir -p /etc/weixin-gate
COPY /weixin-gate /bin/weixin-gate
COPY /etc/weixin-gate/tmpl /etc/weixin-gate/
RUN  chmod +x /bin/weixin-gate
ENTRYPOINT [ "weixin-gate" ]
