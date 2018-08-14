.PHONY: default install test

PKG_NAME = github.com/AdhityaRamadhanus/minerva
TEST_PKG = ${PKG_NAME}/redis

# target #

default: install test

install:
	go get -t ./...

# Test Packages

test:
	go test -v --cover ${TEST_PKG}