FROM iron/base

EXPOSE 6767

ADD imageservice-linux-amd64 /
ADD testimages/*.jpg testimages/
ADD healthchecker-linux-amd64 /

HEALTHCHECK --interval=1s --timeout=3s CMD ["./healthchecker-linux-amd64", "-port=6767"] || exit 1

ENTRYPOINT ["./imageservice-linux-amd64", "-profile=test"]