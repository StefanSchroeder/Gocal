// Copyright (c) 2014 Stefan Schroeder, NY, 2014-03-10
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

/*
This is gocal a tool to generate calendars in PDF for printing.

https://github.com/StefanSchroeder/Gocal

See LICENSE for license.

* Inspired by pcal
* Simplicity: Create a nice calendar with minimum effort.
* No argument: Creates a calendar for this year.
* One argument: Year
* Two argument: Month Year
* Two argument: MonthBegin MonthEnd Year
* Week number
* Moonphase
* Add events from configuration file
* Set papersize
* Choose fonts (limited)
* Several languages
* Day of year
* background image
* Photo calendar
*
*/

package main

import (
	_ "code.google.com/p/go-charset/data"
	"code.google.com/p/gofpdf"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonphase"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// Default values for cmdline parameters.
	DEFAULTCONFIGFILE  = "gocal.xml"
	DEFAULTFOOTER      = "Gocal"
	DEFAULTPAPERSIZE   = "A4"
	DEFAULTORIENTATION = "L"
	DEFAULTOUTPUT      = "output.pdf"
	DEFAULTFONT        = "serif"

	// Layout parameters
	LINES      = 6
	COLUMNS    = 7
	MARGIN     = 10.0 // MM
	CELLMARGIN = 1.0

	// Colors
	DARKGREY  = 150
	LIGHTGREY = 170
	BLACK     = 0
)

var (
	// Font sizes
	EVENTFONTSIZE    = 10.0
	HEADERFONTSIZE   = 32.0
	WEEKFONTSIZE     = 12.0
	WEEKDAYFONTSIZE  = 16.0
	DOYFONTSIZE      = 12.0
	MONTHDAYFONTSIZE = 32.0
	FOOTERFONTSIZE   = 12.0
)

var optFont = flag.String("font", DEFAULTFONT, "Font")
var optFooter = flag.String("footer", DEFAULTFOOTER, "Footer note")
var optHideDOY = flag.Bool("nodoy", false, "Hide day of year (false)")
var optPlain = flag.Bool("plain", false, "Hide everything")
var optHideEvents = flag.Bool("noevents", false, "Hide events from config file (false)")
var optHideMoon = flag.Bool("nomoon", false, "Hide moon phases (false)")
var optHideWeek = flag.Bool("noweek", false, "Hide week number (false)")
var optLocale = flag.String("lang", "", "Language")
var optOrientation = flag.String("p", DEFAULTORIENTATION, "Orientation (L)andscape/(P)ortrait")
var optPaper = flag.String("paper", DEFAULTPAPERSIZE, "Paper format (A3 A4 A5 Letter Legal)")
var optPhoto = flag.String("photo", "", "Show photo (single image PNG JPG GIF)")
var optConfig = flag.String("config", DEFAULTCONFIGFILE, "Configuration file")
var optPhotos = flag.String("photos", "", "Show photos (directory PNG JPG GIF)")
var optWallpaper = flag.String("wall", "", "Show wallpaper PNG JPG GIF")
var outfilename = flag.String("o", DEFAULTOUTPUT, "Output filename")
var optSmall = flag.Bool("small", false, "Smaller fonts")
var optHideOtherMonths = flag.Bool("noother", false, "Hide neighboring month days")
var optNoclear = flag.Bool("noclear", false, "Don't delete temporary directory.")

var moonSize = 4.0
var photoList [13]string
var calFont string
var wallpaperFilename string
var fontTempdir string

// Gocaldate is a type to store single events
type gDate struct {
	Month   time.Month
	Day     int
	Text    string
	Weekday string
	Image   string
}

// Gocaldate is an XML type to store single events
type Gocaldate struct {
	Date  string `xml:"date,attr"`
	Text  string `xml:"text,attr"`
	Image string `xml:"image,attr"`
	//	Month   time.Month
	//	Day     int
	//	Weekday string
}

// TelegramStore is a container to read XML event-list
type TelegramStore struct {
	XMLName   xml.Name `xml:"Gocal"`
	Gocaldate []Gocaldate
}

// monthRange stores begin and end month of the year
type monthRange struct {
	begin int
	end   int
}

// myPdf is an anonymous struct that allows to define methods on non-local types.
type myPdf struct {
	*gofpdf.Fpdf
}

func (pdf myPdf) FullMoon(x, y float64) {
	pdf.Circle(x, y, moonSize, "F")
}

func (pdf myPdf) NewMoon(x, y float64) {
	pdf.Circle(x, y, moonSize, "D")
}

func (pdf myPdf) FirstQuarter(x, y float64) {
	pdf.Arc(x, y, moonSize, moonSize, 0.0, 90.0, 270.0, "F")
}

