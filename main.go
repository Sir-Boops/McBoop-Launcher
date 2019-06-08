package main

import "os"
import "fmt"
import "strings"
import "os/user"
import "runtime"
import "os/exec"
import "io/ioutil"
import "path/filepath"
import "github.com/mholt/archiver"

func main() {

  // What this does
  // 0: Make sure this is running on Linux
  // 0.5: Check for launcher updates
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

  // Check for launcher updates
  expath, _ := filepath.Abs(os.Args[0])
  ex_remote_sum := ReadRemoteText("https://s3.amazonaws.com/boops-deploy/McBoop/McBoop.sha256")
  ex_sum := Sha256SumFile(expath)
  if ex_sum != ex_remote_sum {
    fmt.Println("A launcher update has been found and will now be installed")
    fmt.Println("This can take a little bit of time!")
    os.Remove(expath)
    DownloadFile("https://s3.amazonaws.com/boops-deploy/McBoop/McBoop", expath)
    os.Chmod(expath, os.ModePerm)
    fmt.Println("Done! please re-run your launch command again!")
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

  // Check Java version
  // If java.tar.gz is not there or sum is bad redownload it
  // Yes this is where the java version is set
  current := "96d24d94c022b3e414b612cae8829244329d71ad2cce790f099c020f33247e7e"
  url := "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u212-b04/OpenJDK8U-jre_x64_linux_hotspot_8u212b04.tar.gz"

  // Make sure the java.tar.gz file is there
  if _, err := os.Stat(homedir + ".mcboop/java.tar.gz"); err != nil {
    // New install download current Java
    os.RemoveAll(homedir + ".mcboop/java")
    fmt.Println("Downloading Java.tar.gz please wait this could take a bit")
    DownloadFile(url, homedir + ".mcboop/java.tar.gz")
  } else {
    // It is there to check sum
    sha := Sha256SumFile(homedir + ".mcboop/java.tar.gz")
    if sha != current {
      // We got a bad java.tar.gz
      os.RemoveAll(homedir + ".mcboop/java")
      os.RemoveAll(homedir + ".mcboop/java.tar.gz")
      fmt.Println("Downloading Java.tar.gz please wait this could take a bit")
      DownloadFile(url, homedir + ".mcboop/java.tar.gz")
    }
  }

  // Extract the java.tar.gz file if it's not already
  if _, err := os.Stat(homedir + ".mcboop/java"); err != nil {
    // We extracting bois
    archiver.Unarchive(homedir + ".mcboop/java.tar.gz", homedir + ".mcboop")
    files, _ := filepath.Glob(homedir + ".mcboop/jdk*")
    os.Rename(files[0], homedir + ".mcboop/java")
  }

  // Finally download McBoop
  // If it's not already there else check for an update
  if _, err := os.Stat(homedir + ".mcboop/McBoop.jar"); err != nil {
    DownloadFile("https://s3.amazonaws.com/boops-deploy/McBoop/McBoop.jar", homedir + ".mcboop/McBoop.jar")
  } else {
    // ShaSUM it and check for updates
    remote_sum := ReadRemoteText("https://s3.amazonaws.com/boops-deploy/McBoop/McBoop.jar.sha256")
    sum := Sha256SumFile(homedir + ".mcboop/McBoop.jar")
    if sum != remote_sum {
      fmt.Println("A newer version of McBoop has been found and will now be installed")
      fmt.Println("This should only take a moment to download")
      os.Remove(homedir + ".mcboop/McBoop.jar")
      DownloadFile("https://s3.amazonaws.com/boops-deploy/McBoop/McBoop.jar", homedir + ".mcboop/McBoop.jar")
    }
  }

  mcboop_launch_cmd := []string{"-jar", homedir + ".mcboop/McBoop.jar"}

  // Gen McBoop command
  for i := 1; i < len(os.Args); i++ {
    mcboop_launch_cmd = append(mcboop_launch_cmd, os.Args[i])
  }

  // Launch MC from the command McBoop.jar Generated
  // TODO: FInd a better way to handle this
  cmd := exec.Command(homedir + ".mcboop/java/bin/java", mcboop_launch_cmd...)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run()

  // Make sure the .launch file is there
  if _, err := os.Stat(homedir + ".mcboop/.launch"); err != nil {
    fmt.Println("McBoop.jar didn't exit right check above to see what happened")
  } else {
    // We have the file!
    launch_file, _ := os.Open(homedir + ".mcboop/.launch")
    launch_cmd, _ := ioutil.ReadAll(launch_file)
    launch_file.Close()
    os.Remove(homedir + ".mcboop/.launch")

    mc := exec.Command(homedir + ".mcboop/java/bin/java", strings.Split(string(launch_cmd), " ")...)
    mc.Stdout = os.Stdout
    mc.Stderr = os.Stderr
    fmt.Println("")
    fmt.Println("Game logging starts here")
    fmt.Println("")
    mc.Run()
  }
}
