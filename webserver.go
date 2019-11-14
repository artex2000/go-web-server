package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
)

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/gif", gifer)
    http.HandleFunc("/svg", svger)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
    for k, v := range r.Header {
        fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
    }
    fmt.Fprintf(w, "Host = %q\n", r.Host)
    fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    for k, v := range r.Form {
        fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
    }
}

func gifer(w http.ResponseWriter, r *http.Request) {
    size := 100
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    if v, ok := r.Form["size"]; !ok {
        log.Print("Size not passed")
    } else {
        var err error
        if size, err = strconv.Atoi(v[0]); err != nil {
            log.Print(err)
            size = 100
        }
    }
    lissajous(w, size)
}
    
func svger(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "image/svg+xml")
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    if v, ok := r.Form["f"]; !ok {
        log.Print("Size not passed")
        drawSvg(w, "")
    } else {
        drawSvg(w, v[0])
    }
}
