// Copyright (c) 2014 Stefan Schroeder, NY, 2014-03-10
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package gocal_test

import (
	"gocal"
	"os"
	"runtime"
	"testing"
)

func Test_Example01(t *testing.T) {
  g := gocal.New(1,12,2010)
  g.CreateCalendar("test-example01.pdf")
}

func Test_Example02(t *testing.T) {
  g := gocal.New(1,1,2011)
  g.SetNocolor()
  g.SetOrientation("L")
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

func Test_Example05(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetOrientation("L")
  g.SetFont("mono")
  g.SetLocale("de_DE")
  g.CreateCalendar("test-example05.pdf")
}

func Test_Example06(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetOrientation("P")
  g.SetPlain()
  g.SetLocale("nl_NL")
  g.CreateCalendar("test-example06.pdf")
}

func Test_Example07(t *testing.T) {
  g := gocal.New(3,4,2013)
  if runtime.GOOS == "windows" {
    g.SetFont("c:\\windows\\Fonts\\cabalett.ttf")
  }
  g.SetFooter("Windows specific Font inclusion example")
  g.CreateCalendar("test-example07.pdf")
}

func Test_Example08(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetPhoto("gocalendar" + string(os.PathSeparator) + "pics" + string(os.PathSeparator) + "taxi.JPG")
  g.SetOrientation("P")
  g.CreateCalendar("test-example08.pdf")
}

func Test_Example09(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetLocale("fi_FI")
  g.CreateCalendar("test-example09.pdf")
}

func Test_Example10(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetFontScale(0.5)
  g.CreateCalendar("test-example10.pdf")
}

func Test_Example11(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetSmall()
  g.CreateCalendar("test-example11.pdf")
}

func Test_Example12(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetHideOtherMonth()
  g.CreateCalendar("test-example12.pdf")
}

func Test_Example13(t *testing.T) {
  g := gocal.New(3,4,2013)
  g.SetHideWeek()
  g.AddEvent(15, 3, "One Event", "")
  g.AddEvent(16, 3, "Another Event", "")
  g.AddEvent(17, 4, "Hegdehog\\nvisits", "")
  g.AddEvent(18, 4, "Day\\nof the\\nDay", "")
  g.SetHideMoon()
  g.CreateCalendar("test-example13.pdf")
}

func Test_Example14(t *testing.T) {
  g := gocal.New(1,12,2013)
  g.SetConfig("test-gocal.xml")
  g.CreateCalendar("test-example14.pdf")
}

