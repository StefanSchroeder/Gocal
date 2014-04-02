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

Gocal.New()
Gocal.Write()


*/
package gocal

import (
	_ "code.google.com/p/go-charset/data"
	"code.google.com/p/gofpdf"
	"encoding/xml"
	"fmt"
	"github.com/soniakeys/meeus/julian"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// Layout parameters
	LINES      = 6
	COLUMNS    = 7
	MARGIN     = 10.0 // MM
	CELLMARGIN = 1.0

	// Colors
	DARKGREY  = 150
	LIGHTGREY = 170
	BLACK     = 0

  MOONSIZE = 4.0

	// Font sizes
	EVENTFONTSIZE    = 10.0
	HEADERFONTSIZE   = 32.0
	WEEKFONTSIZE     = 12.0
	WEEKDAYFONTSIZE  = 16.0
	DOYFONTSIZE      = 12.0
	MONTHDAYFONTSIZE = 32.0
	FOOTERFONTSIZE   = 12.0
)

type Calendar struct {
	WantBeginMonth     int
	WantEndMonth       int
	WantYear           int
	OptFont            string
	OptFooter          string
	OptOrientation     string
	OptSmall           bool
	OptPaperformat     string
	OptLocale          string
	OptHideOtherMonths bool
	OptWallpaper       string
	OptHideMoon        bool
	OptHideWeek        bool
	OptHideDOY         bool
	OptPhoto           string
	OptPlain           bool
	OptConfig          string
	OptPhotos          string
	OptFontScale          float64
}

func New(b int, e int, y int) *Calendar {
	return &Calendar{b, e, y,
		"",      // OptFont
		"",      // OptFooter
		"P",     // OptOrientation P=portrait
		false,   // OptSmall
		"A4",    // OptPaperformat
		"en_US", // OptLocale
		false,   // OptHideOtherMonths
		"",      // OptWallpaper
		false,   // OptHideMoon
		false,   // OptHideWeek
		false,   // OptHideDOY
		"",      // OptPhoto
		false,   // OptPlain
		"",      // OptConfig
		"",      // OptPhotos
		1.0,      // OptFontScale
	}
}

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

// myPdf is an anonymous struct that allows to define methods on non-local type.
type myPdf struct {
	*gofpdf.Fpdf
  moonSize float64
}

func (pdf myPdf) fullMoon(x, y float64) {
	pdf.Circle(x, y, pdf.moonSize, "F")
}

func (pdf myPdf) newMoon(x, y float64) {
	pdf.Circle(x, y, pdf.moonSize, "D")
}

func (pdf myPdf) firstQuarter(x, y float64) {
	pdf.Arc(x, y, pdf.moonSize, pdf.moonSize, 0.0, 90.0, 270.0, "F")
}

func (pdf myPdf) lastQuarter(x, y float64) {
	pdf.Arc(x, y, pdf.moonSize, pdf.moonSize, 0.0, 270.0, 270.0+180.0, "F")
}

type pdfWriter struct {
	pdf *gofpdf.Fpdf
	fl  *os.File
  pdfFilename string
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
		fmt.Printf("Generated '%v'.\n", pw.pdfFilename)
	} else {
		fmt.Printf("%s\n", pw.pdf.Error())
	}
	return
}

func docWriter(pdf *gofpdf.Fpdf, fname string) *pdfWriter {
	pw := new(pdfWriter)
  pw.pdfFilename = fname
	pw.pdf = pdf
	if pdf.Ok() {
		var err error
		pw.fl, err = os.Create(pw.pdfFilename)
		if err != nil {
			pdf.SetErrorf("# Error opening output file.")
		}
	}
	return pw
}

func (g *Calendar) SetPlain() {
	g.OptPlain = true
}

func (g *Calendar) SetHideOtherMonth() {
	g.OptHideOtherMonths = true
}

func (g *Calendar) SetHideDOY() {
	g.OptHideDOY = true
}

func (g *Calendar) SetHideMoon() {
	g.OptHideMoon = true
}

func (g *Calendar) SetHideWeek() {
	g.OptHideWeek = true
}

