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
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"fmt"
	"github.com/HayoVanLoon/go-commons/logjson"
	"github.com/HayoVanLoon/naive-image-transforms"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strings"
)

type message struct {
	srcUrl  string
	scale   float64
	rotate  float64
	destUrl string
}

func _close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err.Error())
	}
}

func handle(m message) error {
	ctx := context.Background()

	r, err := getReader(ctx, m.srcUrl)
	if err != nil {
		return nil
	}
	defer _close(r)
	img, err := jpeg.Decode(r)
	if err != nil {
		return err
	}

	img2 := naive.Transform(img, m.scale, m.rotate)

	w, err := getWriter(ctx, m.destUrl)
	if err != nil {
		return err
	}
	defer _close(w)
	err = jpeg.Encode(w, img2, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}

	return nil
}

func splitGsUrl(url string) (string, string, error) {
	rem := url[5:]
	b := strings.Index(rem, "/")
	if b <= 0 {
		return "", "", fmt.Errorf("malformed Google Storage url %s", url)
	}
	if len(rem) == b {
		return "", "", fmt.Errorf("malformed Google Storage url %s", url)
	}
	return rem[:b], rem[b:], nil
}

func getReader(ctx context.Context, src string) (io.ReadCloser, error) {
	if src[:5] == "gs://" {
		b, o, err := splitGsUrl(src)
		if err != nil {
			return nil, err
		}
		cl, err := storage.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		r, err := cl.Bucket(b).Object(o).NewReader(ctx)
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		f, err := os.Open(src)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

func getWriter(ctx context.Context, dest string) (io.WriteCloser, error) {
	if dest[:5] == "gs://" {
		b, o, err := splitGsUrl(dest)
		if err != nil {
			return nil, err
		}
		cl, err := storage.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		w := cl.Bucket(b).Object(o).NewWriter(ctx)
		return w, nil
	} else {
		f, err := os.Create(dest)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

func main() {
	src := flag.String("src", "", "Source image file")
	dest := flag.String("dest", "", "Source image file")
	scale := flag.Float64("scale", 1, "Scale by x")
	rotate := flag.Float64("rotate", 0, "Rotate by x degrees")

	if *src == "" {
		log.Fatalln("no source file specified")
	}
	if *dest == "" {
		log.Fatalln("no destination file specified")
	}

	m := message{
		srcUrl:  *src,
		scale:   *scale,
		rotate:  *rotate,
		destUrl: *dest,
	}

	err := handle(m)
	if err != nil {
		logjson.Critical(err.Error())
	}
}
