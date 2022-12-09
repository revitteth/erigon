#!/bin/sh
# run the debug binary
/go/bin/dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec /usr/local/bin/erigon --