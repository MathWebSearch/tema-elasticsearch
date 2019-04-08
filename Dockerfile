FROM golang:1-alpine as builder

# Build dependencies
RUN apk add --no-cache make git

# Build this updater
ADD . /go/src/github.com/MathWebSearch/tema-elasticsearch
WORKDIR /go/src/github.com/MathWebSearch/tema-elasticsearch
RUN make build-local

# Start with elasticsearch
FROM elasticsearch:6.7.0 as final

# Add all the files
ADD /scripts/ /mws/
COPY --from=builder /go/src/github.com/MathWebSearch/tema-elasticsearch/out/tema_hook /mws/tema_hook

# Set a single instanmce
ENV discovery.type=single-node

# update the control ports
EXPOSE 9200
EXPOSE 9300

# The tema-search index
VOLUME /index/

# and update the entry point
ENTRYPOINT "/mws/tema_entry.sh"