func (pdf myPdf) LastQuarter(x, y float64) {
	pdf.Arc(x, y, moonSize, moonSize, 0.0, 270.0, 270.0+180.0, "F")
}

type pdfWriter struct {
	pdf *gofpdf.Fpdf
	fl  *os.File
}

func (pw *pdfWriter) Write(p []byte) (n int, err error) {
	if pw.pdf.Ok() {
		return pw.fl.Write(p)
	}
	return
}

func (pw *pdfWriter) Close() (err error) {
	if pw.fl != nil {
		pw.fl.Close()
		pw.fl = nil
	}
	if pw.pdf.Ok() {
		fmt.Printf("Generated '%v'.\n", *outfilename)
	} else {
		fmt.Printf("%s\n", pw.pdf.Error())
	}
	return
}

func docWriter(pdf *gofpdf.Fpdf) *pdfWriter {
	pw := new(pdfWriter)
	pw.pdf = pdf
	if pdf.Ok() {
		var err error
		pw.fl, err = os.Create(*outfilename)
		if err != nil {
			pdf.SetErrorf("# Error opening output file.")
		}
	}
	return pw
}

// computeMoonphases fills a map with moonphase information.
func computeMoonphases(moon map[int]string, da int, mo int, yr int) {
	daysInYear := 365
	if julian.LeapYearGregorian(yr) {
		daysInYear = 366
	}
	// Look at every day and check if it has any of the Moon Phases.
	for i := 0; i < 32; i++ {
		dayOfYear := julian.DayOfYearGregorian(yr, mo, int(da)+i)
		decimalYear := float64(yr) +
			float64(dayOfYear-1)/float64(daysInYear)
		jdeNew := moonphase.New(decimalYear)
		y, m, d := julian.JDToCalendar(jdeNew)
		if (y == yr) && (m == mo) && (int(d) == i) {
			//fmt.Printf("New moon on %d\n", int(d))
			moon[int(d)] = "New"
		}
		jdeNew = moonphase.Full(decimalYear)
		y, m, d = julian.JDToCalendar(jdeNew)
		if (y == yr) && (m == mo) && (int(d) == i) {
			//fmt.Printf("Full moon on %d\n", int(d))
			moon[int(d)] = "Full"
		}
		jdeNew = moonphase.First(decimalYear)
		y, m, d = julian.JDToCalendar(jdeNew)
		if (y == yr) && (m == mo) && (int(d) == i) {
			//fmt.Printf("First Q moon on %d\n", int(d))
			moon[int(d)] = "First"
		}
		jdeNew = moonphase.Last(decimalYear)
		y, m, d = julian.JDToCalendar(jdeNew)
		if (y == yr) && (m == mo) && (int(d) == i) {
			moon[int(d)] = "Last"
			//fmt.Printf("Last Q moon on %d\n", int(d))
		}
	}
}

func removeTempdir(d string) {
	if *optNoclear == true {
		return
	}
	os.RemoveAll(d)
}

func processFont(fontFile string) (fontName, tempDirname string) {
	var err error
	tempDirname, err = ioutil.TempDir("", "")
	if fontFile == "mono" {
		fontFile = tempDirname + string(os.PathSeparator) + "freemonobold.ttf"
		ioutil.WriteFile(fontFile, getFreeMonoBold(), 0700)
	} else if fontFile == "serif" {
		fontFile = tempDirname + string(os.PathSeparator) + "freeserifbold.ttf"
		ioutil.WriteFile(fontFile, getFreeSerifBold(), 0700)
	} else if fontFile == "sans" {
		fontFile = tempDirname + string(os.PathSeparator) + "freesansbold.ttf"
		ioutil.WriteFile(fontFile, getFreeSansBold(), 0700)
	}
	err = ioutil.WriteFile(tempDirname+string(os.PathSeparator)+"cp1252.map", []byte(codepageCP1252), 0700)
	err = gofpdf.MakeFont(fontFile, tempDirname+string(os.PathSeparator)+"cp1252.map", tempDirname, os.Stderr, true)
	_ = err
	// FIXME Do some error checking here.
	fontName = filepath.Base(fontFile)
	fontName = strings.TrimSuffix(fontName, filepath.Ext(fontName))
	// fmt.Printf("Using external font: %v\n", fontName)
	return fontName, tempDirname
}

