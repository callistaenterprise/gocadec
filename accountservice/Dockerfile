FROM iron/base
ENV SERVER_PORT 7777
EXPOSE $SERVER_PORT

ADD accountservice-linux-amd64 accountservice-linux-amd64
ENTRYPOINT ["./accountservice-linux-amd64", "-profile=test"]