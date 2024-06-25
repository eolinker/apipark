# 名称：apinto通用镜像
# 创建时间：2022-10-25
FROM centos:7.9.2009
MAINTAINER liujian

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ARG APP

RUN mkdir -p /${APP}

COPY cmd/* /${APP}/
COPY resource/* /${APP}/

WORKDIR /$APP

CMD ./docker_run.sh