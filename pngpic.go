package main

import (
    "image"
    "image/color"
    "image/png"
    "os"
    "io"
    "encoding/binary"
    "math"
)

type Glyph struct {
    codepoint uint32
    left      uint32
    top       uint32
    width     uint32
    height    uint32
    data      [][]byte
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
            img.Set(x+off_x, y+off_y, blend(fg, bg, glyph.data[y][x]))
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
    glyph.data = make([][]byte, glyph.height)
    for i:= 0; i < int(glyph.height); i++ {
        glyph.data[i] = make([]byte, glyph.width)
        _, err = io.ReadFull(f, glyph.data[i])
        if (err != nil) {
            panic(err)
        }
    }
    getSdf(&glyph)
}

func b2f(b byte) float64 {
    return float64(b) / 255.0
}

func f2b(f float64) byte {
    f = 0.5 - f
    if (f < 0.0) {
        f = 0.0
    } else if (f > 1.0) {
        f = 1.0
    }
    return byte(f * 255.0)
}

func getSdf(glyph *Glyph) {
    var left, right, top, bottom bool
    // These variables surround current point of interest as follows:
    // a  b  c
    // d  X  e
    // f  g  h
    var a,b,c,d,e,f,g,h float64
    var df float64
    const SQRT2 = 1.4142136

    src := glyph.data
    sdf := make([][]byte, glyph.height)
    for i:= 0; i < int(glyph.height); i++ {
        sdf[i] = make([]byte, glyph.width)
    }
    for y := 0; y < int(glyph.height); y++ {
        if (y == 0) {
            top = true
        } else if (y == int(glyph.height - 1)) {
            bottom = true
        } else {
            top, bottom = false, false
        }

        for x := 0; x < int(glyph.width); x++ {
            if (x == 0) {
                left = true
            } else if (x == int(glyph.width - 1)) {
                right = true
            } else {
                left, right = false, false
            }

            //take care of border pixels (apron simulation)
            if (top) {
                a, b, c = 0.0, 0.0, 0.0
            } else {
                a, b, c = b2f(src[y-1][x-1]), b2f(src[y-1][x]), b2f(src[y-1][x+1])
            }

            if (bottom) {
                f, g, h = 0.0, 0.0, 0.0
            } else {
                f, g, h = b2f(src[y+1][x-1]), b2f(src[y+1][x]), b2f(src[y+1][x+1])
            }

            if (left) {
                a, d, f = 0.0, 0.0, 0.0
            } else {
                a, d, f = b2f(src[y-1][x-1]), b2f(src[y][x-1]), b2f(src[y+1][x-1])
            }

            if (right) {
                c, e, h = 0.0, 0.0, 0.0
            } else {
                c, e, h = b2f(src[y-1][x+1]), b2f(src[y][x+1]), b2f(src[y+1][x+1])
            }

            //current point of interest
            xx := b2f(src[y][x])
            if (xx == 1.0) {
                sdf[y][x] = 255
                continue
            } else if (xx == 0.0) {
                var hor, vert bool
                if (b == 1.0 || g == 1.0) {
                    vert = true
                }

                if (d == 1.0 || e == 1.0) {
                    hor = true
                }

                if (!vert && !hor) {
                    sdf[y][x] = 0
                    continue
                }
            }

            gx := -a - d * SQRT2 - f + c + e * SQRT2 + h
            gy := -a - b * SQRT2 - c + f + g * SQRT2 + h

            gx, gy = math.Abs(gx), math.Abs(gy)
            if (gx < 0.0001 || gy < 0.0001) {
                df = (0.5 - xx) * SQRT2
            } else {
                glen := gx*gx + gy*gy
                glen = 1.0 / math.Sqrt(glen)

                gx *= glen
                gy *= glen

                if (gx < gy) {
                    gx, gy = gy, gx
                }

                a1 := 0.5 * gy / gx

                if (xx < a1) {
                    df = 0.5 * (gx + gy) - math.Sqrt(2.0 * gx * gy * xx)
                } else if (xx < (1.0 - a1)) {
                    df = (0.5 - xx) * gx
                } else {
                    df = -0.5 * (gx + gy) + math.Sqrt(2.0 * gx * gy * (1 - xx))
                }
            }
            df *= 1.0 / SQRT2
            sdf[y][x] = f2b(df)
        }
    }
    glyph.data = sdf
}

            



