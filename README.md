# tema-elasticsync

[![Build Status](https://travis-ci.org/MathWebSearch/tema-elasticsync.svg?branch=master)](https://travis-ci.org/MathWebSearch/tema-elasticsync)

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

The program creates a and updates an Elasticsearch Index for Tema-Search. 

A Tema-Search Index is a set of JSON Objects conforming to the Temasearch schema. 
Each set of documents belonging to a single `.harvest` file (which in and of itself belongs to one source file) should be contained on one line of an index file ending in `.json`. 
For backward compatibility, in between each lines of items in the index, a document containing meta-information should be stored. 
In the following, we will call the collection of these `.json` file a `TemaSearch Index`. 

The Tema-Index should be kept in sync with an appropriate ElasticSearch Index. 
This means that the content of all documents needs to be indexed by elasticsearch. 
In principle, upon syncing one could:
1. Delete all existing indexed documents from ElasticSearch (if any)
2. Read each document in the entire index and add the documents contained inside of it to ElasticSearch
This approach does not scale well with large datasets. 
Having to delete the entire database, only to add the same content back is too slow. 

Instead we split the index into different files, which we call segments below. 
To syncronize an updated index into Elasticsearch, we roughly do the following:
1. Mark all existing segments in the database as 'untouched' within the syncronization
2. For each segment from the ElasticSearch index to be added:
  - compute a hash of the segment
  - check if this segment with the same name is already stored in the database by comparing the hash
    - if yes, we do not need to do anything as it has not changed
    - if no, we remove the old segment documents (if any) and add the new documents belonging to this hash
  - mark the segment as 'touched' within this syncronization process
3. Delete the documents belonging to any segment still marked as 'untouched'
This process is far more efficient -- only updating documents in the database that have actually been changed. 

In practice, this process requires that two seperate ElasticSearch indexes are maintained. 
The first index -- called `tema` by convention -- contains the TemaSearch Index Documents and is most obvious. 
The second index is called `tema-segments` and contains a list of known segments as well as their hashes. 
As a hash implementation we use `SHA256`. 

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