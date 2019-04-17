# pull in the api
FROM mathwebsearch/mwsapi as mwsapi

# Start with elasticsearch
FROM elasticsearch:6.7.0 as final

# Add all the files
ADD /scripts/ /mws/
COPY --from=mwsapi /elasticsync /mws/elasticsync

# Set a single instanmce
ENV discovery.type=single-node

# update the control ports
EXPOSE 9200

# The data volumes
VOLUME /index/
VOLUME /usr/share/elasticsearch/data

# and update the entry point
ENTRYPOINT "/mws/tema_entry.sh"