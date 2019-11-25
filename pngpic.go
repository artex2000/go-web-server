package main

import (
    "image"
    "image/color"
    "image/png"
    "os"
    "io"
    "encoding/binary"
    "log"
    "math"
)

type Glyph struct {
    codepoint uint32
    left      uint32
    top       uint32
    width     uint32
    height    uint32
    data      []byte
}

var glyph Glyph
var bg = color.NRGBA { R: 30, G: 50, B: 100, A: 255 }
var fg = color.NRGBA { R: 70, G: 90, B: 200, A: 255 }

func drawPng(out io.Writer) {
    const width, height = 800, 600
    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, bg)
        }
    }

    off_x := int((width - glyph.width) >> 1)
    off_y := int((height - glyph.height) >> 1)
    idx := 0
    for y := 0; y < int(glyph.height); y++ {
        for x := 0; x < int(glyph.width); x++ {
            img.Set(x+off_x, y+off_y, blend(fg, bg, glyph.data[idx]))
            idx++
        }
    }

    png.Encode(out, img)
}

func rgb2linear(c uint8) float64 {
    inv255 := 1.0 / 255.0
    r := float64(c) * inv255
    return r * r
}

func linear2rgb(f float64) uint8 {
    r := math.Sqrt(f)
    r *= 255.0
    return uint8(r + 0.5)
}


func blend(src, dst color.NRGBA, alpha uint8) color.NRGBA {
    inv255 := 1.0 / 255.0
    sr := rgb2linear(src.R)
    sg := rgb2linear(src.G)
    sb := rgb2linear(src.B)
    dr := rgb2linear(dst.R)
    dg := rgb2linear(dst.G)
    db := rgb2linear(dst.B)
    a := float64(alpha) * inv255

    rr := sr * a + dr * (1 - a)
    rg := sg * a + dg * (1 - a)
    rb := sb * a + db * (1 - a)

    return color.NRGBA {
        R: linear2rgb(rr),
        G: linear2rgb(rg),
        B: linear2rgb(rb),
        A: 255 }
}

func init() {
    f, err := os.Open("./hack.bin")
    if (err != nil) {
        panic(err)
    }
    tmp := make([]byte, 20)
    _, err = io.ReadFull(f, tmp[:])
    if (err != nil) {
        panic(err)
    }
    glyph.codepoint = binary.LittleEndian.Uint32(tmp[0:4])
    glyph.left = binary.LittleEndian.Uint32(tmp[4:8])
    glyph.top = binary.LittleEndian.Uint32(tmp[8:12])
    glyph.width = binary.LittleEndian.Uint32(tmp[12:16])
    glyph.height = binary.LittleEndian.Uint32(tmp[16:20])
    glyph.data = make([]byte, glyph.width * glyph.height)
    _, err = io.ReadFull(f, glyph.data)
    if (err != nil) {
        panic(err)
    }
    log.Printf("width %g, height %g\n", glyph.width, glyph.height)
}



