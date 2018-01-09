PROJECT := rai-srcmanager

.PHONY: default dep build run

default: build

dep:
	dep ensure
	bazel run //:gazelle_fix

build:
	bazel build ${PROJECT}

run:
	bazel run ${PROJECT}
