/*

This is gocal a tool to generate calendars in PDF for printing.

* Inspired by pcal
* Simplicity: create a nice calendar programmatically with minimum effort.
* One argument: Year
* Two argument: Month Year
* Two argument: MonthBegin MonthEnd Year
* Week number
* Moonphase
* Add events from configuration file
* Set papersize
* Choose fonts (limited)
* Several languages
* day of year/remaining
* bg image
* Photo calendar
*
* add month notes
*/

package main

import (
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"code.google.com/p/gofpdf"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/goodsign/monday"
	"github.com/soniakeys/meeus/base"
	"github.com/soniakeys/meeus/julian"
	"github.com/soniakeys/meeus/moonillum"
	"github.com/soniakeys/meeus/moonphase"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	CONFIGFILE    = "gocal.xml"
	LINES         = 6
	cnGofpdfDir   = "."
	COLUMNS       = 7
	darkgrey      = 150
	lightgrey     = 170
	black         = 0
	MARGIN        = 10.0
	eventFontsize = 10
	loc           = "de_DE"
)

var optOrientation = flag.String("p", "L", "Orientation (L)andscape/(P)ortrait")
var outfilename = flag.String("o", "output.pdf", "Output filename")
var optFooter = flag.String("f", "Gocal", "Footer note")
var optLocale = flag.String("lang", "", "Language")
var optFont = flag.String("font", "Cabalett", "Main font (Times/Arial/Courier)")
var optHideWeek = flag.Bool("noweek", false, "Hide week number (false)")
var optHideEvents = flag.Bool("noevents", false, "Hide events from config file (false)")
var optHideMoon = flag.Bool("nomoon", false, "Hide moon phases (false)")
var optHideDOY = flag.Bool("nodoy", false, "Hide day of year (false)")
var optWallpaper = flag.String("wall", "", "Show wallpaper PNG JPG GIF")
var optPhoto = flag.String("photo", "", "Show photo (single image PNG JPG GIF)")
var optPhotos = flag.String("photos", "", "Show photos (directory PNG JPG GIF)")
var optPaper = flag.String("paper", "A4", "Paper format (A3 A4 A5 Letter Legal)")
var optNewfont = flag.String("newfont", "", "Nothing here.")
var moonSize = 4.0
var photoList [13]string
var calFont string
var	cnFontDir string = cnGofpdfDir + "/font"
var fontTempdir string

type gocalDate struct {
	Month   time.Month
	Day     int
	Text    string
	Weekday string
}