func (g *Calendar) SetSmall() {
	g.OptSmall = true
}

func (g *Calendar) SetFont(f string) {
	g.OptFont = f
}

func (g *Calendar) SetFontScale(f float64) {
	g.OptFontScale = f
}

func (g *Calendar) SetPhotos(f string) {
	g.OptPhotos = f
}

func (g *Calendar) SetPhoto(f string) {
	g.OptPhoto = f
}

func (g *Calendar) SetConfig(f string) {
	g.OptConfig = f
}

func (g *Calendar) SetLocale(f string) {
	g.OptLocale = f
}

func (g *Calendar) SetOrientation(f string) {
	g.OptOrientation = f
}

func (g *Calendar) SetWallpaper(f string) {
	g.OptWallpaper = f
}

func (g *Calendar) SetFooter(f string) {
	g.OptFooter = f
}

func (g *Calendar) SetPaperformat(f string) {
	g.OptPaperformat = f
}

func (g *Calendar) CreateCalendar(fn string) {

	var fontTempdir string
  var fontScale = g.OptFontScale

	if g.OptPlain == true {
		g.SetHideOtherMonth()
		g.SetHideDOY()
		g.SetHideMoon()
		g.SetHideWeek()
	}

	if g.OptSmall == true {
    fontScale = 0.75
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
	if g.OptLocale != "" {
		currentLanguage = g.OptLocale
	}

	// if we don't know that language, fall back to en.
	if testedLanguage[currentLanguage] != true {
		currentLanguage = "en_US"
	}

	var eventList = make([]gDate, 10000) // Maximum number of events

	if g.OptConfig != "" {
		eventList = readConfigurationfile(g.OptConfig)
	}

	wantyear := g.WantYear
	wantmonths := monthRange{g.WantBeginMonth, g.WantEndMonth}
	localizedMonthNames := getLocalizedMonthNames(currentLanguage)
	localizedWeekdayNames := getLocalizedWeekdayNames(currentLanguage)

  var calFont string
	if g.OptFont != "" {
		calFont = g.OptFont
	}

	calFont, fontTempdir = processFont(calFont)

	pdf := gofpdf.New(g.OptOrientation, "mm", g.OptPaperformat, fontTempdir)
	pdf.SetTitle("Created with Gocal", true)
	pdf.AddFont(calFont, "", calFont+".json")

	PAGEWIDTH, PAGEHEIGHT, _ := pdf.PageSize(0)
	if g.OptOrientation != "P" {
		PAGEWIDTH, PAGEHEIGHT = PAGEHEIGHT, PAGEWIDTH
	}

	cw := (PAGEWIDTH - 2*MARGIN) / COLUMNS // cellwidth w margin
	ch := PAGEHEIGHT / (LINES + 2)         // cellheight

  var photoList [12]string
	if g.OptPhoto != "" {
		ch *= 0.5
		for i := 0; i < 12; i++ {
			photoname := g.OptPhoto
			if strings.HasPrefix(photoname, "http://") {
				photoname = downloadFile(photoname, fontTempdir)
			}
			photoList[i] = photoname
		}
	}

	if g.OptPhotos != "" {
		ch *= 0.5
		fileList, err := filepath.Glob(g.OptPhotos + string(os.PathSeparator) + "*")
		if err == nil {
			for i := 0; i < 12; i++ {
				photoList[i] = fileList[i%len(fileList)]
			}
		} else {
			fmt.Printf("# There is an error in your path to photos: %v\n", err)
		}
	}

	calendarTable := func(mymonth int, myyear int) {
		pdf.SetFont(calFont, "", WEEKDAYFONTSIZE * fontScale)
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
		if g.OptHideMoon == false {
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

				if g.OptHideOtherMonths == true && nd.Month() != time.Month(mymonth) {
					pdf.SetX(pdf.GetX() + cw)
					day++
					continue
				}
				pdf.SetCellMargin(CELLMARGIN)

				// Add moon icon
				if m, ok := moon[int(day)]; ok == true {
					x, y := pdf.GetXY()
					moonLocX, moonLocY := x+cw*0.82, y+ch*0.2

          moonsize := MOONSIZE
          if g.OptPhoto != "" || g.OptPhotos != ""  {
            moonsize *= 0.6
          }
					myMoonPDF := myPdf{pdf, moonsize}
					switch m {
					case "Full":
						myMoonPDF.fullMoon(moonLocX, moonLocY)
					case "New":
						myMoonPDF.newMoon(moonLocX, moonLocY)
					case "First":
						myMoonPDF.firstQuarter(moonLocX, moonLocY)
					case "Last":
						myMoonPDF.lastQuarter(moonLocX, moonLocY)
					}
				}

				// Day of year, lower right
				if g.OptHideDOY == false && int(nd.Month()) == mymonth {
					doy := julian.DayOfYearGregorian(myyear, mymonth, int(nd.Day()))
					pdf.SetFont(calFont, "", DOYFONTSIZE * fontScale)
					pdf.CellFormat(cw, ch, fmt.Sprintf("%d", doy), "1", 0, "BR", fill, 0, "")
					pdf.SetX(pdf.GetX() - cw) // reset
				}

				// Add week number, lower left
				if nd.Weekday() == time.Monday && g.OptHideWeek == false {
					pdf.SetFont(calFont, "", WEEKFONTSIZE * fontScale)
					_, weeknr := nd.ISOWeek()
					pdf.CellFormat(cw, ch, fmt.Sprintf("W %d", weeknr), "1", 0, "BL", fill, 0, "")
					pdf.SetX(pdf.GetX() - cw) // reset
				}

				// Add event text
				for _, ev := range eventList {
					if nd.Day() == ev.Day && nd.Month() == ev.Month {
						x, y := pdf.GetXY()
						pdf.SetFont(calFont, "", EVENTFONTSIZE * fontScale)

						if ev.Image != "" {
							pdf.Image(ev.Image, x, y, cw, ch, false, "", 0, "")
						}
						for i, j := range strings.Split(ev.Text, "\\n") {
							pdf.Text(x+0.02*cw, y+0.70*ch+float64(i)*EVENTFONTSIZE * fontScale/4.0, fmt.Sprintf("%s", j))
						}
					}
				}

				// day of the month, big number
				pdf.SetFont(calFont, "", MONTHDAYFONTSIZE * fontScale)
				pdf.CellFormat(cw, ch, fmt.Sprintf("%d", nd.Day()), "1", 0, "TL", fill, 0, "")
				day++
			}
			pdf.Ln(-1)
		}
	}

	for mo := wantmonths.begin; mo <= wantmonths.end; mo++ {
		//fmt.Printf("Printing page %d\n", page)
		pdf.AddPage()
		if g.OptWallpaper != "" {
      wallpaperFilename := g.OptWallpaper
			if strings.HasPrefix(wallpaperFilename, "http://") {
				wallpaperFilename = downloadFile(g.OptWallpaper, fontTempdir)
			}
			pdf.Image(wallpaperFilename, 0, 0, PAGEWIDTH, PAGEHEIGHT, false, "", 0, "")
		}

    if g.OptPhoto != "" || g.OptPhotos != ""  {
			photo := photoList[mo-1] // this list is zero-based.
			if photo != "" {
				pdf.Image(photo, 0, PAGEHEIGHT*0.5, PAGEWIDTH, PAGEHEIGHT*0.5, false, "", 0, "")
			}
		}

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		pdf.SetFont(calFont, "", HEADERFONTSIZE * fontScale)
		pdf.CellFormat(PAGEWIDTH-MARGIN, MARGIN, localizedMonthNames[mo]+" "+fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		calendarTable(mo, wantyear)

		pdf.Ln(-1)
		pdf.SetFont(calFont, "", FOOTERFONTSIZE * fontScale)
		pdf.Text(0.50*PAGEWIDTH, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", g.OptFooter))
	}
	pdf.OutputAndClose(docWriter(pdf, fn))
	removeTempdir(fontTempdir)
}

