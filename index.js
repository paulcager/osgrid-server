import OsGridRef from 'geodesy/osgridref.js';
import { LatLon } from 'geodesy/osgridref.js';
const url = require('url');


const express = require('express');
const app = express();


// Convert between Ordnance Survey grid reference and WGS-84 lat/lon.
// See https://www.movable-type.co.uk/scripts/latlong-os-gridref.html
//
// Usage: http://server/?gridRef=osgb-grid-ref or ://server/?latlon=lat,lon
app.get('/', (request, response) => {
    var urlParts = url.parse(request.url, true);
    var parameters = urlParts.query;

    if (parameters.gridRef) {
        const gridref = OsGridRef.parse(parameters.gridRef);
        const wgs84 = gridref.toLatLon();
        response.send(String(wgs84.lat) + "," + String(wgs84.lon));
    } else if (parameters.latLon) {
        const parts = parameters.latLon.split(",");
        const wgs84 = new LatLon(parts[0], parts[1]);
        const gridref = wgs84.toOsGrid();
        response.send(gridref.toString());
    } else {
        response.send(404);
    }
});

app.listen(9090, () => console.log('Listening on port 9090'));

