/*
 * Copyright 2019 Hayo van Loon
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package naive

import (
	"github.com/HayoVanLoon/go-commons/mathext"
	"image"
	"math"
	"testing"
)

func TestRound(t *testing.T) {
	cases := []struct {
		msg      string
		input    float64
		expected int
	}{
		{"failed", 0, 0},
		{"failed", .1, 0},
		{"failed", .49, 0},
		{"failed", .5, 1},
		{"failed", .51, 1},
		{"failed", .9, 1},
		{"failed", 100.9, 101},
		{"failed", -.1, 0},
		{"failed", -.49, 0},
		{"failed", -.5, -1},
		{"failed", -.51, -1},
		{"failed", -.9, -1},
		{"failed", -100.9, -101},
	}
	for i, c := range cases {
		if actual := mathext.Round(c.input); actual != c.expected {
			t.Errorf("case %v: %s (expected: %v, got: %v)", i, c.msg, c.expected, actual)
		}
	}
}

func TestCalcBounds(t *testing.T) {
	sqr := image.Rect(0, 0, 10, 10)
	rect := image.Rect(0, 0, 20, 10)
	cases := []struct {
		msg      string
		input    image.Rectangle
		scale    float64
		rot      float64
		expected image.Rectangle
	}{
		{"no-op", rect, 1, 0, rect},
		{"45-sqr", sqr, 1, math.Pi * .25,
			image.Rect(0, 0,
				mathext.Round(mathext.Eucl(10, 10)),
				mathext.Round(mathext.Eucl(10, 10)))},
		{"45", rect, 1, math.Pi * .25,
			image.Rect(0, 0,
				mathext.Round(20*math.Sin(.25*math.Pi)+10*math.Cos(.25*math.Pi)),
				mathext.Round(20*math.Cos(.25*math.Pi)+10*math.Sin(.25*math.Pi)))},
		{"90", rect, 1, math.Pi * .5, image.Rect(0, 0, 10, 20)},
		{"180", rect, 1, math.Pi, rect},
		{"270", rect, 1, math.Pi * 1.5, image.Rect(0, 0, 10, 20)},
		{"360", rect, 1, math.Pi * 2, rect},
	}
	for i, c := range cases {
		if actual := calcBounds(c.input, c.scale, c.rot); actual != c.expected {
			t.Errorf("case %v: %s (expected: %v, got: %v)", i, c.msg, c.expected, actual)
		}
	}
}

func diff(a, b, d float64) bool {
	return b < a - d || a + d < b
}

func TestCalcSrc(t *testing.T) {
	cases := []struct {
		msg      string
		x1, y1 int
		scale, rad float64
		x2, y2 float64
	}{
		{"no-op", 10, 0, 1, 0, 10, 0},
		{"no-op2", -10, 0, 1, 0, -10, 0},
		{"scale-", 10, 0, .5, 0, 20, 0},
		{"scale+", 10, 0, 2, 0, 5, 0},
		{"scale-", 10, 10, .5, 0, 20, 20},
		{"scale+", 10, 10, 2, 0, 5, 5},
		{"-90", 10, 0, 1, -math.Pi * .5, 0, 10},
		{"+90", 10, 0, 1, math.Pi * .5, 0, -10},
		{"+45", mathext.Round(mathext.Eucl(10, 10)), 0, 1, math.Pi * .25, 10, -10},
		{"scale-,+45", mathext.Round(mathext.Eucl(10, 10)), 0, .5, math.Pi * .25, 20, -20},
		{"scale+,+45", mathext.Round(mathext.Eucl(10, 10)), 0, 2, math.Pi * .25, 5, -5},
	}
	for i, c := range cases {
		if x, y := calcSrc(c.x1, c.y1, c.scale, c.rad); diff(x, c.x2, .25) || diff(y, c.y2, .25) {
			t.Errorf("case %v: %s (expected: (%v,%v), got: (%v,%v))", i, c.msg, c.x2, c.y2, x, y)
		}
	}
}
