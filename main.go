package main

import "os"
import "fmt"
import "os/user"
import "runtime"

func main() {
  // TODO: Auto update this boi

  // What this does
  // 0: Make sure this is running on Linux
  // 1: Creates a ~/.McBoop folder
  // 2: Downloads latest good openjdk version
  // ( Checking the meta.json file at https://git.sergal.org/Sir-Boops/McBoop-Launcher )
  // 3: Downloads and puts McBoop.jar into ~/.mcboop
  // 4: IF McBoop.jar was already present then just check for updates for it and auto update

  // Check current running OS
  if runtime.GOOS != "linux" {
    fmt.Println("Sorry, this tool is only for Linux at the moment")
    os.Exit(0)
  }

  // Get Users home dir path
  usr, _ := user.Current()
  homedir := usr.HomeDir + "/"

  // Check if ~/.mcboop is there or not
  // and if it's not create it
  if _, err := os.Stat(homedir + ".mcboop"); err != nil {
    // New Install!
    os.MkdirAll(homedir + ".mcboop", os.ModePerm)
  }

  // TODO: Check Meta File
}
