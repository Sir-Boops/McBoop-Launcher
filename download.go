package main

import "io"
import "os"
import "net/http"

func DownloadFile(URL string, DL_PATH string) {

  // Download the new java.tar.gz
  resp, _ := http.Get(URL)
  defer resp.Body.Close()
  out, _ := os.Create(DL_PATH)
  defer out.Close()
  io.Copy(out, resp.Body)

}
