# osgrid-server - _Ordnance Survey Grid Ref Converter_

This project implements a simple REST server to expose the excellent https://www.movable-type.co.uk/scripts/latlong-os-gridref.html
Ordnance Survey grid reference conversion utilities.

## Usage

The easiest way to avoid environmental or Node versioning problems is to start the server using Docker:

    docker run -d --rm --name osgrid-server -p 9090:9090 paulcager/osgrid-server

A grid reference may then be converted to a lat /lon:

    $ curl http://localhost:9090/gridref/SU3724515505
    50.93779289256305,-1.471309033422234

Or a lat/lon converted to a grid reference:

    $ curl http://localhost:9090/latlon/50.9378,-1.4713
    SU 37245 15505

