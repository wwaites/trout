/*
   Copyright (C) 2015 William Waites

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
   trout - a program that reactively diagnoses network faults
*/
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
		fmt.Fprintf(os.Stderr, `
trout -- Copyright (C) 2015 William Waites.
This program comes with ABSOLUTELY NO WARRANTY. This is free software, and
you are welcome to redistribute it under the terms of the GNU General Public
License version 3 or later.
`)
		fmt.Fprintf(os.Stderr, "\nUsage: %s [flags] hostname\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
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

func Alive(host string) (alive bool) {
	cmd := exec.Command(fping, "-q", "-r", "0", "-c", "1", host)
	err := cmd.Run()
	alive = err == nil
	return
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
