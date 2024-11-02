FROM node:18 AS build-js

RUN npm install -g gulp gulp-cli --loglevel verbose  # Add verbose logging

WORKDIR /build

COPY . .

RUN npm install --only=dev

FROM golang:latest AS build-golang

WORKDIR /go/src/github.com/gophish/gophish

COPY . .

RUN go get -v && go build -v

FROM debian:stable-slim

RUN useradd -m -d /opt/gophish -s /bin/bash app

RUN apt-get update && \
	apt-get install --no-install-recommends -y jq libcap2-bin ca-certificates && \
	apt-get clean && \
	rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /opt/gophish

COPY --from=build-golang /go/src/github.com/gophish/gophish/ ./
COPY --from=build-js /build/static/js/dist/ ./static/js/dist/
COPY --from=build-js /build/static/css/dist/ ./static/css/dist/
COPY --from=build-golang /go/src/github.com/gophish/gophish/config.json ./

RUN chown app. config.json

RUN setcap 'cap_net_bind_service=+ep' /opt/gophish/gophish

USER app

ENV ADMIN_LISTEN_URL=""
ENV ADMIN_USE_TLS=""
ENV ADMIN_CERT_PATH=""
ENV ADMIN_KEY_PATH=""
ENV ADMIN_TRUSTED_ORIGINS=""
ENV PHISH_LISTEN_URL=""
ENV PHISH_USE_TLS=""
ENV PHISH_CERT_PATH=""
ENV PHISH_KEY_PATH=""
ENV CONTACT_ADDRESS=""
ENV DB_NAME=""
ENV DB_FILE_PATH=""

RUN if [ -n "${ADMIN_LISTEN_URL}" ] ; then \
	jq -r --arg ADMIN_LISTEN_URL "${ADMIN_LISTEN_URL}" '.admin_server.listen_url = $ADMIN_LISTEN_URL' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${ADMIN_USE_TLS}" ] ; then \
	jq -r --argjson ADMIN_USE_TLS "${ADMIN_USE_TLS}" '.admin_server.use_tls = $ADMIN_USE_TLS' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${ADMIN_CERT_PATH}" ] ; then \
	jq -r --arg ADMIN_CERT_PATH "${ADMIN_CERT_PATH}" '.admin_server.cert_path = $ADMIN_CERT_PATH' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${ADMIN_KEY_PATH}" ] ; then \
	jq -r --arg ADMIN_KEY_PATH "${ADMIN_KEY_PATH}" '.admin_server.key_path = $ADMIN_KEY_PATH' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${ADMIN_TRUSTED_ORIGINS}" ] ; then \
	jq -r --arg ADMIN_TRUSTED_ORIGINS "${ADMIN_TRUSTED_ORIGINS}" '.admin_server.trusted_origins = ($ADMIN_TRUSTED_ORIGINS|split(","))' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${PHISH_LISTEN_URL}" ] ; then \
	jq -r --arg PHISH_LISTEN_URL "${PHISH_LISTEN_URL}" '.phish_server.listen_url = $PHISH_LISTEN_URL' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${PHISH_USE_TLS}" ] ; then \
	jq -r --argjson PHISH_USE_TLS "${PHISH_USE_TLS}" '.phish_server.use_tls = $PHISH_USE_TLS' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${PHISH_CERT_PATH}" ] ; then \
	jq -r --arg PHISH_CERT_PATH "${PHISH_CERT_PATH}" '.phish_server.cert_path = $PHISH_CERT_PATH' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${PHISH_KEY_PATH}" ] ; then \
	jq -r --arg PHISH_KEY_PATH "${PHISH_KEY_PATH}" '.phish_server.key_path = $PHISH_KEY_PATH' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${CONTACT_ADDRESS}" ] ; then \
	jq -r --arg CONTACT_ADDRESS "${CONTACT_ADDRESS}" '.contact_address = $CONTACT_ADDRESS' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${DB_NAME}" ] ; then \
	jq -r --arg DB_NAME "${DB_NAME}" '.db_name = $DB_NAME' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi && \
	if [ -n "${DB_FILE_PATH}" ] ; then \
	jq -r --arg DB_FILE_PATH "${DB_FILE_PATH}" '.db_path = $DB_FILE_PATH' config.json > config.json.tmp && \
	mv config.json.tmp config.json; \
	fi

RUN sed -i 's/127.0.0.1/0.0.0.0/g' config.json

EXPOSE 3333 8080 8443 80

CMD ["./gophish"]
