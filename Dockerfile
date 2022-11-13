###############################
# DOCKER START STAGE
###############################
FROM golang:1.19.3-buster
WORKDIR /usr/src/goapp/
USER ${USER}
ADD ./go.mod /usr/src/goapp/
ADD . /usr/src/goapp/

###############################
# DOCKER ENVIRONMENT STAGE
###############################
ENV GO111MODULE="on" \
  CGO_ENABLED="0" \
  GO_GC="off"

###############################
# DOCKER UPGRADE STAGE
###############################
RUN apt-get autoremove \
  && apt-get autoclean \
  && apt-get update --fix-missing \
  && apt-get upgrade -y \
  && apt-get install curl \
  build-essential -y

################################
# DOCKER INSTALL NODEJS STAGE
################################
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install nodejs

###############################
# DOCKER INSTALL & BUILD STAGE
###############################
RUN go mod download \
  && go mod tidy \
  && go mod verify \
  && go build -o main .
RUN npm i pm2 -g --verbose --no-audit

###############################
# DOCKER FINAL STAGE
###############################
EXPOSE 3000
CMD pm2-runtime pm2.config.js