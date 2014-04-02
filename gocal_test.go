// Copyright (c) 2014 Stefan Schroeder, NY, 2014-03-10
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package gocal_test

import (
	"gocal"
	"os"
	"testing"
)

func Test_Example01(t *testing.T) {
  g := gocal.New(1,12,2010)
  g.CreateCalendar("test-example01.pdf")
}

func Test_Example02(t *testing.T) {
  g := gocal.New(1,1,2011)
  g.SetNocolor()
  g.CreateCalendar("test-example02.pdf")
}

func Test_Example03(t *testing.T) {
  g := gocal.New(1,1,2015)
  g.SetOrientation("P")
  g.SetLocale("fr_FR")
  g.SetFont("sans")
  g.CreateCalendar("test-example03.pdf")
}

func Test_Example04(t *testing.T) {
  g := gocal.New(1,1,2015)
  g.SetOrientation("P")
  g.SetPhotos("gocalendar" + string(os.PathSeparator) + "pics")
  g.CreateCalendar("test-example04.pdf")
}

    /*
$E -o example01.pdf -p P -photos pics 1 2014
$E -o example04.pdf -lang de_DE -font mono 2 2014
$E -o example05.pdf -lang nl_NL -plain 3 2014
$E -o example06.pdf -font c:/windows/Fonts/cabalett.ttf -lang en_US 4 2014
$E -o example07.pdf -p P -lang fr_FR -photo pics/taxi.JPG  4 2014
$E -o example09.pdf -lang fi_FI -font serif -p L  4 2014
$E -o example10.pdf -lang fi_FI -font mono -p L 12 2013
$E -o example11.pdf -lang de_DE -font sans -p L 6 2014
$E -o example13.pdf -font sans -noother 7 2014
$E -o example14.pdf -small 2 2014
     */
