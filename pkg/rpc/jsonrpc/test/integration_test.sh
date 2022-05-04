#!/usr/bin/env bash

# Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

set -e

# Get the directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# Change into that dir because we expect that.
pushd "$DIR"

echo "==> Building the server"
go build -o rpcserver main.go

echo "==> (Re)starting the server"
PID=$(pgrep rpcserver || echo "")
if [[ $PID != "" ]]; then
	kill -9 "$PID"
fi
./rpcserver &
PID=$!
sleep 2

echo "==> simple request"
R1=$(curl -s 'http://localhost:8008/hello_world?name="my_world"&num=5')
R2=$(curl -s --data @data.json http://localhost:8008)
if [[ "$R1" != "$R2" ]]; then
	echo "responses are not identical:"
	echo "R1: $R1"
	echo "R2: $R2"
	echo "FAIL"
	exit 1
else
	echo "OK"
fi

echo "==> request with 0x-prefixed hex string arg"
R1=$(curl -s 'http://localhost:8008/hello_world?name=0x41424344&num=123')
R2='{"jsonrpc":"2.0","id":"","result":{"Result":"hi ABCD 123"},"error":""}'
if [[ "$R1" != "$R2" ]]; then
	echo "responses are not identical:"
	echo "R1: $R1"
	echo "R2: $R2"
	echo "FAIL"
	exit 1
else
	echo "OK"
fi

echo "==> request with missing params"
R1=$(curl -s 'http://localhost:8008/hello_world')
R2='{"jsonrpc":"2.0","id":"","result":{"Result":"hi  0"},"error":""}'
if [[ "$R1" != "$R2" ]]; then
	echo "responses are not identical:"
	echo "R1: $R1"
	echo "R2: $R2"
	echo "FAIL"
	exit 1
else
	echo "OK"
fi

echo "==> request with unquoted string arg"
R1=$(curl -s 'http://localhost:8008/hello_world?name=abcd&num=123')
R2="{\"jsonrpc\":\"2.0\",\"id\":\"\",\"result\":null,\"error\":\"Error converting http params to args: invalid character 'a' looking for beginning of value\"}"
if [[ "$R1" != "$R2" ]]; then
	echo "responses are not identical:"
	echo "R1: $R1"
	echo "R2: $R2"
	echo "FAIL"
	exit 1
else
	echo "OK"
fi

echo "==> request with string type when expecting number arg"
R1=$(curl -s 'http://localhost:8008/hello_world?name="abcd"&num=0xabcd')
R2="{\"jsonrpc\":\"2.0\",\"id\":\"\",\"result\":null,\"error\":\"Error converting http params to args: Got a hex string arg, but expected 'int'\"}"
if [[ "$R1" != "$R2" ]]; then
	echo "responses are not identical:"
	echo "R1: $R1"
	echo "R2: $R2"
	echo "FAIL"
	exit 1
else
	echo "OK"
fi

echo "==> Stopping the server"
kill -9 $PID

rm -f rpcserver

popd
exit 0