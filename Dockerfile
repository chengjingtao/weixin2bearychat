FROM debian:jessie
ENV TZ=Asia/Shanghai
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
	&& echo $TZ > /etc/timezone

RUN mkdir -p /etc/weixin2bearychat
COPY /weixin2bearychat /bin/weixin2bearychat
COPY /etc/weixin2bearychat/tmpl /etc/weixin2bearychat/
RUN  chmod +x /bin/weixin2bearychat
ENTRYPOINT [ "weixin2bearychat" ]
