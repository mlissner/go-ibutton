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

package main

import (
	w1 "./w1"
	"fmt"
	"flag"
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
			if err != nil {
				fmt.Printf("could not open iButton (%v)\n", err)
				os.Exit(1)
			}
			status, err := button.Status()
			if err != nil {
				fmt.Printf("could not get iButton status (%v)\n", err)
				os.Exit(1)
			}
			fmt.Printf("time: %v\n", status.Time())
			button.Close()
		case "help":
			flag.Usage()
			os.Exit(2)
	}

}