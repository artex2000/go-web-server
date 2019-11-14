package main

import (
    "fmt"
    "math"
    "io"
)

const (
    width, height = 800, 600
    cells = 100
    xyrange = 30
    xyscale = width / 2 / xyrange
    zscale = height * 0.4
    angle = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)
type f3d func(float64, float64) float64
var f3dlist = map[string]f3d {
    "sin" : sin3d,
    "hyp" : hyp3d,
}

func drawSvg(out io.Writer, plot string) {
    var f f3d
    var ok bool
    if f, ok = f3dlist[plot]; !ok {
        f = sin3d
    }

    fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
                     "style='stroke: grey; fill: white; stroke-width: 0.7' "+
                     "width='%d' height='%d'>", width, height)
    for i := 0; i < cells; i++ {
        for j := 0; j < cells; j++ {
            poly := corner(i, j, f)
            fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
                             poly[0], poly[1], poly[2], poly[3],
                             poly[4], poly[5], poly[6], poly[7])
        }
    }
    fmt.Fprintln(out, "</svg>")
}

func corner(i, j int, f f3d) *[8]float64 {
    var r [8]float64
    skewi := []int{ 1, 0, 0, 1 }
    skewj := []int{ 0, 0, 1, 1 }
    for k := 0; k < 4; k++ {
        x := xyrange * (float64(i + skewi[k]) / cells - 0.5)
        y := xyrange * (float64(j + skewj[k]) / cells - 0.5)

        z := f(x, y)

        r[k * 2] = width / 2 + (x - y)*cos30*xyscale
        r[k * 2 +1] = height / 2 + (x + y)*sin30*xyscale - z*zscale
    }
    return &r
}

func sin3d(x, y float64) float64 {
    r := math.Hypot(x, y)
    return math.Sin(r) / r
}

func hyp3d(x, y float64) float64 {
    return x*x - y*y
}
