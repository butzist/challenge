# Solution to challenge

This repository contains my suggested solution to a coding challenge.

## Prerequisites

* Golang 1.9+
* Git

Additionally, for running against a test data set you will need:

* Docker
* docker-compose

## Installation
The program can be installed using the standard procedure for Golang:

    go get -u github.com/butzist/challenge

This will download, build, and install the program. If successful, it can be found in $GOPATH/bin/challenge. This usually defaults to ~/go/bin/challenge.

## Usage

The implementation is composed of four components, that each have multiple implementations, that can be selected via command line:

* Source
    * kafka: Read input data from Kafka topic
    * canned: Download and read data from S3 url
* Output
    * kafka: Output counts to Kafka topic
    * console: Output to stdout
* Processing
    * simple: Accept only monotonically increasing timestamps and output data when timestamp within next minute received
    * advanced: Accept only timestamps within a window around the current system time and output after timeout of 5s
* Counters
    * exact: Use exact, but memory-hungry counter
    * probabilistic: Use inexact, but friendly, counter
   
    
### Kafka specific environment variables

* KAFKA_BROKER: The broker to use, defaults to "localhost:9092"
* KAFKA_TOPIC: The topic to read from, defaults to "mytopic"
* KAFKA_PARTITION: The partition to read from, defaults to 0
* KAFKA_OUTPUT_TOPIC: The topic to write results to, defaults to "mycounts"