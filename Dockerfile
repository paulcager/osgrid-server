FROM paulcager/go-base:latest

EXPOSE 9090
WORKDIR /go/src/osgrid-server

COPY go.mod *.go ./
RUN go install -v ./... && sha256sum /go/bin/osgrid-server && cd /go && rm -rf /go/src /go/pkg

CMD [ "/go/bin/osgrid-server" ]


