package main

import "os"
import "io"
import "encoding/hex"
import "crypto/sha256"

func Sha256SumFile(FILE string) (string) {
  f, _ := os.Open(FILE)
  defer f.Close()
  h := sha256.New()
  io.Copy(h, f)
  return hex.EncodeToString(h.Sum(nil))
}
