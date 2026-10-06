#!/bin/sh
set -eu

SITE_PACKAGES="$(
    python -c 'import site; print(site.getsitepackages()[0])'
)"

NVIDIA_PYTHON_LIBS="$(
    find "${SITE_PACKAGES}/nvidia" \
        -type d \
        -name lib \
        -print 2>/dev/null |
    sort |
    paste -sd:
)"

export LD_LIBRARY_PATH="${NVIDIA_PYTHON_LIBS}:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH:-}"

exec python main.py "$@"
