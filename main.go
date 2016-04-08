package main

import (
  "os"
  "fmt"
  "log"
  "errors"
  "strconv"
  "strings"
  "io/ioutil"
  "net/http"
  "net/url"
  "github.com/daddye/vips"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    log.Fatal("$PORT must be set")
  }

  http.HandleFunc("/resize", resizeHandler)
  http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}


func resizeHandler(w http.ResponseWriter, r *http.Request) {
  log.Printf("New Request: %s", r.URL)

  uri := r.URL.Query().Get("src")
  size := r.URL.Query().Get("size")
  mode := r.URL.Query().Get("mode")

  if !strings.Contains(size, "x") {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "size should be of format (width)x(height), e.g. 200x100")
    return
  }

  splited := strings.Split(size, "x")
  width, _ := strconv.Atoi(splited[0])
  height, _ := strconv.Atoi(splited[1])

  crop := false
  embed := false
  if mode == "crop" {
    crop = true
  } else {
    embed = true
  }

  data, err := download(uri)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Failed to download image from %s (%s)", uri, err)
    return
  }

  resized, _ := resize(data, width, height, crop, embed)

  w.Header().Set("Content-Type", "image/jpeg")
  w.Header().Set("Content-Length", strconv.Itoa(len(resized)))
  if _, err := w.Write(resized); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    log.Printf("Unable to write image response")
  }
}

func download(uri string) (data []byte, err error) {
  _, err = url.Parse(uri)
  if err != nil {
    log.Printf("Malformatted Image URL %s: %s", uri, err)
    return
  }

  res, err := http.Get(uri)
  if err != nil {
    log.Printf("Failed to download image from %s: %s", uri, err)
    return
  } else if res.StatusCode != http.StatusOK {
    log.Printf("Cannot find image on %s:", uri)
    err = errors.New(fmt.Sprintf("Cannot find image with status %s (%s)", strconv.Itoa(res.StatusCode), http.StatusText(res.StatusCode)))
    return
  }

  data, err = ioutil.ReadAll(res.Body)
  return
}

func resize(original []byte, w, h int, crop, embed bool) (resized []byte, err error) {
  options := vips.Options{
    Width: w,
    Height: h,
    Crop: crop,
    Embed: embed,
    Extend: vips.EXTEND_WHITE,
    Interpolator: vips.BILINEAR,
    Gravity: vips.CENTRE,
    Quality: 82,
  }

  resized, err = vips.Resize(original, options)
  return
}
