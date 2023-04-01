// Copyright (c) 2014 Stefan Schroeder, NY, 2014-04-13
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"flag"
	"fmt"
	"github.com/StefanSchroeder/Gocal"
	"os"
	"strconv"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myFlags arrayFlags

const VERSION = "0.9 the Unready"

var optFont = flag.String("font", "serif", "Font")
var optFontScale = flag.Float64("fontscale", 1.0, "Font")
var optYearSpread = flag.Int("spread", 1, "Spread year over multiple pages")
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
var optYearA = flag.Bool("yearA", false, "Year calendar (design A)")
var optYearB = flag.Bool("yearB", false, "Year calendar (design B)")
var optCheckers = flag.Bool("checker", false, "Fill grid with checkerboard.")
var optFillpattern = flag.String("fill", "", "Set grid fill pattern.")
var optVersion = flag.Bool("v", false, "Version.")
var optICS = flag.String("ics", "", "ICS-file.")

func main() {
	flag.Var(&myFlags, "list1", "Some description for this param.")
	flag.Parse()

	if *optVersion {
		fmt.Printf("# Gocal version %s\n", VERSION)
		os.Exit(0)
	}

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
	g.SetYearSpread(*optYearSpread)
	if *optYearSpread != 1 && (!*optYearA && !*optYearB) {
		fmt.Printf("WARN: Option 'spread' ignored. Only valid for year-mode.\n")
	}

	if len(*optICS) > 0 {
		g.AddICS(*optICS)
	}

	g.AddConfig(*optConfig)
	for _, i := range myFlags {
		g.AddConfig(i)
	}
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
	if *optHideOtherMonths == true {
		g.SetHideOtherMonth()
	}
	g.SetFontScale(*optFontScale)
	g.SetWallpaper(*optWallpaper)
	g.SetPhotos(*optPhotos)
	g.SetPhoto(*optPhoto)
	g.SetFooter(*optFooter)
	g.SetFillpattern(*optFillpattern)
	/*
	  g.AddEvent(31, 1, "one", "")
	  g.AddEvent(28, 2, "two", "")
	  g.AddEvent(31, 3, "three", "")
	*/
	if *optYearA == true {
		g.CreateYearCalendar(*outfilename)
	} else if *optYearB == true {
		g.CreateYearCalendarInverse(*outfilename)
	} else {
		g.CreateCalendar(*outfilename)
	}
}
