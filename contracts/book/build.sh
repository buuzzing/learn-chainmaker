#!/bin/bash

contractName=Book

go build -ldflags="-s -w" -o $contractName
7z a ../build/$contractName $contractName
rm $contractName