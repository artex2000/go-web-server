package main

import (
    "image"
    "image/color"
    "image/png"
    "io"
)

func drawPng(out io.Writer) {
    const width, height = 800, 600
    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.NRGBA {
                            R: uint8((x + y) & 255),
                            G: uint8((x + y) << 1 & 255),
                            B: uint8((x + y) << 2 & 255),
                            A: 255,
                        })
        }
    }
    png.Encode(out, img)
}
