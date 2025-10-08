#!/usr/bin/sh
CGO_ENABLED=0 GOOS=linux go build -o server main.go &&
cp server public/ &&
cp -R bot public/ &&
cp -R images public/ &&
cp -R lib public/ &&
cp -R templates public/ &&
cp .env_prod public/.env