#!/bin/sh

curl 127.0.0.1/index.html
# or use verbose mode: view sending and receiving headers
curl -v 127.0.0.1/index.html 

http GET 127.0.0.1/humans.txt
