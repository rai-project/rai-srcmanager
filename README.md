
# README [![Build Status](https://travis-ci.org/rai-project/rai-srcmanager.svg?branch=master)](https://travis-ci.org/rai-project/rai-srcmanager)

## What is this repository for?

`rai-srcmanager` is the utility for managing the `rai` repositories.
An introduction to the design of the RAI framework and associated applications
can be found in the wiki for this repo.

## How do I get set up?

### Install and Set Up Golang

Either use the [Go Version Manager](https://github.com/moovweb/gvm) or
navigate to the [Golang Site](https://golang.org/) and set it up manually.
It is prefered that you use the Go version manager.

### Install glide

Glide is available on github [here](https://github.com/Masterminds/glide).

### Build from source

    go get -u -v github.com/rai-project/rai-srcmanager

### Clone this Repository using GIT

Navigate to where Go will expect to find the source for this repo. Make the path if it does not exist.

    mkdir -p $GOPATH/src/github.com/rai-project
    cd $GOPATH/src/github.com/rai-project

Clone this repository there.

    git clone git@github.com:rai-project/rai-srcmanager.git
    cd rai-srcmanager

Pull the rest of the RAI repositories (read on for how to do that)

## Install and use the rai-srcmanager

First, install the dependencies for rai-srcmanager. From within the `rai-srcmanager` directory, run

    glide install

Now, install `rai-srcmanager`.

    go install

Now `rai-srcmanager` may be used on the command line to pull the other `rai` repositories. Run

    rai-srcmanager

To see what it can do for you. You'll probably want to run

    rai-srcmanager update
    rai-srcmanager goget

to get started. By default `rai-srcmanager` checks out the repos using the `ssh` protocol. You can change that (checking out using the `https` protocol) by using `rai-srcmanager [[cmd]] --public` (for example `rai-srcmanager update --public`)


## Troubleshooting

Manually update this repository regularly to keep the source manager and rai repository list up to date.

    git pull

`rai-srcmanager update` runs `go get -u -v` in all repositories. You can run that command in a troublesome
repository to get more information about a problem.

`rai-srcmanager goget` is very verbose. Things may appear to be failures that are not failures.

### Contribution guidelines ###

[To Do]
