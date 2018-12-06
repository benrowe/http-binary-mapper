FROM alpine:latest

MAINTAINER Ben Rowe <ben.rowe.83@gmail.com>

EXPOSE 80

RUN mkdir /src
WORKDIR /src

VOLUME /src/data

COPY bin/proxy /src/proxy

RUN chmod +x /src/proxy

CMD /src/proxy -output=/src/data/output.log -cfg=/src/data/mappings.yaml

