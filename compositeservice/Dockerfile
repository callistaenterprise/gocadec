FROM iron/base

EXPOSE 6565
EXPOSE 8181

ADD compositeservice-linux-amd64 /
ADD healthchecker-linux-amd64 /

HEALTHCHECK --interval=1s --timeout=3s CMD ["./healthchecker-linux-amd64", "-port=6565"] || exit 1

ENTRYPOINT ["./compositeservice-linux-amd64", "-profile=test"]