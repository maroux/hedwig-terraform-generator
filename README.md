# Hedwig Terraform Generator

[![Build Status](https://travis-ci.org/Automatic/hedwig-terraform-generator.svg?branch=master)](https://travis-ci.org/Automatic/hedwig-terraform-generator)

[Hedwig](https://github.com/Automatic/hedwig) is a inter-service communication bus that works on AWS SQS/SNS, while 
keeping things pretty simple and straight forward. It uses [json schema](http://json-schema.org/) draft v4 for 
schema validation so all incoming and outgoing messages are validated against pre-defined schema.

Hedwig Terraform Generator is a CLI utility that makes the process of managing 
[Hedwig Terraform modules](https://registry.terraform.io/search?q=hedwig&verified=false) easier by abstracting 
away details about [Terraform](https://www.terraform.io/).

## Usage 

### Installation

Download the latest version of the release from 
[Github releases](https://github.com/Automatic/hedwig-terraform-generator/releases) - 
it's distributed as a zip containing a Go binary file.

### Configuration

Configuration is specified as a JSON file. Run 

```sh
./hedwig-terraform-generator config-file-structure
```

to get the sample configuration file.

**Advanced usage**: The config *may* contain references to other terraform resources, as long as they resolve to 
an actual resource at runtime. 

### How to use

Run 

```sh
./hedwig-terraform-generator apply-config <config file path>
```

to create Terraform modules. The module is named `hedwig` by default in the current directory.

Re-run on any changes.

## Development

### Prerequisites

Install go1.11.x 

### Getting Started

Assuming that you have go installed, set up your environment:

```sh
$ # in a location NOT in your GOPATH:
$ git checkout github.com/Automatic/hedwig-terraform-generator
$ cd hedwig-terraform-generator
$ go get -mod=readonly -v ./...
$ GO111MODULE=off go get github.com/go-bindata/go-bindata/...
$ GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
$ GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck

```

### Running Tests

You can run tests in using ``make test``. By default, it will run all of the unit and functional tests, but you can 
also run specific tests directly using go test:

```sh
$ go test ./...
$ go test -run TestGenerate ./...
```

## Release Notes

[Github Releases](https://github.com/Automatic/hedwig-terraform-generator/releases)

## How to publish


```sh
make clean build

cd bin/linux-amd64 && zip hedwig-terraform-generator-linux-amd64.zip hedwig-terraform-generator; cd -
cd bin/darwin-amd64 && zip hedwig-terraform-generator-darwin-amd64.zip hedwig-terraform-generator; cd -
```

Upload to Github and attach the zip files created in above step.
