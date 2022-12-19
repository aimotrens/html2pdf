#!/bin/bash

if [ ! -f /tmp/swag ]; then
    curl -L https://github.com/swaggo/swag/releases/download/v1.8.9/swag_1.8.9_Linux_x86_64.tar.gz | \
        tar xvfz - -C /tmp/ swag
fi

/tmp/swag init
