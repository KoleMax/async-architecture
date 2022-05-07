#!/bin/bash
set -e

case "$1" in
    start)
        /app/api
        ;;
    lint)
        cd /build
        make ci-lint
        ;;
    tests)
        cd /build
        make ci-test
        ;;
    *)
        exec "$@"
esac
