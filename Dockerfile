FROM paulcager/go-base:latest as build

WORKDIR /go

COPY go.mod *.go ./src/osgrid-server/
RUN cd /go/src/osgrid-server && \
    go install -v ./... && \
    sha256sum /go/bin/osgrid-server

FROM debian:stable-slim
RUN apt-get update && apt-get -y upgrade
WORKDIR /app
COPY --from=build /go/bin/osgrid-server .
EXPOSE 9090

CMD [ "/app/osgrid-server" ]


