BINARY = bin/starc
FILE ?= tesc.starc
VERSION = 2.0.0
VERS_STATE = PRE-ALPHA



.PHONY: ignite build install help clean version test

build:
    go build -o $(BINARY)
    @echo "Star-C, Successfull operation"

ignite:
    build
    ./$(BINARY) ignite $(FILE)