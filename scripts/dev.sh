#!/bin/bash
set -e

if [[ $1 == "init" ]]; then
    echo "Init the dev setup..."
    echo `pwd`
    bash scripts/setup.sh
    echo "Successfully installed dependencies.. "
    echo "Start the app using - dev up"
    
elif [[ $1 == "up" ]]; then
    echo "Building Astra app.."
    go build
    echo "Starting the Astra app.."
    
    port=$2
    if [[ -z $port ]]; then
        port=3000
    fi
    ./astra server --port=$port
elif [[ $1 == "test" ]]; then
    echo "Unit Testing Astra app.."
    go test ./...
else
    echo "Invalid choice. Select between dev init or dev up or dev test"
fi
