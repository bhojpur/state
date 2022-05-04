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

VERSION=v12.9.0
NODE_FULL=node-${VERSION}-linux-x64

mkdir -p ~/.local/bin
mkdir -p ~/.local/node
wget http://nodejs.org/dist/${VERSION}/${NODE_FULL}.tar.gz -O ~/.local/node/${NODE_FULL}.tar.gz
tar -xzf ~/.local/node/${NODE_FULL}.tar.gz -C ~/.local/node/
ln -s ~/.local/node/${NODE_FULL}/bin/node ~/.local/bin/node
ln -s ~/.local/node/${NODE_FULL}/bin/npm ~/.local/bin/npm
export PATH=~/.local/bin:$PATH
npm i -g dredd
ln -s ~/.local/node/${NODE_FULL}/bin/dredd ~/.local/bin/dredd