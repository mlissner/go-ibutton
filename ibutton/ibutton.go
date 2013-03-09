// This file is part of ibutton.
//
// Copyright (C) 2013 Max Hille <mh@lambdasoup.com>
//
// ibutton is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// at your option) any later version.
//
// ibutton is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ibutton.  If not, see <http://www.gnu.org/licenses/>.

// Package main includes the runnable commands for the ibutton project
package main

import (
	"flag"
	"fmt"
	"github.com/maxhille/go-ibutton/w1"
	"os"
)

// parse arguments
var command = flag.String("command", "help", "displays general help")

func main() {

	flag.Parse()

	switch *command {
	case "status":
		button := new(w1.Button)
		err := button.Open()
		defer button.Close()
		if err != nil {
			fmt.Printf("could not open iButton (%v)\n", err)
			os.Exit(1)
		}
		status, err := button.Status()
		if err != nil {
			fmt.Printf("could not get iButton status (%v)\n", err)
			os.Exit(1)
		}
		fmt.Printf("time:           %v\n", status.Time())
		fmt.Printf("model:          %v\n", status.Name())
		fmt.Printf("timestamp:      %v\n", status.MissionTimestamp())
		fmt.Printf("count:          %v\n", status.SampleCount())
		fmt.Printf("running:        %v\n", status.MissionInProgress())
		fmt.Printf("memory cleared: %v\n", status.MemoryCleared())
		fmt.Printf("resolution:     %v\n", func() string {
			if status.HighResolution() {
				return "0.0625°C"
			}
			return "0.5°C"
		}())
		fmt.Printf("rate:           %v\n", status.SampleRate())
	case "clear":
		button := new(w1.Button)
		err := button.Open()
		defer button.Close()
		if err != nil {
			fmt.Printf("could not open button (%v)\n", err)
			os.Exit(1)
		}
		err = button.ClearMemory()
		if err != nil {
			fmt.Printf("could not clear memory (%v)\n", err)
			os.Exit(1)
		}
		fmt.Printf("Cleared Memory.\n")
	case "read":
		button := new(w1.Button)
		err := button.Open()
		defer button.Close()
		if err != nil {
			fmt.Printf("could not open button (%v)\n", err)
			os.Exit(1)
		}
		samples, err := button.ReadLog()
		if err != nil {
			fmt.Printf("could not read log (%v)\n", err)
			os.Exit(1)
		}
		for _, sample := range samples {
			fmt.Printf("%v\t%3.3f°C\n", sample.Time, sample.Temp)
		}
	case "stop":
		button := new(w1.Button)
		err := button.Open()
		defer button.Close()
		if err != nil {
			fmt.Printf("could not open button (%v)\n", err)
			os.Exit(1)
		}
		err = button.StopMission()
		if err != nil {
			fmt.Printf("could not stop mission (%v)\n", err)
			os.Exit(1)
		}
		fmt.Printf("Stopped mission.\n")
	case "help":
		flag.Usage()
		os.Exit(2)
	}

}
