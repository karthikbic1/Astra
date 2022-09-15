#!/bin/bash
set -e

mkdir -p ~/.go
echo "GOPATH=$HOME/.go" >> ~/.bashrc
echo "export GOPATH" >> ~/.bashrc
echo "PATH=\$PATH:\$GOPATH/bin # Add GOPATH/bin to PATH for scripting" >> ~/.bashrc
source ~/.bashrc


if [[ $OSTYPE == 'darwin'* ]]; then
    echo 'macOS'
    brew update
    brew install golang
else
    sudo snap install go --classic
fi



