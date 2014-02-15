/*

This is gocal a tool to generate calendars in PDF for printing.

* inspired by pcal
* simplicity: create a nice calendar programmatically with minimum effort.
* show week number
* moonphase
* add events from configuration file
* set papersize
* add month notes
* choose fonts
* add images to days
* multilang
* sunrise/sunset for location

*/

package main

import (
	"code.google.com/p/gofpdf"
	"flag"
	"fmt"
	"github.com/soniakeys/meeus/base"
	"github.com/soniakeys/meeus/globe"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonillum"
	"github.com/soniakeys/meeus/moonphase"
	"github.com/soniakeys/meeus/rise"
	"github.com/soniakeys/meeus/sidereal"
	"github.com/soniakeys/meeus/solar"
	"math"
	"os"
	"strconv"
	"time"
)

// Absolute path needed for gocov tool; relative OK for test
const (
	CONFIGFILE = ".gocal.json"

	LINES       = 6
	cnGofpdfDir = "."
	cnFontDir   = cnGofpdfDir + "/font"
	COLUMNS     = 7
	darkgrey    = 120
	lightgrey   = 170
	black       = 0
	moonSize    = 4.0
	MARGIN      = 10.0
	MAINFONT    = "Times"
)

var (
	// USNO gives the coordiates it uses in decimal degrees, so I do a simple
	// conversion to radians like this rather than use base.NewAngle as in
	// examples in the package doc.
	bostonLon = 71.1 * math.Pi / 180
	bostonLat = 42.3 * math.Pi / 180

	h0 = rise.Stdh0Solar // events of interest are sunrise, sunset

	tz = -5 * 3600. // EST time zone correction, in seconds
)

var orientation = flag.String("p", "L", "Orientation (L)andscape/(P)ortrait")
var outfilename = flag.String("o", "output.pdf", "Output filename")

var entryDates = map[string]string{
	"January":      "Winter",
	"1.1.":         "New Year",
	"15. January":  "Jule",
	"14. February": "Enno",
	"1/6":          "Allerheiligen",
	"5/23":         "Geburtstag",
	"5/20":         "Birgit",
	"9/18":         "Tomke",
	"Sunday":       "Flea market",
}

type pdfWriter struct {
	pdf *gofpdf.Fpdf
	fl  *os.File
}

type calDate struct {
	day  string
	text string
}

func printSunrise(y, m, d int) (sunrise_sunset string) {
	p := globe.Coord{bostonLat, bostonLon}
	jd := julian.CalendarGregorianToJD(y, m, float64(d)) // date of interest
	Th0 := sidereal.Apparent0UT(jd)
	ra, dec := solar.ApparentEquatorial(jd) // position of sun
	mRise, _, mSet, err := rise.ApproxTimes(p, h0, Th0, ra, dec)
	if err != nil {
		fmt.Println(err)
		return
	}
	rise := fmt.Sprintf("%02d", int(base.NewFmtTime(mRise+tz).Hour())) + ":" +
		fmt.Sprintf("%02d", int(base.NewFmtTime(mRise+tz).Min())%60) + "\n"
	sunset := fmt.Sprintf("%02d", int(base.NewFmtTime(mSet+tz).Hour())) + ":" +
		fmt.Sprintf("%02d", int(base.NewFmtTime(mSet+tz).Min())%60) + "\n"

	sunrise_sunset = rise + "/" + sunset
	fmt.Println(y, m, d, base.NewFmtTime(mRise+tz))
	return
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
		fmt.Printf("Successfully generated \n")
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
			pdf.SetErrorf("Error opening output file ")
		}
	}
	return pw
}

func getPhase(y, m, d int) (r float64) { // future use
	i := moonillum.PhaseAngle3(julian.CalendarGregorianToJD(y, m, float64(d)))
	k := base.Illuminated(i)
	r = k
	return r
}

