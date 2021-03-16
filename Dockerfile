FROM paulcager/go-base:latest as build

COPY * ./
RUN CGO_ENABLED=0 go install -v ./... && \
    sha256sum /go/bin/osgrid-server

FROM  scratch
WORKDIR /app
COPY --from=build /go/bin/osgrid-server .
EXPOSE 9090

CMD [ "/app/osgrid-server" ]


