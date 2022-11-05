###############################
# DOCKER START STAGE
###############################
FROM golang:1.19.3-buster
WORKDIR /usr/src/goapp/
USER ${USER}
COPY ./go.mod \
  ./go.sum /usr/src/goapp/
COPY . /usr/src/goapp/

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
  && apt-get update \
  && apt-get upgrade -y \
  && apt-get install build-essential -y

###############################
# DOCKER INSTALL & BUILD STAGE
###############################
RUN go mod tidy \
  && go mod download \
  && go build -o main .

###############################
# DOCKER FINAL STAGE
###############################
EXPOSE 3000
CMD ["./main"]