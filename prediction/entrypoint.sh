#!/bin/sh
set -e

NVIDIA_PYTHON_LIBS="$(
    find /usr/local/lib/python3.13/site-packages/nvidia \
        -type d \
        -name lib \
        -print |
        paste -sd:
)"

export LD_LIBRARY_PATH="/usr/local/nvidia/lib64:${NVIDIA_PYTHON_LIBS}:${LD_LIBRARY_PATH}"

exec python main.py "$@"