func main() {
	flag.Parse()

	var wantyear int = int(time.Now().Year())
	var wantmonth int = int(time.Now().Month())

	if flag.NArg() == 1 {
		dummyyear, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		wantyear = int(dummyyear)
	} else if flag.NArg() == 2 {
		dummymonth, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		dummyyear, _ := strconv.ParseInt(flag.Arg(1), 10, 32)
		wantmonth = int(dummymonth)
		wantyear = int(dummyyear)
	}

	paperformat := "A4"

	pdf := gofpdf.New(*orientation, "mm", paperformat, cnFontDir)
	PAGEWIDTH, PAGEHEIGHT, _ := pdf.PageSize(0)
	if *orientation != "P" {
		PAGEWIDTH, PAGEHEIGHT = PAGEHEIGHT, PAGEWIDTH
	}

	cw := (PAGEWIDTH - 2*MARGIN) / COLUMNS // cellwidth w 20mm margin
	ch := PAGEHEIGHT / (LINES + 2)         // cellheight
	calendarTable := func(mymonth int) {
		wantmonth = mymonth
		for titleDay := 1; titleDay < 8; titleDay++ {
			pdf.CellFormat(cw, 7, fmt.Sprintf("%s", time.Weekday(titleDay%7)), "0", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		var day int64 = 1

		t := time.Date(wantyear, time.Month(wantmonth), 1, 0, 0, 0, 0, time.UTC)
		day -= int64(t.Weekday())
		if day > 0 { // adjust silly exception where month starts w/ Sunday.
			day -= 7
		}

		daysInYear := 365
		if julian.LeapYearGregorian(wantyear) {
			daysInYear = 366
		}

		moon := make(map[int]string)

		// This is a perfect example for the difference btw. efficient and
		// effective. It's effective, but not efficient:
		// Look at every day and check it has any of the Moon Phases.
		for i := 0; i < 32; i++ {
			decimalYear := float64(wantyear) +
				float64(julian.DayOfYearGregorian(wantyear, wantmonth, int(day)+i)-1)/float64(daysInYear)
			jdeNew := moonphase.New(decimalYear)
			y, m, d := julian.JDToCalendar(jdeNew)
			if (y == wantyear) && (m == wantmonth) && (int(d) == i) {
				//fmt.Printf("New moon on %d\n", int(d))
				moon[int(d)] = "New"
			}
			jdeNew = moonphase.Full(decimalYear)
			y, m, d = julian.JDToCalendar(jdeNew)
			if (y == wantyear) && (m == wantmonth) && (int(d) == i) {
				//fmt.Printf("Full moon on %d\n", int(d))
				moon[int(d)] = "Full"
			}
			jdeNew = moonphase.First(decimalYear)
			y, m, d = julian.JDToCalendar(jdeNew)
			if (y == wantyear) && (m == wantmonth) && (int(d) == i) {
				//fmt.Printf("First Q moon on %d\n", int(d))
				moon[int(d)] = "First"
			}
			jdeNew = moonphase.Last(decimalYear)
			y, m, d = julian.JDToCalendar(jdeNew)
			if (y == wantyear) && (m == wantmonth) && (int(d) == i) {
				moon[int(d)] = "Last"
				//fmt.Printf("Last Q moon on %d\n", int(d))
			}
		}

		pdf.SetFont(MAINFONT, "B", 24)
		for i := 0; i < LINES; i++ {
			for j := 0; j < COLUMNS; j++ {
				var fill bool = false
				nd := time.Date(wantyear, time.Month(wantmonth), 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(day) * 24 * 60 * 60 * time.Second)

				if nd.Month() != time.Month(wantmonth) { // GREY
					pdf.SetTextColor(darkgrey, darkgrey, darkgrey)
					pdf.SetFillColor(lightgrey, lightgrey, lightgrey)
					fill = true
				} else if nd.Weekday() == time.Saturday || nd.Weekday() == time.Sunday {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(black, black, black) // BLACK
				}

				sunrise_sunset := ""
				if nd.Weekday() == time.Sunday {
					sunrise_sunset = printSunrise(nd.Year(), int(nd.Month()), int(nd.Day()))
				}

				pdf.SetCellMargin(1)

				cellString := fmt.Sprintf("%d", nd.Day())
				m, ok := moon[int(day)]
				if ok {
					// Moonphase as Text: cellString += " " + m
					x, y := pdf.GetXY()
					moonLocX, moonLocY := x+cw*0.8, y+ch*0.3

					switch m {
					case "Full":
						pdf.Circle(moonLocX, moonLocY, moonSize, "F") // Full
					case "New":
						pdf.Circle(moonLocX, moonLocY, moonSize, "D") // New
					case "First":
						pdf.Arc(moonLocX, moonLocY, moonSize, moonSize, 0.0, 90.0, 270.0, "F") // 1st Q
					case "Last":
						pdf.Arc(moonLocX, moonLocY, moonSize, moonSize, 0.0, 270.0, 270.0+180.0, "F") // Last Q
					}
				}
				if nd.Weekday() == time.Monday {
					x, y := pdf.GetXY()
					pdf.SetFont(MAINFONT, "", 12)
					_, weeknr := nd.ISOWeek()
					pdf.Text(x+0.02*cw, y+0.95*ch, fmt.Sprintf("Week %d", weeknr))
				}

				if nd.Weekday() == time.Sunday {
					x, y := pdf.GetXY()
					pdf.SetFont(MAINFONT, "", 12)

					pdf.Text(x+0.02*cw, y+0.95*ch, fmt.Sprintf("%s", sunrise_sunset))
				}

				pdf.SetFont(MAINFONT, "B", 32)
				pdf.CellFormat(cw, ch, cellString, "1", 0, "TL", fill, 0, "")
				day++
			}
			pdf.Ln(-1)
		}

	}

	for page := 1; page < 13; page++ {
		//fmt.Printf("Printing page %d\n", page)
		wantmonth = page
		pdf.AddPage()
		pdf.SetTextColor(black, black, black) // BLACK
		pdf.SetFont(MAINFONT, "B", 32)
		pdf.CellFormat(PAGEWIDTH-10.0, 10.0, fmt.Sprintf("%s", time.Month(wantmonth))+" "+fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetTextColor(black, black, black) // BLACK
		pdf.SetFont(MAINFONT, "B", 16)
		calendarTable(page)
	}

	pdf.OutputAndClose(docWriter(pdf))
}
