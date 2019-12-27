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
)

// Transform the image with the given parameters.
// rot(ation) must be specified in degrees.
func Transform(img image.Image, scale float64, rot float64) image.Image {
	rad := rot / 180 * math.Pi
	img2 := image.NewRGBA64(calcBounds(img.Bounds(), scale, rad))

	// rotate around centre of image
	// TODO: verify behaviour when 'Dx != Max.X', etc
	ox := img.Bounds().Dx() / 2
	oy := img.Bounds().Dy() / 2
	ox2 := img2.Bounds().Dx() / 2
	oy2 := img2.Bounds().Dy() / 2

	// iterate over new image rather than source; easier scaling and rounding
	for x2 := -ox2; x2 < ox2; x2 += 1 {
		for y2 := -oy2; y2 < oy2; y2 += 1 {
			x1, y1 := calcSrc(x2, y2, scale, rad)
			img2.Set(ox2+x2, oy2+y2, img.At(ox+mathext.Round(x1), oy+mathext.Round(y1)))
		}
	}

	return img2
}

func calcBounds(r image.Rectangle, scale float64, rad float64) image.Rectangle {
	dX := float64(r.Dx())
	dY := float64(r.Dy())
	dX2 := (dX*math.Cos(rad) + dY*math.Sin(rad)) * scale
	dY2 := (dY*math.Cos(rad) + dX*math.Sin(rad)) * scale
	return image.Rect(0, 0, mathext.Round(math.Abs(dX2)), mathext.Round(math.Abs(dY2)))
}

func calcSrc(x, y int, scale, rad float64) (float64, float64) {
	r, arc := mathext.ToPolar(float64(x), float64(y))
	x2, y2 := mathext.ToCartesian(r/scale, arc-rad)
	return x2, y2
}
