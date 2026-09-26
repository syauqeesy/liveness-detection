#!/bin/sh
set -e

export LD_LIBRARY_PATH="$(
    find /usr/local/lib/python3.13/site-packages/nvidia \
        -type d \
        -name lib \
        -print |
        paste -sd:
):${LD_LIBRARY_PATH}"

exec python main.py
