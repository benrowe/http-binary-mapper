FROM alpine:latest

MAINTAINER Ben Rowe <ben.rowe.83@gmail.com>

EXPOSE 8000

VOLUME [ "/data" ]

ADD proxy /
CMD ["/proxy" "-output=/data/output.log -cfg=/data/mappings.yml"]

