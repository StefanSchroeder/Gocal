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
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/soniakeys/meeus/base"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonillum"
	"github.com/soniakeys/meeus/moonphase"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CONFIGFILE    = "gocal.xml"
	LINES         = 6
	cnGofpdfDir   = "."
	cnFontDir     = cnGofpdfDir + "/font"
	COLUMNS       = 7
	darkgrey      = 150
	lightgrey     = 170
	black         = 0
	moonSize      = 4.0
	MARGIN        = 10.0
	MAINFONT      = "Times"
	eventFontsize = 10
)

var orientation = flag.String("p", "L", "Orientation (L)andscape/(P)ortrait")
var outfilename = flag.String("o", "output.pdf", "Output filename")
var footer = flag.String("f", "Gocal", "Footer note")

type gocalDate struct {
	Month   time.Month
	Day     int
	Text    string
	Weekday string
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

type Gocaldate struct {
	Date string `xml:"date,attr"`
	Text string `xml:"text,attr"`
	//	Month   time.Month
	//	Day     int
	//	Weekday string
}

type TelegramStore struct {
	XMLName   xml.Name `xml:"Gocal"`
	Gocaldate []Gocaldate
}

func readConfigurationfile(filename string) (eL []gocalDate) {

	var v TelegramStore

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	v = TelegramStore{}
	err2 := xml.Unmarshal([]byte(data), &v)
	if err2 != nil {
		log.Fatal("# ERROR: when trying to unmarshal the XML configuration file: %v", err2)
		return
	}

	for _, m := range v.Gocaldate {

		if strings.Index(m.Date, "/") != -1 { // Is this Month/Day ?

			textArray := strings.Split(m.Date, "/")

			if textArray[0] == "*" {
				d, _ := strconv.ParseInt(textArray[1], 10, 32)
				for j := 1; j < 13; j++ {
					gcd := gocalDate{time.Month(j), int(d), m.Text, ""}
					eL = append(eL, gcd)
				}
			} else {
				mo, _ := strconv.ParseInt(textArray[0], 10, 32)
				d, _ := strconv.ParseInt(textArray[1], 10, 32)

				gcd := gocalDate{time.Month(mo), int(d), m.Text, ""}
				eL = append(eL, gcd)
			}
		}
	}

	return eL
}

func main() {
	flag.Parse()

	eventList := make([]gocalDate, 1000)

	eventList = readConfigurationfile(CONFIGFILE)

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
		for titleDay := 1; titleDay < 8; titleDay++ { // Print Weekdays
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
		// Look at every day and check if it has any of the Moon Phases.
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

				// Determine color
				if nd.Month() != time.Month(wantmonth) { // GREY
					pdf.SetTextColor(darkgrey, darkgrey, darkgrey)
					pdf.SetFillColor(lightgrey, lightgrey, lightgrey)
					fill = true
				} else if nd.Weekday() == time.Saturday || nd.Weekday() == time.Sunday {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(black, black, black) // BLACK
				}

				pdf.SetCellMargin(1)

				// Add day-of-month number
				cellString := fmt.Sprintf("%d", nd.Day())

				// Add moon icon
				if m, ok := moon[int(day)]; ok == true {
					// Moonphase as Text: cellString += " " + m
					x, y := pdf.GetXY()
					moonLocX, moonLocY := x+cw*0.82, y+ch*0.2

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

				// Add week number
				if nd.Weekday() == time.Monday {
					x, y := pdf.GetXY()
					pdf.SetFont(MAINFONT, "", 12)
					_, weeknr := nd.ISOWeek()
					pdf.Text(x+0.02*cw, y+0.95*ch, fmt.Sprintf("Week %d", weeknr))
				}

				// Add event text
				for _, ev := range eventList {
					if nd.Day() == ev.Day && nd.Month() == ev.Month {
						x, y := pdf.GetXY()
						pdf.SetFont(MAINFONT, "", eventFontsize)

						textArray := strings.Split(ev.Text, "\n")

						for i, j := range textArray {
							pdf.Text(x+0.02*cw, y+0.50*ch+float64(i)*eventFontsize/4.0, fmt.Sprintf("%s", j))
						}
					}
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
		pdf.Ln(-1)
		pdf.SetFont(MAINFONT, "B", 12)
		pdf.Text(0.50*PAGEWIDTH, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", *footer))
	}

	pdf.OutputAndClose(docWriter(pdf))
}
