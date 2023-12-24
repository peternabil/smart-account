#!/bin/sh

openssl req -x509 -nodes -newkey rsa:2048 -keyout key.pem -out cert.pem -sha256 -days 365 \
    -subj "/C=EG/ST=Cairo/L=Cairo/O=Peter/OU=IT Department/CN=localhost"

# docker build . -t my_app

docker compose up -d
