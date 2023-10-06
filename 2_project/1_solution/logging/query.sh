#!/bin/bash

while true;
do
    curl 10.0.1.173:8080/v1/color/red
    sleep 1
    curl 10.0.1.173:8080/v1/error/5xx
    sleep 1
    curl 10.0.1.173:8080/v1/color/orange
    sleep 1
    curl 10.0.1.173:8080/v1/error/4xx
    sleep 1
    curl 10.0.1.173:8080/v1/color/melon
    sleep 1
    curl 10.0.1.173:8080/v1/error/3xx
    sleep 1
done