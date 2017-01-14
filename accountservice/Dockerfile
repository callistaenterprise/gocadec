FROM iron/base
ENV SERVER_PORT 7777
EXPOSE $SERVER_PORT

ADD accountservice-linux-amd64 /
ADD healthchecker-linux-amd64 /

HEALTHCHECK --interval=1s --timeout=3s CMD ["./healthchecker-linux-amd64", "-port=7777"] || exit 1

ENTRYPOINT ["./accountservice-linux-amd64", "-profile=test"]