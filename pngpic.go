package main

import (
    "image"
    "image/color"
    "image/png"
    "os"
    "io"
    "encoding/binary"
    "log"
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

func drawPng(out io.Writer) {
    const width, height = 800, 600
    img := image.NewNRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.NRGBA {
                            R: 255,
                            G: 255,
                            B: 255,
                            A: 255,
                        })
        }
    }

    off_x := int((width - glyph.width) >> 1)
    off_y := int((height - glyph.height) >> 1)
    idx := 0
    for y := 0; y < int(glyph.height); y++ {
        for x := 0; x < int(glyph.width); x++ {
            img.Set(x+off_x, y+off_y, color.NRGBA {
                            R: 0,
                            G: 0,
                            B: 0,
                            A: glyph.data[idx],
                        })
            idx++
        }
    }

    png.Encode(out, img)
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



