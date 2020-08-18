FROM golang


WORKDIR /go/src/drops-backend

ADD . /go/src/drops-backend
RUN ./install_packages.sh

CMD ["go", "run", "server.go"]
