package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var fping string
var traceroute string
var numeric bool
func init() {
	flag.BoolVar(&numeric, "n", false, "Do not resolve addresses")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] hostname\n", os.Args[0])
		flag.PrintDefaults()
	}

	var err error
	fping, err = exec.LookPath("fping")
	if err != nil {
		panic(err)
	}
	
	traceroute, err = exec.LookPath("traceroute")
	if err != nil {
		panic(err)
	}
}

func Alive(host string) bool {
	cmd := exec.Command(fping, "-q", "-r", "0", "-c", "1", host)
	err := cmd.Run()
	if err == nil {
		return true
	} else {
		return false
	}
}

func Trace(host string) ([]byte, error) {
	var cmd *exec.Cmd
	if numeric {
		cmd = exec.Command(traceroute, "-n", host)
	} else {
		cmd = exec.Command(traceroute, host)
	}
	return cmd.CombinedOutput()
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(255)
	}

	hostname := flag.Arg(0)

	c := time.Tick(1 * time.Second)
	for _ = range c {
		if !Alive(hostname) {
			out, _ := Trace(hostname)
			fmt.Printf("%s --> %s\n%s\n", time.Now(), hostname, out)
			time.Sleep(10 * time.Second)
		}
	}
}