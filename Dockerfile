FROM golang:1.12-alpine
LABEL maintainer "UshioShugo<ushio.s@gmail.com>"


ENV APP_PATH=/app

COPY . ${APP_PATH}
WORKDIR ${APP_PATH}

RUN apk add --no-cache --virtual .goget \
	git && \
	go get -v && \
	apk del .goget

RUN apk add --no-cache gcc g++
