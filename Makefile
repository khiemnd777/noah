SHELL := /bin/bash

API_DIR := api
MAKE_API := $(MAKE) -C $(API_DIR)

.PHONY: docker-up docker-down docker-restart docker-stop docker-log

docker-up:
	$(MAKE_API) docker-up

docker-down:
	$(MAKE_API) docker-down

docker-restart:
	$(MAKE_API) docker-restart

docker-stop:
	$(MAKE_API) docker-stop

docker-log:
	$(MAKE_API) docker-log
