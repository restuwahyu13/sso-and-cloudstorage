#########################
# GLOBAL VARIABLE
#########################
nodemon := @nodemon
go := @go
docker := @docker
compose := @docker-compose
soda := @soda
db_config_path := ./database.yaml
db_path := ./databases/migrations
db_env := development

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

###############################
# DOCKER NDP BUILD ENVIRONMENT
###############################
dc-ndpb:
	${compose} up -d --remove-orphans --no-deps --build

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

##########################
# GOLANG BUILD ENVIRONMENT
##########################
build:
	${go} build -o main .

#########################################
# GOLANG INSTALL GLOBAL CLI ENVIRONMENT
#########################################
install-cli:
	${go} install github.com/gobuffalo/pop/v6/soda@latest

##############################
# MIGRATION CREATE ENVIRONMENT
##############################
db-create:
ifdef name
	${soda} create -e ${db_env} -c ${db_config_path} -p ${db_path} -d ${name}
endif

##############################
# MIGRATION GEN ENVIRONMENT
##############################
db-gen:
ifdef name
	${soda} generate -e ${db_env} -c ${db_config_path} -p ${db_path} -d sql ${name}
endif

##############################
# MIGRATION DROP ENVIRONMENT
##############################
db-drop:
	${soda} drop -e ${db_env} -c ${db_config_path} -p ${db_path} -d -a

#######################################
# MIGRATION UP|DOWN|STATUS ENVIRONMENT
#######################################
db-mig:
ifdef type
	${soda} migrate -e ${db_env} -c ${db_config_path} -p ${db_path} -d -s 0 ${type}
endif

db-mig-stat:
	${soda} migrate -e ${db_env} -c ${db_config_path} -p ${db_path} -d status