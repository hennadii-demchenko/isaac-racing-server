#!/bin/bash

cd "$GOPATH/src/github.com/Zamiell/isaac-racing-server/src"
go install
if [ $? -eq 0 ]; then
        mv "$GOPATH/bin/src" "$GOPATH/bin/isaac-racing-server"
	supervisorctl restart isaac-racing-server
else
	echo "isaac-racing-server - Go compilation failed!"
fi