type monthRange struct {
	begin int
	end   int
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
			pdf.SetErrorf("Error opening output file.")
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

/// This function returns an array of Monthnames already in the
// right locale.
func getLocalizedMonthNames(locale string) (monthnames [13]string) {

	for page := 1; page < 13; page++ {
		t := time.Date(2013, time.Month(page), 1, 0, 0, 0, 0, time.UTC)
		buf := new(bytes.Buffer)
		w, err := charset.NewWriter("windows-1252", buf)
		if err != nil {
			log.Fatal(err)
		}

		monthString := fmt.Sprintf("%s", monday.Format(t, "January", monday.Locale(locale)))
		fmt.Fprintf(w, monthString)
		w.Close()

		monthnames[page] = fmt.Sprintf("%s", buf)
	}

	return monthnames
}

func getLocalizedWeekdayNames(locale string) (wdnames [8]string) {
	for i := 1; i <= 7; i++ {
		// Some arbitrary date, that allows us to pickup Weekday-Strings.
		t := time.Date(2013, 1, 6+i, 0, 0, 0, 0, time.UTC)
		wdnames[i] = monday.Format(t, "Monday", monday.Locale(locale))
	}
	return wdnames
}

func computeMoonphases(moon map[int]string, da int, mo int, yr int) {
	daysInYear := 365
	if julian.LeapYearGregorian(yr) {
		daysInYear = 366
	}
	// This is a perfect example for the difference btw. efficient and
	// effective. It's effective, but not efficient:
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

func main() {
	flag.Parse()

	testedLanguage := map[string]bool{
		"en_US": true,
		"en_GB": true,
		"da_DK": false,
		"nl_BE": true,
		"nl_NL": true,
		"fi_FI": true,
		"fr_FR": true,
		"fr_CA": true,
		"de_DE": true,
		"hu_HU": false,
		"it_IT": false,
		"nn_NO": false,
		"nb_NO": false,
		"pt_PT": false,
		"pt_BR": false,
		"ro_RO": false,
		"ru_RU": false,
		"es_ES": false,
		"sv_SE": false,
		"tr_TR": false,
		"bg_BG": false,
		"zh_CN": false,
		"zh_TW": false,
		"zh_HK": false,
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

	eventList := make([]gocalDate, 1000)

	if *optHideEvents == false {
		eventList = readConfigurationfile(CONFIGFILE)
	}

	var wantyear int = int(time.Now().Year())
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

	paperformat := *optPaper

  if *optNewfont != "" {
    var err error
    fontTempdir, err = ioutil.TempDir("", "")
    err = ioutil.WriteFile(fontTempdir + string(os.PathSeparator) + "cp1252.map", []byte(codepageCP1252), 0700)
    err = gofpdf.MakeFont(*optNewfont, fontTempdir + string(os.PathSeparator) + "cp1252.map", fontTempdir, os.Stderr, true)
    _ = err
    // FIXME Do some error checking here.
    cnFontDir = fontTempdir
    calFont = filepath.Base(*optNewfont)
    calFont = strings.TrimSuffix(calFont , filepath.Ext(calFont))
    fmt.Printf("Using external font: %v\n", calFont)
  } else {
    calFont = *optFont
  }

	pdf := gofpdf.New(*optOrientation, "mm", paperformat, cnFontDir)
	pdf.SetTitle("Created with Gocal", true)
  if calFont != "times" && calFont != "arial" && calFont != "courier" && calFont != "helvetica" {
    pdf.AddFont(calFont, "", calFont + ".json")
  }

	PAGEWIDTH, PAGEHEIGHT, _ := pdf.PageSize(0)
	if *optOrientation != "P" {
		PAGEWIDTH, PAGEHEIGHT = PAGEHEIGHT, PAGEWIDTH
	}

	cw := (PAGEWIDTH - 2*MARGIN) / COLUMNS // cellwidth w 20mm margin
	ch := PAGEHEIGHT / (LINES + 2)         // cellheight

	YOFFSET := 0.0
	if *optPhoto != "" {
		ch *= 0.5
		moonSize *= 0.6 // make moon smaller on photopage
		for i := 0; i < 13; i++ {
			photoList[i] = *optPhoto
		}
	}
	if *optPhotos != "" {
		ch *= 0.5
		moonSize *= 0.6 // make moon smaller on photopage
		a, _ := filepath.Glob(*optPhotos + "/*")
		for i := 0; i < 13; i++ {
			photoList[i] = a[i%len(a)]
		}
	}

	calendarTable := func(mymonth int, myyear int) {
		for weekday := 1; weekday < 8; weekday++ { // Print Weekdays
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
		pdf.SetFont(calFont, "", 24)
		for i := 0; i < LINES; i++ {
			for j := 0; j < COLUMNS; j++ {
				var fill bool = false
				nd := time.Date(myyear, time.Month(mymonth), 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(day) * 24 * 60 * 60 * time.Second)

				// Determine color
				if nd.Month() != time.Month(mymonth) { // GREY
					pdf.SetTextColor(darkgrey, darkgrey, darkgrey)
					pdf.SetFillColor(lightgrey, lightgrey, lightgrey)
					fill = false // FIXME
				} else if nd.Weekday() == time.Saturday || nd.Weekday() == time.Sunday {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(black, black, black) // BLACK
				}

				pdf.SetCellMargin(1)

				// Add moon icon
				if m, ok := moon[int(day)]; ok == true {
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

				// Day of year, lower right
				if *optHideDOY == false && int(nd.Month()) == mymonth {
					doy := julian.DayOfYearGregorian(myyear, mymonth, int(nd.Day()))
					x, y := pdf.GetXY()
					pdf.SetFont(calFont, "", 12)
					pdf.Text(x+0.82*cw, y+0.95*ch, fmt.Sprintf("%d", doy))
				}

				// Add week number, lower left
				if nd.Weekday() == time.Monday && *optHideWeek == false {
					x, y := pdf.GetXY()
					pdf.SetFont(calFont, "", 12)
					_, weeknr := nd.ISOWeek()
					pdf.Text(x+0.02*cw, y+0.95*ch, fmt.Sprintf("W %d", weeknr))
				}

				// Add event text
				for _, ev := range eventList {
					if nd.Day() == ev.Day && nd.Month() == ev.Month {
						x, y := pdf.GetXY()
						pdf.SetFont(calFont, "", eventFontsize)

						textArray := strings.Split(ev.Text, "\n")

						for i, j := range textArray {
							pdf.Text(x+0.02*cw, y+0.70*ch+float64(i)*eventFontsize/4.0, fmt.Sprintf("%s", j))
						}
					}
				}

				// day of the month, big number
				pdf.SetFont(calFont, "", 32)
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
			pdf.Image(*optWallpaper, 0, 0, PAGEWIDTH, PAGEHEIGHT, false, "", 0, "")
		}
		if *optPhoto != "" || *optPhotos != "" {
			pdf.Image(photoList[mo-1], 0, PAGEHEIGHT*0.5, cw*8, ch*8, false, "", 0, "")
		}
		pdf.SetTextColor(black, black, black)
		pdf.SetFont(calFont, "", 32)

		pdf.CellFormat(PAGEWIDTH-10.0, YOFFSET+10, localizedMonthNames[mo]+" "+fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetFont(calFont, "", 16)
		calendarTable(mo, wantyear)
		pdf.Ln(-1)
		pdf.SetFont(calFont, "", 12)
		pdf.Text(0.50*PAGEWIDTH, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", calFont))
	}

	pdf.OutputAndClose(docWriter(pdf))
}