func main() {
	flag.Parse()

	if *optPlain == true {
		*optHideEvents = true
		*optHideMoon = true
		*optHideDOY = true
		*optHideWeek = true
	}

	if *optSmall == true {
		EVENTFONTSIZE *= 0.75
		HEADERFONTSIZE *= 0.75
		WEEKFONTSIZE *= 0.75
		WEEKDAYFONTSIZE *= 0.75
		DOYFONTSIZE *= 0.75
		MONTHDAYFONTSIZE *= 0.75
		FOOTERFONTSIZE *= 0.75
	}

	testedLanguage := map[string]bool{
		"en_US": true,
		"en_GB": true,
		"da_DK": true,
		"nl_BE": true,
		"nl_NL": true,
		"fi_FI": true,
		"fr_FR": true,
		"fr_CA": true,
		"de_DE": true,
		"hu_HU": true,
		"it_IT": true,
		"nn_NO": true,
		"nb_NO": true,
		"pt_PT": true,
		"pt_BR": true,
		"ro_RO": true,
		"ru_RU": true,
		"es_ES": true,
		"sv_SE": true,
		"tr_TR": true,
		"bg_BG": true,
		"zh_CN": true,
		"zh_TW": true,
		"zh_HK": true,
	}

	// First try Environment
	currentLanguage := os.Getenv("LANG")

	// If set on the cmdline, override
	if *optLocale != "" {
		currentLanguage = *optLocale
	}

	// if we don't know that language, fall back to en.
	if testedLanguage[currentLanguage] != true {
		currentLanguage = "en_US"
	}

	eventList := make([]gDate, 1000)

	if *optHideEvents == false {
		eventList = readConfigurationfile(*optConfig)
	}

	var wantyear = int(time.Now().Year())
	wantmonths := monthRange{1, 12}

	if flag.NArg() == 1 {
		dummyyear, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		wantyear = int(dummyyear)
	} else if flag.NArg() == 2 {
		dummymonth, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		dummyyear, _ := strconv.ParseInt(flag.Arg(1), 10, 32)
		wantmonths.begin = int(dummymonth)
		wantmonths.end = int(dummymonth)
		wantyear = int(dummyyear)
	} else if flag.NArg() == 3 {
		dummymonthBegin, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		dummymonthEnd, _ := strconv.ParseInt(flag.Arg(1), 10, 32)
		dummyyear, _ := strconv.ParseInt(flag.Arg(2), 10, 32)
		wantmonths.begin = int(dummymonthBegin)
		wantmonths.end = int(dummymonthEnd)
		wantyear = int(dummyyear)
	}

	localizedMonthNames := getLocalizedMonthNames(currentLanguage)
	localizedWeekdayNames := getLocalizedWeekdayNames(currentLanguage)

	if *optFont != "" {
		calFont = *optFont
	}

	calFont, fontTempdir = processFont(calFont)

	pdf := gofpdf.New(*optOrientation, "mm", *optPaper, fontTempdir)
	pdf.SetTitle("Created with Gocal", true)
	pdf.AddFont(calFont, "", calFont+".json")

	PAGEWIDTH, PAGEHEIGHT, _ := pdf.PageSize(0)
	if *optOrientation != "P" {
		PAGEWIDTH, PAGEHEIGHT = PAGEHEIGHT, PAGEWIDTH
	}

	cw := (PAGEWIDTH - 2*MARGIN) / COLUMNS // cellwidth w margin
	ch := PAGEHEIGHT / (LINES + 2)         // cellheight

	if *optPhoto != "" {
		ch *= 0.5
		moonSize *= 0.6 // make moon smaller on photopage
		for i := 0; i < 13; i++ {
			photoname := *optPhoto
			if strings.HasPrefix(photoname, "http://") {
				photoname = downloadFile(photoname, fontTempdir)
			}
			photoList[i] = photoname
		}
	}
	if *optPhotos != "" {
		ch *= 0.5
		moonSize *= 0.6 // make moon smaller on photopage
		fileList, err := filepath.Glob(*optPhotos + string(os.PathSeparator) + "*")
		if err == nil {
			for i := 0; i < 13; i++ {
				photoList[i] = fileList[i%len(fileList)]
			}
		} else {
			fmt.Printf("# There is an error in your path to photos: %v\n", err)
		}
	}

	calendarTable := func(mymonth int, myyear int) {
		pdf.SetFont(calFont, "", WEEKDAYFONTSIZE)
		for weekday := 1; weekday < 8; weekday++ { // Print weekdays in first row
			pdf.CellFormat(cw, 7, localizedWeekdayNames[weekday], "0", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)

		// Figure out the first day in the calendar which depends on the weekday
		// of the first day
		var day int64 = 1
		t := time.Date(myyear, time.Month(mymonth), 1, 0, 0, 0, 0, time.UTC)

		day -= int64(t.Weekday())
		if day > 0 { // adjust silly exception where month starts w/ Sunday.
			day -= 7
		}

		moon := make(map[int]string)
		if *optHideMoon == false {
			computeMoonphases(moon, int(day), mymonth, myyear)
		}

		for i := 0; i < LINES; i++ {
			for j := 0; j < COLUMNS; j++ {
				fill := false
				nd := time.Date(myyear, time.Month(mymonth), 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(day) * 24 * 60 * 60 * time.Second)

				// Determine color
				if nd.Month() != time.Month(mymonth) { // GREY
					pdf.SetTextColor(DARKGREY, DARKGREY, DARKGREY)
					pdf.SetFillColor(LIGHTGREY, LIGHTGREY, LIGHTGREY)
					fill = false // FIXME, do we want fill here?
				} else if nd.Weekday() == time.Saturday || nd.Weekday() == time.Sunday {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(BLACK, BLACK, BLACK)
				}

				if *optHideOtherMonths == true && nd.Month() != time.Month(mymonth) { // GREY
					pdf.SetX(pdf.GetX() + cw)
					day++
					continue
				}
				pdf.SetCellMargin(CELLMARGIN)

				// Add moon icon
				if m, ok := moon[int(day)]; ok == true {
					x, y := pdf.GetXY()
					moonLocX, moonLocY := x+cw*0.82, y+ch*0.2

					myMoonPDF := myPdf{pdf}
					switch m {
					case "Full":
						myMoonPDF.FullMoon(moonLocX, moonLocY)
					case "New":
						myMoonPDF.NewMoon(moonLocX, moonLocY)
					case "First":
						myMoonPDF.FirstQuarter(moonLocX, moonLocY)
					case "Last":
						myMoonPDF.LastQuarter(moonLocX, moonLocY)
					}
				}

				// Day of year, lower right
				if *optHideDOY == false && int(nd.Month()) == mymonth {
					doy := julian.DayOfYearGregorian(myyear, mymonth, int(nd.Day()))
					pdf.SetFont(calFont, "", DOYFONTSIZE)
					pdf.CellFormat(cw, ch, fmt.Sprintf("%d", doy), "1", 0, "BR", fill, 0, "")
					pdf.SetX(pdf.GetX() - cw) // reset
				}

				// Add week number, lower left
				if nd.Weekday() == time.Monday && *optHideWeek == false {
					pdf.SetFont(calFont, "", WEEKFONTSIZE)
					_, weeknr := nd.ISOWeek()
					pdf.CellFormat(cw, ch, fmt.Sprintf("W %d", weeknr), "1", 0, "BL", fill, 0, "")
					pdf.SetX(pdf.GetX() - cw) // reset
				}

				// Add event text
				for _, ev := range eventList {
					if nd.Day() == ev.Day && nd.Month() == ev.Month {
						x, y := pdf.GetXY()
						pdf.SetFont(calFont, "", EVENTFONTSIZE)

            if ev.Image != "" {
              pdf.Image(ev.Image, x, y, cw, ch, false, "", 0, "")
            }
						for i, j := range strings.Split(ev.Text, "\\n") {
							pdf.Text(x+0.02*cw, y+0.70*ch+float64(i)*EVENTFONTSIZE/4.0, fmt.Sprintf("%s", j))
						}
					}
				}

				// day of the month, big number
				pdf.SetFont(calFont, "", MONTHDAYFONTSIZE)
				pdf.CellFormat(cw, ch, fmt.Sprintf("%d", nd.Day()), "1", 0, "TL", fill, 0, "")
				day++
			}
			pdf.Ln(-1)
		}
	}

	for mo := wantmonths.begin; mo <= wantmonths.end; mo++ {
		//fmt.Printf("Printing page %d\n", page)
		pdf.AddPage()
		if *optWallpaper != "" {
			wallpaperFilename = *optWallpaper
			if strings.HasPrefix(wallpaperFilename, "http://") {
				wallpaperFilename = downloadFile(*optWallpaper, fontTempdir)
			}
			pdf.Image(wallpaperFilename, 0, 0, PAGEWIDTH, PAGEHEIGHT, false, "", 0, "")
		}
		if *optPhoto != "" || *optPhotos != "" {
			photo := photoList[mo-1]
			if photo != "" {
				pdf.Image(photo, 0, PAGEHEIGHT*0.5, PAGEWIDTH, PAGEHEIGHT*0.5, false, "", 0, "")
			}
		}

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		pdf.SetFont(calFont, "", HEADERFONTSIZE)
		pdf.CellFormat(PAGEWIDTH-MARGIN, MARGIN, localizedMonthNames[mo]+" "+fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		calendarTable(mo, wantyear)

		pdf.Ln(-1)
		pdf.SetFont(calFont, "", FOOTERFONTSIZE)
		pdf.Text(0.50*PAGEWIDTH, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", *optFooter))
	}

	pdf.OutputAndClose(docWriter(pdf))
	removeTempdir(fontTempdir)
}
