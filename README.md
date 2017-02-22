[![wercker status](https://app.wercker.com/status/94dd4b8a0c490fccd899e59cee34c671/m/master "wercker status")](https://app.wercker.com/project/byKey/94dd4b8a0c490fccd899e59cee34c671)

# README #

## What is this repository for?

`p3sr-srcmanager` is the utility for managing the `p3sr` repositories.
An introduction to the design of the P3SR framework and associated applications
can be found in the wiki for this repo.

## How do I get set up?

### Install and Set Up Golang

Either use the [Go Version Manager](https://github.com/moovweb/gvm) or
navigate to the [Golang Site](https://golang.org/) and set it up manually.
It is prefered that you use the Go version manager, and instructions are on how to do that is on the [wiki](https://bitbucket.org/c3sr/p3sr-srcmanager/wiki/Install%20Go%20Environment).


### Install and Set Up Javascript

For Javascript development you can use the [Node Version Manager](https://github.com/creationix/nvm) or naviage to the [NodeJS Site](https://nodejs.org/en/) and install V6+ manually. 
It is prefered that you use the Node version manager, and instructions are on how to do that is on the [wiki](https://bitbucket.org/c3sr/p3sr-srcmanager/wiki/Install%20Javascript%20Environment).

### Install glide

Glide is available on github [here](https://github.com/Masterminds/glide).

### Clone this Repository

Due to how Go's package system works, it really only works for pulling publically-available code.
Since P3SR is not public yet, you'll need to clone the repository manually while following Go's conventions.

Navigate to where Go will expect to find the source for this repo. Make the path if it does not exist.

    cd $GOPATH/src/bitbucket.org/c3sr

Clone this repository there.

    git clone git@bitbucket.org:c3sr/p3sr-srcmanager.git
    cd p3sr-srcmanager

Pull the rest of the P3SR repositories (read on for how to do that)

## Install and use the p3sr-srcmanager

First, install the dependencies for p3sr-srcmanager. From within the `p3sr-srcmanager` directory, run

    glide install

Now, install `p3sr-srcmanager`.

    go install

Now `p3sr-srcmanager` may be used on the command line to pull the other `p3sr` repositories. Run

    p3sr-srcmanager

Install the protoc compiler. You can find precompiled binaries on [github](https://github.com/google/protobuf/releases).

Install some more dependencies

    go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

To see what it can do for you. You'll probably want to run

    p3sr-srcmanager update
    p3sr-srcmanager goget

to get started.

## Set up Additional Repositories

Some P3SR repositories will require additional steps. A non-exhaustive list is here.
Consult each of those repositories to set them up as needed.

* [`p3sr-services`](https://bitbucket.org/c3sr/p3sr-services)
* `p3sr-cermine`
* `p3sr-pdftoxml`
* `p3sr-yeoman-generator`

## How do I actually do anything with P3SR?

Consult [p3sr-cmd](https://bitbucket.org/c3sr/p3sr-cmd).

## How do I add a new microservice?

Consult [p3sr-yeoman-generator](https://bitbucket.org/c3sr/p3sr-yeoman-generator).

## Troubleshooting

Manually update this repository regularly to keep the source manager and p3sr repository list up to date.

    git pull

`p3sr-srcmanager update` runs `go get -u -v` in all repositories. You can run that command in a troublesome
repository to get more information about a problem.

`p3sr-srcmanager goget` is very verbose. Things may appear to be failures that are not failures.

### Contribution guidelines ###

[To Do]

### Who do I talk to? ###

Please contact Carl Pearson, pearson@illinois.edu.