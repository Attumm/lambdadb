#!/usr/bin/env bash

set -x
set -e
set -u

curl -vv 'http://127.0.0.1:8000/list/?groupby=woning_type&reduce=count'
curl -vv 'http://127.0.0.1:8000/list/?match-wijkcode=WK036394&groupby=woning_type&reduce=count'
