package gocal

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

/*
This is gocal a tool to generate calendars in PDF for printing.

https://github.com/StefanSchroeder/Gocal
Copyright (c) 2014 Stefan Schroeder, NY, 2014-03-10

*/

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/soniakeys/meeus/v3/julian"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// LINES has to cover the number of weeks in one month.
	LINES = 6
	// COLUMNS equals the number of days in a week.
	COLUMNS = 7
	// MARGIN is the padding around the calendar on the page.
	MARGIN = 10.0 // MM
	// CELLMARGIN is the padding in the cell.
	CELLMARGIN = 1.0

	// DARKGREY is the intensity of grey in cell backgrounds.
	DARKGREY = 150
	// LIGHTGREY is the intensity of grey in neighbor month days.
	LIGHTGREY = 170
	// BLACK is black.
	BLACK = 0

	// MOONSIZE is the size of the moon icon.
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

var testedLanguage = map[string]bool{
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
	OptConfigs         []string
	OptPhotos          string
	OptFontScale       float64
	OptNocolor         bool
	EventList          []gDate
	OptCutWeekday      int
	OptCheckers        bool
	OptFillpattern     string
	OptYearSpread      int
}

func New(b int, e int, y int) *Calendar {
	return &Calendar{b, e, y,
		"serif", // OptFont
		"",      // OptFooter
		"L",     // OptOrientation P=portrait
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
		nil,     // OptConfigs
		"",      // OptPhotos
		1.0,     // OptFontScale
		false,   // OptNocolor
		nil,     // EventList
		0,       // OptCutWeekday
		false,   // OptCheckers
		"",      // OptFillpattern
		1,       // OptYearSpread
	}
}

// gDate is a type to store single events
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

// monthRange stores begin and end month of the year
type monthRange struct {
	begin int
	end   int
}

// myPdf is an anonymous struct that allows to define methods on non-local types
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
	pdf         *gofpdf.Fpdf
	fl          *os.File
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

