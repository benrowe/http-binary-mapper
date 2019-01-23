FROM alpine:latest

MAINTAINER Ben Rowe <ben.rowe.83@gmail.com>

EXPOSE 80

RUN apk update && \
    apk add -y openssl && \
    openssl genrsa -des3 -passout pass:x -out server.pass.key 2048 && \
    openssl rsa -passin pass:x -in server.pass.key -out server.key && \
    rm server.pass.key && \
    openssl req -new -x509 -sha256 -key server.key -out server.csr \
        -subj "/C=UK/ST=Warwickshire/L=Leamington/O=OrgName/OU=IT Department/CN=example.com" && \
    openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt


RUN mkdir /src
WORKDIR /src

VOLUME /src/data

COPY bin/proxy /src/proxy

RUN chmod +x /src/proxy

CMD /src/proxy -output=/src/data/output.log -cfg=/src/data/mappings.yaml

