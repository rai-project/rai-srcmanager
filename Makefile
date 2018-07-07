PROJECT := rai-srcmanager

.PHONY: default dep build run

default: build

dep:
	dep ensure -v

build:
	go build

run:
	bazel run ${PROJECT}
