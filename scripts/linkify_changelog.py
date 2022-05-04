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

import fileinput
import re

# This script goes through the provided file, and replaces any " \#<number>",
# with the valid mark down formatted link to it. e.g.
# " [\#number](https://github.com/bhojpur/state/pull/<number>)
# Note that if the number is for a an issue, github will auto-redirect you when you click the link.
# It is safe to run the script multiple times in succession. 
#
# Example usage $ python3 linkify_changelog.py ../CHANGELOG_PENDING.md
for line in fileinput.input(inplace=1):
    line = re.sub(r"\s\\#([0-9]*)", r" [\\#\1](https://github.com/bhojpur/state/pull/\1)", line.rstrip())
    print(line)