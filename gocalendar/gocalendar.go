// Copyright (c) 2014 Stefan Schroeder, NY, 2014-03-10
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"flag"
	"github.com/StefanSchroeder/gocal"
	"strconv"
	"time"
)

var optFont = flag.String("font", "serif", "Font")
var optFontScale = flag.Float64("fontscale", 1.0, "Font")
var optFooter = flag.String("footer", "Gocal", "Footer note")
var optHideDOY = flag.Bool("nodoy", false, "Hide day of year (false)")
var optPlain = flag.Bool("plain", false, "Hide everything")
var optHideEvents = flag.Bool("noevents", false, "Hide events from config file (false)")
var optHideMoon = flag.Bool("nomoon", false, "Hide moon phases (false)")
var optHideWeek = flag.Bool("noweek", false, "Hide week number (false)")
var optLocale = flag.String("lang", "", "Language")
var optOrientation = flag.String("p", "P", "Orientation (L)andscape/(P)ortrait")
var optPaper = flag.String("paper", "A4", "Paper format (A3 A4 A5 Letter Legal)")
var optPhoto = flag.String("photo", "", "Show photo (single image PNG JPG GIF)")
var optConfig = flag.String("config", "gocal.xml", "Configuration file")
var optPhotos = flag.String("photos", "", "Show photos (directory PNG JPG GIF)")
var optWallpaper = flag.String("wall", "", "Show wallpaper PNG JPG GIF")
var outfilename = flag.String("o", "output.pdf", "Output filename")
var optSmall = flag.Bool("small", false, "Smaller fonts")
var optHideOtherMonths = flag.Bool("noother", false, "Hide neighboring month days")
var optNocolor = flag.Bool("nocolor", false, "Sundays and Saturdays in black, instead of red.")

func main() {
	flag.Parse()

	wantyear := int(time.Now().Year())
	beginmonth := 1
	endmonth := 12

	if flag.NArg() == 1 {
		dummyyear, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		wantyear = int(dummyyear)
	} else if flag.NArg() == 2 {
		dummymonth, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		dummyyear, _ := strconv.ParseInt(flag.Arg(1), 10, 32)
		beginmonth = int(dummymonth)
		endmonth = int(dummymonth)
		wantyear = int(dummyyear)
	} else if flag.NArg() == 3 {
		dummymonthBegin, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		dummymonthEnd, _ := strconv.ParseInt(flag.Arg(1), 10, 32)
		dummyyear, _ := strconv.ParseInt(flag.Arg(2), 10, 32)
		beginmonth = int(dummymonthBegin)
		endmonth = int(dummymonthEnd)
		wantyear = int(dummyyear)
	}

	g := gocal.New(beginmonth, endmonth, wantyear)
	g.SetFont(*optFont)
	g.SetOrientation(*optOrientation)
	g.SetPaperformat(*optPaper)
	g.SetLocale(*optLocale)
	g.SetConfig(*optConfig)
	if *optPlain == true {
		g.SetPlain()
	}
	if *optHideDOY == true {
		g.SetHideDOY()
	}
	if *optHideWeek == true {
		g.SetHideWeek()
	}
	if *optHideMoon == true {
		g.SetHideMoon()
	}
	if *optSmall == true {
		g.SetSmall()
	}
	if *optNocolor == true {
		g.SetNocolor()
	}
	g.SetFontScale(*optFontScale)
	g.SetWallpaper(*optWallpaper)
	g.SetPhotos(*optPhotos)
	g.SetPhoto(*optPhoto)
	g.SetFooter(*optFooter)
	/*
	  g.AddEvent(31, 1, "HALLO", "")
	  g.AddEvent(28, 2, "HALLO", "")
	  g.AddEvent(31, 3, "HALLO", "")
	  g.AddEvent(30, 4, "HALLO", "")
	*/
	g.CreateCalendar(*outfilename)
}
