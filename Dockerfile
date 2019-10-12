FROM harbor-dev.hfjy.com:5050/base/hfjy-base:alpine-3.8

MAINTAINER liukun <liukun@hfjy.com>

add container-orchestration-service /opt/app/
add *.yaml /opt/app/



WORKDIR /opt/app
EXPOSE 8080

CMD ["./container-orchestration-service"]
