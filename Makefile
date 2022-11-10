#########################
# GLOBAL VARIABLE
#########################
nodemon := @nodemon
go := @go
docker := @docker
compose := @docker-compose

#########################
# GOLANG DEV ENVIRONMENT
#########################
dev:
	${nodemon} -V -e go} -x ${go} run . --race --signal SIGTERM
#########################
# GOLANG PROD ENVIRONMENT
#########################
prod:
	./main

#########################
# GOLANG BUILD ENVIRONMENT
#########################
build:
	${go} build -o main .

#########################
# DOCKER UP ENVIRONMENT
#########################
dc-up:
	${compose} up -d --remove-orphans --build

#########################
# DOCKER NDP ENVIRONMENT
#########################
dc-ndp:
	${compose} up -d --remove-orphans --no-deps

###########################
# DOCKER BUILD ENVIRONMENT
###########################
dc-down:
	${compose} down

##################################
# DOCKER MANUAL BUILD ENVIRONMENT
##################################
dc-build:
	${docker} build -t goapp-sso --compress .