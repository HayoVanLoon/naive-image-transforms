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
package main

import (
	"github.com/HayoVanLoon/naive-image-transforms"
	"image/jpeg"
	"log"
	"os"
)

type message struct {
	scale    float64
	rotation float64
}

func handle(m message) {
	f, err := os.Open("tmp/in.jpg")
	if err != nil {
		log.Fatal(err)
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	img2 := naive.Transform(img, m.scale, m.rotation)

	w, err := os.Create("tmp/out.jpg")
	err = jpeg.Encode(w, img2, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	handle(message{1, 15})
}
