#!/usr/bin/env bash

set -x
set -e
set -u

# should be cached.
curl -vv 'http://127.0.0.1:8000/list/?groupby=woning_type&reduce=count'

# should not be cached.(using bitmaps)
curl -vv 'http://127.0.0.1:8000/list/?match-wijkcode=WK036394&groupby=woning_type&reduce=count'
