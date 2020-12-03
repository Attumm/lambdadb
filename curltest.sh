#!/usr/bin/env bash

set -x
set -e
set -u


curl \
	--data-urlencode 'geojson={
				"type": "Polygon",
				"coordinates": [
					[
					    [4.902321, 52.428306],
					    [4.90127, 52.427024],
					    [4.905281, 52.426069],
					    [4.906782, 52.426226],
					    [4.906418, 52.427469],
					    [4.902321, 52.428306]
					]
				]
			}' \
	'http://127.0.0.1:8000/list/?groupby=postcode&reduce=count'