func (g *Calendar) WantFillMode(s string) bool {
	if strings.Index(g.OptFillpattern, s) != -1 {
		return true
	}
	return false
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

func (g *Calendar) SetNocolor() {
	g.OptNocolor = true
}

func (g *Calendar) SetSmall() {
	g.OptSmall = true
}

func (g *Calendar) SetYearSpread(frac int) {
	switch frac {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 12:
		g.OptYearSpread = frac
	}
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

func (g *Calendar) AddConfig(f string) {
	g.OptConfigs = append(g.OptConfigs, f)
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

func (g *Calendar) SetFillpattern(f string) {
	g.OptFillpattern = f
}

func (g *Calendar) AddEvent(day int, month int, text string, image string) {
	gcd := gDate{time.Month(month), int(day), text, "", image}
	g.EventList = append(g.EventList, gcd)
}

func (g *Calendar) SetPaperformat(f string) {
	g.OptPaperformat = f
}

func (g *Calendar) WantFill(i int, j int, wd time.Weekday) bool {

	if wd == time.Sunday && g.WantFillMode("S") {
		return true
	}

	if wd == time.Saturday && g.WantFillMode("s") {
		return true
	}

	if i%2 == 0 && g.WantFillMode("Y") {
		return true
	}

	if (i+1)%2 == 0 && g.WantFillMode("y") {
		return true
	}

	if j%2 == 0 && g.WantFillMode("X") {
		return true
	}

	if (j+1)%2 == 0 && g.WantFillMode("x") {
		return true
	}

	if (j+i)%2 == 0 && g.WantFillMode("c") {
		return true
	}

	if (j+i+1)%2 == 0 && g.WantFillMode("C") {
		return true
	}

	return false
}

func getLanguage(inLanguage string) (outLanguage string) {
	// First try Environment
	outLanguage = os.Getenv("LANG")

	// If set on the cmdline, override
	if inLanguage != "" {
		outLanguage = inLanguage
	}

	// if we don't know that language, fall back to en.
	if testedLanguage[outLanguage] != true {
		outLanguage = "en_US"
	}
	return
}

func (g *Calendar) AddWallpaper(pdf *gofpdf.Fpdf, fontTempdir string, PAGEWIDTH float64, PAGEHEIGHT float64) {
	wallpaperFilename := g.OptWallpaper
	if strings.HasPrefix(wallpaperFilename, "http://") {
		wallpaperFilename = downloadFile(g.OptWallpaper, fontTempdir)
	}
	pdf.Image(wallpaperFilename, 0, 0, PAGEWIDTH, PAGEHEIGHT, false, "", 0, "")
}

func (g *Calendar) CreateYearCalendarInverse(fn string) {

	var fontTempdir string
	var fontScale = g.OptFontScale
	var calFont = g.OptFont

	if g.OptSmall == true {
		fontScale = 0.75
	}

	wantyear := g.WantYear

	calFont, fontTempdir = processFont(calFont)

	pdf := gofpdf.New(g.OptOrientation, "mm", g.OptPaperformat, fontTempdir)
	pdf.AddFont(calFont, "", calFont+".json")

	pdf.SetFillColor(LIGHTGREY, LIGHTGREY, LIGHTGREY)
	pdf.SetMargins(10.0, 5.0, 10.0)
	pdf.SetTitle("Created with Gocal", true)

	PAGEWIDTH, PAGEHEIGHT, _ := pdf.PageSize(0)
	if g.OptOrientation != "P" {
		PAGEWIDTH, PAGEHEIGHT = PAGEHEIGHT, PAGEWIDTH
	}

	cw := (PAGEWIDTH - 2*MARGIN) / 12.5
	ch := (PAGEHEIGHT - 2*MARGIN) / 32
	currentLanguage := getLanguage(g.OptLocale)
	localizedWeekdayNames := getLocalizedWeekdayNames(currentLanguage, 2)
	localizedMonthNames := getLocalizedMonthNames(currentLanguage)

	monthFracture := g.OptYearSpread
	cw = cw * float64(monthFracture)
	monthOnePage := 12 / monthFracture

	for pageCount := 0; pageCount < monthFracture; pageCount++ {
		pdf.AddPage()

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		pdf.SetFont(calFont, "", HEADERFONTSIZE*fontScale)
		pdf.CellFormat(PAGEWIDTH-MARGIN, MARGIN, fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")

		if g.OptWallpaper != "" {
			g.AddWallpaper(pdf, fontTempdir, PAGEWIDTH, PAGEHEIGHT)
		}

		pdf.Ln(-1)

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		pdf.CellFormat(cw*0.5/float64(monthFracture), ch*0.75, "", "1", 0, "C", false, 0, "")

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		pdf.SetFont(calFont, "", FOOTERFONTSIZE*fontScale*0.8)
		pdf.SetFillColor(LIGHTGREY, LIGHTGREY, LIGHTGREY)
		for mo := pageCount*monthOnePage + 1; mo <= pageCount*monthOnePage+monthOnePage; mo++ {
			pdf.CellFormat(cw, ch*0.75, fmt.Sprintf("%s", localizedMonthNames[mo]), "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale*0.25)
		for i := 1; i <= 31; i++ {
			pdf.SetTextColor(BLACK, BLACK, BLACK)
			pdf.CellFormat(cw*0.5/float64(monthFracture), ch*0.9, fmt.Sprintf("%d", i), "1", 0, "C", false, 0, "")
			for j := pageCount*monthOnePage + 1; j <= pageCount*monthOnePage+monthOnePage; j++ {
				tDay := time.Date(wantyear, time.Month(j), i, 0, 0, 0, 0, time.UTC)
				wd := localizedWeekdayNames[(tDay.Weekday()+1)%7]

				if (tDay.Weekday() == time.Saturday || tDay.Weekday() == time.Sunday) && !g.OptNocolor {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(BLACK, BLACK, BLACK)
				}

				_, readbackMonth, _ := tDay.Date()
				if int(readbackMonth) == int(j) {

					// Day of year, lower right
					if g.OptHideDOY == false && int(tDay.Month()) == j {
						doy := julian.DayOfYearGregorian(wantyear, int(time.Month(j)), int(tDay.Day()))
						pdf.SetFont(calFont, "", DOYFONTSIZE*fontScale*0.5)
						pdf.CellFormat(cw, ch*0.9, fmt.Sprintf("%d", doy), "1", 0, "BR", false, 0, "")
						pdf.SetX(pdf.GetX() - cw) // reset
					}
					// Add week number, lower left
					if tDay.Weekday() == time.Monday && g.OptHideWeek == false {
						pdf.SetFont(calFont, "", WEEKFONTSIZE*0.5*fontScale)
						_, weeknr := tDay.ISOWeek()
						pdf.CellFormat(cw, ch*0.9, fmt.Sprintf("W %d", weeknr), "1", 0, "BL", false, 0, "")
						pdf.SetX(pdf.GetX() - cw) // reset
					}

					fillBox := g.WantFill(i, j, tDay.Weekday())

					pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale*0.25)
					pdf.CellFormat(cw, ch*0.9, fmt.Sprintf("%s", wd), "1", 0, "TL", fillBox, 0, "")
				} else {
					// empty cell to skip ahead
					pdf.CellFormat(cw, ch*0.9, "", "1", 0, "TL", false, 0, "")
				}
			}
			pdf.Ln(-1)
		}

		pdf.Ln(-1)
		pdf.SetTextColor(DARKGREY, DARKGREY, DARKGREY)
		pdf.SetFont(calFont, "", FOOTERFONTSIZE*fontScale)
		pdf.Text(0.50*PAGEWIDTH-pdf.GetStringWidth(g.OptFooter)*0.5, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", g.OptFooter))
	}

	pdf.OutputAndClose(docWriter(pdf, fn))
	removeTempdir(fontTempdir)
}

func (g *Calendar) CreateYearCalendar(fn string) {

	var fontTempdir string
	var fontScale = g.OptFontScale
	var calFont = g.OptFont

	if g.OptSmall == true {
		fontScale = 0.75
	}

	wantyear := g.WantYear

	calFont, fontTempdir = processFont(calFont)

	pdf := gofpdf.New(g.OptOrientation, "mm", g.OptPaperformat, fontTempdir)
	pdf.AddFont(calFont, "", calFont+".json")

	pdf.SetFillColor(LIGHTGREY, LIGHTGREY, LIGHTGREY)
	pdf.SetMargins(10.0, 5.0, 10.0)
	pdf.SetTitle("Created with Gocal", true)

	PAGEWIDTH, PAGEHEIGHT, _ := pdf.PageSize(0)
	if g.OptOrientation != "P" {
		PAGEWIDTH, PAGEHEIGHT = PAGEHEIGHT, PAGEWIDTH
	}

	currentLanguage := getLanguage(g.OptLocale)
	localizedWeekdayNames := getLocalizedWeekdayNames(currentLanguage, 2)
	localizedMonthNames := getLocalizedMonthNames(currentLanguage)

	monthFracture := g.OptYearSpread
	monthOnePage := 12 / monthFracture

	cw := (PAGEWIDTH - 2*MARGIN) / 32
	ch := (PAGEHEIGHT - 2*MARGIN) / 14
	ch = ch * float64(monthFracture)
	for pageCount := 0; pageCount < monthFracture; pageCount++ {
		pdf.AddPage()
		pdf.SetTextColor(BLACK, BLACK, BLACK)

		if g.OptWallpaper != "" {
			g.AddWallpaper(pdf, fontTempdir, PAGEWIDTH, PAGEHEIGHT)
		}

		pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale)
		pdf.CellFormat(PAGEWIDTH-MARGIN, MARGIN, fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")
		pdf.Ln(-1)

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		monthTable := func(mymonth int, myyear int) {
			var day int64 = 1

			pdf.CellFormat(cw, ch, "", "1", 0, "C", false, 0, "")
			for j := 1; j < 32; j++ {
				pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale*0.25)

				tDay := time.Date(myyear, time.Month(mymonth), j, 0, 0, 0, 0, time.UTC)
				if (tDay.Weekday() == time.Saturday || tDay.Weekday() == time.Sunday) && !g.OptNocolor {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(BLACK, BLACK, BLACK)
				}
				// if the date is invalid, like 30.2., time.Date will still have a proper date
				// in this case e.g. the 1.3. We have to check if the month we put in is the
				// the month that arrived.
				_, readbackMonth, _ := tDay.Date()
				if int(readbackMonth) == mymonth {

					// Day of year, lower right
					if g.OptHideDOY == false && int(tDay.Month()) == mymonth && tDay.Weekday() != time.Monday {
						doy := julian.DayOfYearGregorian(wantyear, int(mymonth), int(tDay.Day()))
						pdf.SetFont(calFont, "", DOYFONTSIZE*fontScale*0.5)
						pdf.CellFormat(cw, ch, fmt.Sprintf("%d", doy), "1", 0, "BR", false, 0, "")
						pdf.SetX(pdf.GetX() - cw) // reset
					}
					// Add week number, lower left
					if tDay.Weekday() == time.Monday && g.OptHideWeek == false {
						pdf.SetFont(calFont, "", WEEKFONTSIZE*0.5*fontScale)
						_, weeknr := tDay.ISOWeek()
						pdf.CellFormat(cw, ch, fmt.Sprintf("W %d", weeknr), "1", 0, "BL", false, 0, "")
						pdf.SetX(pdf.GetX() - cw) // reset
					}

					fillBox := g.WantFill(mymonth, j, tDay.Weekday())

					pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale*0.25)
					pdf.CellFormat(cw, ch, fmt.Sprintf("%s", localizedWeekdayNames[(tDay.Weekday()+1)%7]), "1", 0, "TL", fillBox, 0, "")
					day++
				}
			}
		}

		var day int64 = 1

		// The header cells shall not scale with the monthFracture. Undo it.
		var ch_header = ch / float64(monthFracture) * 0.3
		pdf.CellFormat(cw, ch_header, "", "1", 0, "C", false, 0, "")

		// top row: 1..31
		for j := 0; j < 31; j++ {
			pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale*0.25)
			pdf.CellFormat(cw, ch_header, fmt.Sprintf("%d", day), "1", 0, "C", false, 0, "")
			day++
		}
		pdf.Ln(-1)

		//for mo := 1; mo <= totalMonth; mo++ {
		for mo := pageCount*monthOnePage + 1; mo <= pageCount*monthOnePage+monthOnePage; mo++ {
			pdf.SetTextColor(BLACK, BLACK, BLACK)
			pdf.SetFont(calFont, "", FOOTERFONTSIZE*fontScale*0.8)
			pdf.TransformBegin()
			x, y := pdf.GetXY()
			pdf.TransformRotate(90, x+cw-CELLMARGIN, y+ch-CELLMARGIN)
			pdf.Text(x+cw-CELLMARGIN, y+ch-CELLMARGIN*2, localizedMonthNames[mo])
			pdf.TransformEnd()
			monthTable(mo, wantyear)
			pdf.Ln(-1)
		}
		pdf.Ln(-1)
		pdf.SetTextColor(DARKGREY, DARKGREY, DARKGREY)
		pdf.SetFont(calFont, "", FOOTERFONTSIZE*fontScale)
		pdf.Text(0.50*PAGEWIDTH-pdf.GetStringWidth(g.OptFooter)*0.5, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", g.OptFooter))
	}

	pdf.OutputAndClose(docWriter(pdf, fn))
	removeTempdir(fontTempdir)
}

func getPhotolist(in string, temp string) (out [12]string) {
	if in != "" {
		for i := 0; i < 12; i++ {
			photoname := in
			if strings.HasPrefix(photoname, "http://") {
				photoname = downloadFile(photoname, temp)
			}
			out[i] = photoname
		}
	}
	return out
}

func getPhotoslist(in string) (out [12]string) {
	if in != "" {
		fileList, err := filepath.Glob(in + string(os.PathSeparator) + "*")
		if err == nil {
			for i := 0; i < 12; i++ {
				out[i] = fileList[i%len(fileList)]
			}
		} else {
			fmt.Printf("# There is an error in your path to photos: %v\n", err)
		}
	}
	return out
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

	currentLanguage := getLanguage(g.OptLocale)

	var fileEventList = make([]gDate, 10000) // Maximum number of events

	if g.OptConfig != "" {
		fileEventList = readConfigurationfile(g.OptConfig)
	}

	if len(g.OptConfigs) > 0 {
		for _, evfile := range g.OptConfigs {
			thiseventList := readConfigurationfile(evfile)
			for _, ev := range thiseventList {
				fileEventList = append(fileEventList, ev)
			}
		}
	}

	eventList := fileEventList
	for _, ev := range g.EventList {
		eventList = append(eventList, ev)
	}

	wantyear := g.WantYear
	wantmonths := monthRange{g.WantBeginMonth, g.WantEndMonth}
	localizedMonthNames := getLocalizedMonthNames(currentLanguage)
	localizedWeekdayNames := getLocalizedWeekdayNames(currentLanguage, 0)

	var calFont = g.OptFont

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
	photoList = getPhotolist(g.OptPhoto, fontTempdir)
	if g.OptPhotos != "" {
		photoList = getPhotoslist(g.OptPhotos)
	}
	if g.OptPhoto != "" || g.OptPhotos != "" {
		ch *= 0.5
	}

	calendarTable := func(mymonth int, myyear int) {
		pdf.SetFont(calFont, "", WEEKDAYFONTSIZE*fontScale)
		for weekday := 0; weekday <= 6; weekday++ { // Print weekdays in first row
			// The week row can be smaller
			pdf.CellFormat(cw, ch*0.33, localizedWeekdayNames[(weekday+2)%7], "0", 0, "C", false, 0, "")
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
				today := time.Date(myyear, time.Month(mymonth), 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(day) * 24 * 60 * 60 * time.Second)
				fill := g.WantFill(i, j, today.Weekday())

				// Determine color
				if today.Month() != time.Month(mymonth) { // GREY
					pdf.SetTextColor(DARKGREY, DARKGREY, DARKGREY)
					pdf.SetFillColor(LIGHTGREY, LIGHTGREY, LIGHTGREY)
					fill = false // FIXME, do we want fill here?
				} else if (today.Weekday() == time.Saturday || today.Weekday() == time.Sunday) && !g.OptNocolor {
					pdf.SetTextColor(255, 0, 0) // RED
				} else {
					pdf.SetTextColor(BLACK, BLACK, BLACK)
				}

				if g.OptHideOtherMonths == true && today.Month() != time.Month(mymonth) {
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
					if g.OptPhoto != "" || g.OptPhotos != "" {
						moonsize *= 0.6
					}
					myMoonPDF := myPdf{pdf, moonsize}
					pdf.SetFillColor(LIGHTGREY, LIGHTGREY, LIGHTGREY)
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
				if g.OptHideDOY == false && int(today.Month()) == mymonth {
					doy := julian.DayOfYearGregorian(myyear, mymonth, int(today.Day()))
					pdf.SetFont(calFont, "", DOYFONTSIZE*fontScale)
					pdf.CellFormat(cw, ch, fmt.Sprintf("%d", doy), "1", 0, "BR", fill, 0, "")
					pdf.SetX(pdf.GetX() - cw) // reset
				}

				// Add week number, lower left
				if today.Weekday() == time.Monday && g.OptHideWeek == false {
					pdf.SetFont(calFont, "", WEEKFONTSIZE*fontScale)
					_, weeknr := today.ISOWeek()
					pdf.CellFormat(cw, ch, fmt.Sprintf("W %d", weeknr), "1", 0, "BL", fill, 0, "")
					pdf.SetX(pdf.GetX() - cw) // reset
				}

				// Add event text
				for _, ev := range eventList {
					if today.Day() == ev.Day && today.Month() == ev.Month {
						x, y := pdf.GetXY()
						pdf.SetFont(calFont, "", EVENTFONTSIZE*fontScale)

						if ev.Image != "" {
							pdf.Image(ev.Image, x, y, cw, ch, false, "", 0, "")
						}
						for i, j := range strings.Split(ev.Text, "\\n") {
							pdf.Text(x+0.02*cw, y+0.50*ch+float64(i)*EVENTFONTSIZE*fontScale/3.0, fmt.Sprintf("%s", j))
						}
					}
				}

				// day of the month, big number
				pdf.SetFont(calFont, "", MONTHDAYFONTSIZE*fontScale)
				pdf.CellFormat(cw, ch, fmt.Sprintf("%d", today.Day()), "1", 0, "TL", fill, 0, "")
				day++
			}
			pdf.Ln(-1)
		}
	}

	for mo := wantmonths.begin; mo <= wantmonths.end; mo++ {
		//fmt.Printf("Printing page %d\n", page)
		pdf.AddPage()
		if g.OptWallpaper != "" {
			g.AddWallpaper(pdf, fontTempdir, PAGEWIDTH, PAGEHEIGHT)
		}

		if g.OptPhoto != "" || g.OptPhotos != "" {
			photo := photoList[mo-1] // this list is zero-based.
			if photo != "" {
				pdf.Image(photo, 0, PAGEHEIGHT*0.5, PAGEWIDTH, PAGEHEIGHT*0.5, false, "", 0, "")
			}
		}

		pdf.SetTextColor(BLACK, BLACK, BLACK)
		pdf.SetFont(calFont, "", HEADERFONTSIZE*fontScale)
		pdf.CellFormat(PAGEWIDTH-MARGIN, MARGIN, localizedMonthNames[mo]+" "+fmt.Sprintf("%d", wantyear), "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		calendarTable(mo, wantyear)

		pdf.Ln(-1)
		pdf.SetTextColor(DARKGREY, DARKGREY, DARKGREY)
		pdf.SetFont(calFont, "", FOOTERFONTSIZE*fontScale)
		pdf.Text(0.50*PAGEWIDTH-pdf.GetStringWidth(g.OptFooter)*0.5, 0.95*PAGEHEIGHT, fmt.Sprintf("%s", g.OptFooter))
	}
	pdf.OutputAndClose(docWriter(pdf, fn))
	removeTempdir(fontTempdir)
}
