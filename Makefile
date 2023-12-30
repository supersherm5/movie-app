

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

