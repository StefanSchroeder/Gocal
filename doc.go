/*
 * Copyright (c) 2014 Stefan Schroeder
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file
 */

/*
Package gocal implements a PDF calendar generator. It consists of a
library gocal and a standalone tool gocalendar.

* Inspired by pcal

* Simplicity: Create a nice calendar with minimum effort.

* Week number

* Moonphase

* Add events from configuration file

* Set papersize

* Choose fonts (limited)

* Several languages

* Day of year

* Background image

* Photo calendar

* More

For build instructions, test instructions, examples, see
README.md.

EXAMPLE

	package main
	import (
		"github.com/StefanSchroeder/gocal"
	)
	func main() {
		g := gocal.New(1,12,2010)
		g.CreateCalendar("test-example01.pdf")
	}


*/
package gocal
