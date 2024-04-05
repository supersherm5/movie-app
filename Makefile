#!/usr/bin/env bash

#############################################
# Metadata Service
#############################################
start-metadata:
		 go run metadata/cmd/main.go
.PHONY: start-metadata


#############################################
# Rating Service
#############################################
start-rating:
		 go run rating/cmd/main.go
.PHONY: start-rating

#############################################
# Movie Service
#############################################
start-movie:
		 go run movie/cmd/main.go
.PHONY: start-movie


#############################################
# Utilities
#############################################

start-all-services:
	make -j start-metadata start-rating start-movie
.PHONY: start-all-services

start-consul:
	@echo "Searching for Hashicorp/Consul image..."
	@if docker container ls -a | grep -q "dev-consul"; then \
  	echo "Dev-consul container was found. Starting consul container."; \
		docker container start dev-consul; \
	elif docker image ls | grep -q "hashicorp/consul"; then \
		echo "Hashicorp/consul image already downloaded. Starting consul container."; \
		docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0; \
	else \
		echo "Hashicorp/Consul image isn't installed. Downloading hashicorp/consul:latest."; \
		docker pull hashicorp/consul; \
		docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0; \
	fi
.PHONY: start-consul

stop-consul:
	@echo "Checking if the dev-consul container is running..."
	@CONSUL_CONTAINER=$$(docker container inspect dev-consul); \
  	if [ -n "$$CONSUL_CONTAINER" ]; then \
		CONTAINER_STATUS=$$(echo "$$CONSUL_CONTAINER" | grep -o '"Status": "[^"]*' | awk -F'"' '{print $$4}'); \
		CONSUL_CONTAINER_Status=$$CONTAINER_STATUS; \
		if [ "$$CONSUL_CONTAINER_Status" = "running" ]; then \
			echo "Stopping dev-consul container."; \
			docker container stop dev-consul; \
		else \
			echo "Container dev-consul is not running."; \
		fi; \
	else \
		echo "Container dev-consul does not exist or an error occurred while inspecting."; \
	fi
.PHONY: stop-consul