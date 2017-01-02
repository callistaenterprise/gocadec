FROM iron/base

EXPOSE 6767

ADD imageservice-linux-amd64 imageservice-linux-amd64
ADD testimages/*.jpg testimages/

ENTRYPOINT ["./imageservice-linux-amd64", "-profile=test"]