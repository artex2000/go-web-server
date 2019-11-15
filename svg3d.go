package main

import (
    "fmt"
    "math"
    "io"
    "log"
)

const (
    width, height = 1600, 1000
    cells = 100
    angle = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)
var minf, maxf float64

type f3d func(float64, float64) float64
type polygon struct {
    poly [8]float64
    color [3] uint8
}

type plot3d struct {
    formula f3d
    xyrange float64
    zfactor float64
}

var plot3dlist = map[string]*plot3d {
    "sin" : &plot3d{ sin3d, 30.0, 0.4 },
    "hyp" : &plot3d{ hyp3d, 6.0, 0.15 },
    "exp" : &plot3d{ exp3d, 3.0, 0.3 },
    "log" : &plot3d{ log3d, 8.0, 0.01 },
    "bum" : &plot3d{ bum3d, 2.0, 1.0 },
    "cub" : &plot3d{ cub3d, 2.0, 0.7 },
    "exe" : &plot3d{ exe3d, 2.0, 0.7 },
    "exc" : &plot3d{ exc3d, 3.0, 0.0008 },
}

func drawSvg(out io.Writer, plot string) {
    var f *plot3d
    var ok bool
    minf, maxf = math.MaxFloat64, math.SmallestNonzeroFloat64

    if f, ok = plot3dlist[plot]; !ok {
        f = plot3dlist["sin"]
    }

    fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
                     "style='stroke: grey; stroke-width: 0.7' "+
                     "width='%d' height='%d'>", width, height)
    for i := 0; i < cells; i++ {
        for j := 0; j < cells; j++ {
            p := corner(i, j, f)
            fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' " +
                             "style='fill:#%02x%02x%02x'/>\n",
                             p.poly[0], p.poly[1], p.poly[2], p.poly[3],
                             p.poly[4], p.poly[5], p.poly[6], p.poly[7],
                             p.color[0], p.color[1], p.color[2])
        }
    }
    fmt.Fprintln(out, "</svg>")
    log.Printf("min %g, max %g\n", minf, maxf)
}

func corner(i, j int, f *plot3d) *polygon {
    var r polygon
    skewi := []int{ 1, 0, 0, 1 }
    skewj := []int{ 0, 0, 1, 1 }
    xyscale := width / 2 / f.xyrange
    zscale := height * f.zfactor
    for k := 0; k < 4; k++ {
        x := f.xyrange * (float64(i + skewi[k]) / cells - 0.5)
        y := f.xyrange * (float64(j + skewj[k]) / cells - 0.5)

        z := f.formula(x, y)
        if z < minf {
            minf = z
        }
        if z > maxf {
            maxf = z
        }
        /*
        if z < -0.5 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0x55, 0x00
        } else if z < -0.4 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0x00, 0x80
        } else if z < -0.3 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0x80, 0x00
        } else if z < -0.2 {
            r.color[0], r.color[1], r.color[2] = 0x80, 0x00, 0x00
        } else if z < -0.1 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0x80, 0x80
        } else if z < 0.0 {
            r.color[0], r.color[1], r.color[2] = 0x80, 0x80, 0x00
        } else if z < 0.1 {
            r.color[0], r.color[1], r.color[2] = 0x80, 0x80, 0x80
        } else if z < 0.2 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0x00, 0xff
        } else if z < 0.3 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0xff, 0x00
        } else if z < 0.4 {
            r.color[0], r.color[1], r.color[2] = 0xff, 0x00, 0x00
        } else {
            r.color[0], r.color[1], r.color[2] = 0x00, 0xff, 0xff
        }
        */
        if z > 0.8 {
            r.color[0], r.color[1], r.color[2] = 0xff, 0x00, 0x00
        } else if z > 0.6 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0xff, 0x00
        } else if z > 0.4 {
            r.color[0], r.color[1], r.color[2] = 0x00, 0x00, 0xff
        } else {
            r.color[0], r.color[1], r.color[2] = 0xff, 0xff, 0xff
        }



        r.poly[k * 2] = width / 2 + (x - y)*cos30*xyscale
        r.poly[k * 2 +1] = height / 2 + (x + y)*sin30*xyscale - z*zscale
    }
    return &r
}

func sin3d(x, y float64) float64 {
    r := math.Hypot(x, y)
    return math.Sin(r) / r
}

func hyp3d(x, y float64) float64 {
    return x*x/7 - y*y/7 
}

func exp3d(x, y float64) float64 {
    return 0.75 / math.Exp(x*x*y*y*81) 
}

func log3d(x, y float64) float64 {
    return 0.7 / (math.Log(x*x + y*y) + 0.6) 
}

func bum3d(x, y float64) float64 {
    return math.Sin(x*5) * math.Cos(y*5) / 5 
}

func cub3d(x, y float64) float64 {
    return x*y*y*y - y*x*x*x;
}

func exe3d(x, y float64) float64 {
    r := x*x + y*y
    return (x*x + 3*y*y) * math.Exp(-r) - 0.5
}

func exc3d(x, y float64) float64 {
    return math.Floor(math.Exp(math.Abs(x*y/2))) + math.Round(1/math.Cos(x*y))-2
}
