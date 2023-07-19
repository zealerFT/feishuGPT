FROM alpine:v0.0.1

MAINTAINER fermi <spuerfei@gmail.com>

    LABEL name="larkGPT" \
description="feishugpt service by fermi" \
owner="spuerfei@gmail.com"

WORKDIR /bin

ADD bin/lark /bin/lark
ADD bin/zoneinfo.zip /bin/zoneinfo.zip
COPY config/*.yaml /bin/config/

COPY docker/entrypoint.sh /entrypoint.sh
# 处理 TimeLocation 无法找到的问题
ENV ZONEINFO=/bin/zoneinfo.zip
ENTRYPOINT ["/entrypoint.sh"]
