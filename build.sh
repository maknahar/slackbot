#!/usr/bin/env bash

#Create a local .env file
source $GOPATH/src/github.com/maknahar/jtbot/.env

echo "Doing some cleaning ..."
go clean
echo "Done."

echo "Running goimport ..."
goimports -w=true .
echo "Done."

echo "Running go vet ..."
go vet ./...
if [ $? != 0 ]; then
  exit
fi
echo "Done."

echo "Running go generate ..."
go generate ./...
echo "Done."

echo "Running go format ..."
gofmt -w .
echo "Done."

echo "Running go build ..."
go build -race
if [ $? != 0 ]; then
  echo "## Build Failed ##"
  exit
fi
echo "Done."

echo "Running unit test ..."
go test -p=1 ./...
if [ $? == 0 ]; then
    echo "Done."
	echo "## Starting service ##"
    ./jtbot -v
fi
