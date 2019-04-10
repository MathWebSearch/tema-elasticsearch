# tema-elasticsync

[![Build Status](https://travis-ci.org/MathWebSearch/mws-temasync.svg?branch=master)](https://travis-ci.org/MathWebSearch/mws-temasync)

A component configuring and maintaining an elasticsearch instance for use with TemaSearch. 

## Building

Tema-elasticsearch is writting in [Go](https://golang.org). 
It can be built using the `go build` tool or alternatively, if Make is installed, using `make`. 

## Usage

```
Usage of ./tema-elasticsync:
  -elastic-host string
        Host to use for elasticsearch (default "0.0.0.0")
  -elastic-port int
        Port to use for elasticsearch (default 9200)
  -index-dir string
        Directory to use for Indexes (default "/index/")
```

## Process

Segmented Sync, to be documented. 


## Dockerfile

For convenienve, a Dockerfile is also provided. 
It can be found as the automated build [mathwebsearch/tema-elasticsync](https://hub.docker.com/r/mathwebsearch/tema-elasticsync) on DockerHub. 

The docker image:
- expects a TemaSearch Index inside the `/index/` volume
- starts a single-node elasticsearch instance listening on port `9200`, maintaining storing data in `/usr/share/elasticsearch/data`
- syncronises the elasticsearch instance with the TemaSearch index on every container start

Use it e.g. as follows:

```
    docker run mathwebsearch/tema-elasticsync -p 9200:9200 -v /path/to/indicies:/index/ -v /usr/share/elasticsearch/data mathwebsearch/tema-elasticsync
```