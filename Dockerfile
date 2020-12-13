FROM paulcager/go-base:latest

EXPOSE 9090
WORKDIR /go/src/osgrid-server

COPY go.mod *.go ./
RUN go install -v ./... && sha256sum /go/bin/osgrid-server

CMD [ "/go/bin/osgrid-server" ]


