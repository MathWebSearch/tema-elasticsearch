# tema-elasticsync

This repository provides a Dockerfile for configuring and maintaining an elasticsearch instance for use with TemaSearch. 
It can be found as the automated build [mathwebsearch/tema-elasticsync](https://hub.docker.com/r/mathwebsearch/tema-elasticsync) on DockerHub. 

The docker image:
- expects a TemaSearch Index inside the `/index/` volume
- starts a single-node elasticsearch instance listening on port `9200`, maintaining storing data in `/usr/share/elasticsearch/data`
- syncronises the elasticsearch instance with the TemaSearch index on every container start

Use it e.g. as follows:

```
    docker run mathwebsearch/tema-elasticsync -p 9200:9200 -v /path/to/indicies:/index/ -v /usr/share/elasticsearch/data mathwebsearch/tema-elasticsync
```