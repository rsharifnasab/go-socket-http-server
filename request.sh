#!/bin/sh

printf "\nget index.html with curl\n"
curl 127.0.0.1/index.html
# or use verbose mode: view sending and receiving headers
printf "\nget humans.txt in verbose mode\n"
curl -v 127.0.0.1/humans.txt

printf "\nget / with http client\n"
http GET 127.0.0.1
