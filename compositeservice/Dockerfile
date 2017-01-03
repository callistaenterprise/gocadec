FROM iron/base

EXPOSE 6565
EXPOSE 8181

ADD compositeservice-linux-amd64 compositeservice-linux-amd64

ENTRYPOINT ["./compositeservice-linux-amd64", "-profile=test"]