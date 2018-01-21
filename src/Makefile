TEST ?= $(shell go list ./... | grep -v sample)
VERSION = $(shell grep 'version =' utils/version.go | sed -E 's/.*"(.+)"$$/\1/')

default: all

all: build

deps:
	go get -d -v github.com/ilovelili/coinmarketcap
	go get -d -v github.com/ahmetb/go-linq

build: deps
	go build -o marketcap 
	cp config.json template ../ansible/playbooks/files/compiled
	mv marketcap ../ansible/playbooks/files/compiled

version:
	@echo $(VERSION)

.PTHONY: all deps build test version