// Convert between Ordnance Survey grid reference and WGS-84 lat/lon.
// See https://www.movable-type.co.uk/scripts/latlong-os-gridref.html

import OsGridRef from 'geodesy/osgridref.js';
import { LatLon } from 'geodesy/osgridref.js';
const url = require('url');
const express = require('express');

const app = express();

// Usage: http://server/gridref=osgb-grid-ref
app.get('/gridref/:ref', (request, response) => {
        const gridref = OsGridRef.parse(request.params.ref);
        const wgs84 = gridref.toLatLon();
        response.send(String(wgs84.lat) + "," + String(wgs84.lon));
});

// Usage: http://server/latlon=lat,lon
app.get('/latlon/:latlon', (request, response) => {
    const parts = request.params.latlon.split(",");
    const wgs84 = new LatLon(parts[0], parts[1]);
    const gridref = wgs84.toOsGrid();
    response.send(gridref.toString());
});

app.listen(9090, () => console.log('Listening on port 9090'));

