package math

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Fraction defined in terms of a numerator divided by a denominator in uint64
// format. Fraction must be positive.
type Fraction struct {
	// The portion of the denominator in the faction, e.g. 2 in 2/3.
	Numerator uint64 `json:"numerator"`
	// The value by which the numerator is divided, e.g. 3 in 2/3.
	Denominator uint64 `json:"denominator"`
}

func (fr Fraction) String() string {
	return fmt.Sprintf("%d/%d", fr.Numerator, fr.Denominator)
}

// ParseFractions takes the string of a fraction as input i.e "2/3" and converts this
// to the equivalent fraction else returns an error. The format of the string must be
// one number followed by a slash (/) and then the other number.
func ParseFraction(f string) (Fraction, error) {
	o := strings.Split(f, "/")
	if len(o) != 2 {
		return Fraction{}, errors.New("incorrect formating: should have a single slash i.e. \"1/3\"")
	}
	numerator, err := strconv.ParseUint(o[0], 10, 64)
	if err != nil {
		return Fraction{}, fmt.Errorf("incorrect formatting, err: %w", err)
	}

	denominator, err := strconv.ParseUint(o[1], 10, 64)
	if err != nil {
		return Fraction{}, fmt.Errorf("incorrect formatting, err: %w", err)
	}
	if denominator == 0 {
		return Fraction{}, errors.New("denominator can't be 0")
	}
	if numerator > math.MaxInt64 || denominator > math.MaxInt64 {
		return Fraction{}, fmt.Errorf("value overflow, numerator and denominator must be less than %d", int64(math.MaxInt64))
	}
	return Fraction{Numerator: numerator, Denominator: denominator}, nil
}
