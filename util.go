package gocal

// Copyright (c) 2014 Stefan Schroeder, NY, 2014-03-10
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
//
// util.go
//
// This file is part of gocal, a PDF calendar generator in Go.
//
// https://github.com/StefanSchroeder/Gocal
//

import _ "embed"


import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/PuloV/ics-golang"
	"github.com/goodsign/monday"
	"github.com/phpdave11/gofpdf"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"github.com/soniakeys/meeus/v3/julian"
	"github.com/soniakeys/meeus/v3/moonphase"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const YmdHis = "2006-01-02 15:04:05"

// TelegramStore is a container to read XML event-list
type TelegramStore struct {
	XMLName   xml.Name `xml:"Gocal"`
	Gocaldate []Gocaldate
}

const (
	CLEARTEMP = true
)

// removeTempdir removes the temoprary directory,
// unless we want to keep it for debugging.
func removeTempdir(d string) {
	if CLEARTEMP == false {
		return
	}
	os.RemoveAll(d)
}

// computeMoonphasesJ populates a map for the entire year.
// Keys are dates in YYYY-MM-DD format,
// Values are strings from the list Full, New, First, Last.
func computeMoonphasesJ(moonJ map[string]string, yr int) {
	daysInYear := 365
	if julian.LeapYearGregorian(yr) {
		daysInYear = 366
	}

	moon_funcs := map[string]func(float64) float64{
		"Full":  moonphase.Full,
		"New":   moonphase.New,
		"First": moonphase.First,
		"Last":  moonphase.Last,
	}

	// For each days of the year we compute the
	// nearest Full/New/First/Last moonphase.
	// The dates for each are crammed into the result map.
	for i := 0; i < daysInYear; i++ {
		decimalYear := float64(yr) +
			float64(i-1)/float64(daysInYear)
		for moonkey, _ := range moon_funcs {
			jd := moon_funcs[moonkey](decimalYear)
			y, m, d := julian.JDToCalendar(jd)
			moonString := fmt.Sprintf("%04d-%02d-%02d", y, m, int(d))
			moonJ[moonString] = moonkey
		}
	}
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
			fmt.Printf("New moon on %d.%d.%d\n", int(d), int(m), int(y))
			moon[int(d)] = "New"
		}
		jdeNew = moonphase.Full(decimalYear)
		y, m, d = julian.JDToCalendar(jdeNew)
		if (y == yr) && (m == mo) && (int(d) == i) {
			fmt.Printf("Full moon on %d.%d.%d\n", int(d), int(m), int(y))
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

//go:embed fonts/FreeSansBold.ttf
var freesansbold []byte

//go:embed fonts/FreeMonoBold.ttf
var freemonobold []byte

//go:embed fonts/FreeSerifBold.ttf
var freeserifbold []byte

// processFont creates a font usable from a TTF.
// It also sets up the temporary directory to store the
// intermediate files.
func processFont(fontFile string) (fontName, tempDirname string) {
	var err error
	tempDirname, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	if fontFile == "mono" {
		fontFile = tempDirname + string(os.PathSeparator) + "freemonobold.ttf"
		ioutil.WriteFile(fontFile, freemonobold, 0700)
	} else if fontFile == "serif" {
		fontFile = tempDirname + string(os.PathSeparator) + "freeserifbold.ttf"
		ioutil.WriteFile(fontFile, freeserifbold, 0700)
	} else if fontFile == "sans" {
		fontFile = tempDirname + string(os.PathSeparator) + "freesansbold.ttf"
		ioutil.WriteFile(fontFile, freesansbold, 0700)
	}
	err = ioutil.WriteFile(tempDirname+string(os.PathSeparator)+"cp1252.map", []byte(codepageCP1252), 0700)
	if err != nil {
		log.Fatal(err)
	}
	err = gofpdf.MakeFont(fontFile, tempDirname+string(os.PathSeparator)+"cp1252.map", tempDirname, nil, true)
	if err != nil {
		log.Fatal(err)
	}
	fontName = filepath.Base(fontFile)
	fontName = strings.TrimSuffix(fontName, filepath.Ext(fontName))
	// fmt.Printf("Using external font: %v\n", fontName)
	return fontName, tempDirname
}

// downloadFile loads a file via http into the tempDir
// and returns the fullpath filename.
func downloadFile(in string, tempDir string) (fileName string) {
	extension := filepath.Ext(in)

	// The filename from the URL might contain colons that are
	// not valid characters in a filename in Windows. THerefore
	// we simply call out image 'image'.
	fileName = "image" + extension
	fileName = tempDir + string(os.PathSeparator) + fileName

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("# Error creating %v\n", output)
		return
	}
	defer output.Close()

	retrieve, err := http.Get(in)
	if err != nil {
		fmt.Printf("# Error downloading %v\n", in)
		return
	}
	defer retrieve.Body.Close()

	_, err = io.Copy(output, retrieve.Body)
	if err != nil {
		fmt.Printf("# Error copying %v\n", in)
		return
	}

	return fileName
}

// This function converts a string into the required
// Codepage.
func convertCP(in string) (out string) {
	buf := new(bytes.Buffer)
	w, err := charset.NewWriter("windows-1252", buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, in)
	w.Close()

	out = fmt.Sprintf("%s", buf)
	return out
}

// This function reads the events XML file and returns a
// list of gDate objects.
func readICSfile(filename string, targetyear int) (eL []gDate) {

	/* There is an ugly hack lurking here. The events in ICS
	contain years, but we wanted the configuration to be
	agnostic of years.*/
	parser := ics.New()

	ics.FilePath = "tmp/new/"

	ics.DeleteTempFiles = true

	inputChan := parser.GetInputChan()

	outputChan := parser.GetOutputChan()

	inputChan <- filename

	go func() {
		for event := range outputChan {
			eventText := convertCP(event.GetSummary())
			year := event.GetStart().Format("2006")
			mon := event.GetStart().Format("01")
			day := event.GetStart().Format("02")

			yr, _ := strconv.ParseInt(year, 10, 32)
			mo, _ := strconv.ParseInt(mon, 10, 32)
			d, _ := strconv.ParseInt(day, 10, 32)
			if int(targetyear) == int(yr) {
				gcd := gDate{time.Month(mo), int(d), eventText, "", ""}
				eL = append(eL, gcd)
			}
		}
	}()
	parser.Wait()

	return eL
}

// This function reads the events XML file and returns a
// list of gDate objects.
func readConfigurationfile(filename string) (eL []gDate) {

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
		log.Fatalf("# ERROR: when trying to unmarshal the XML configuration file: %v", err2)
		return
	}

	for _, m := range v.Gocaldate {

		if strings.Index(m.Date, "/") != -1 { // Is this Month/Day ?

			textArray := strings.Split(m.Date, "/")

			eventText := convertCP(m.Text)

			if textArray[0] == "*" {
				d, _ := strconv.ParseInt(textArray[1], 10, 32)
				for j := 1; j < 13; j++ {
					gcd := gDate{time.Month(j), int(d), eventText, "", m.Image}
					eL = append(eL, gcd)
				}
			} else {
				mo, _ := strconv.ParseInt(textArray[0], 10, 32)
				d, _ := strconv.ParseInt(textArray[1], 10, 32)

				gcd := gDate{time.Month(mo), int(d), eventText, "", m.Image}
				eL = append(eL, gcd)
			}
		} else { // There is no slash, assume weekday

			eventText := convertCP(m.Text)
			gcd := gDate{time.Month(0), int(0), eventText, string(m.Date), m.Image}
			eL = append(eL, gcd)
		}
	}

	return eL
}

// / This function returns an array of Monthnames already in the
// right locale.
func getLocalizedMonthNames(locale string) (monthnames [13]string) {

	for page := 1; page < 13; page++ {
		t := time.Date(2013, time.Month(page), 1, 0, 0, 0, 0, time.UTC)
		monthnames[page] = convertCP(fmt.Sprintf("%s", monday.Format(t, "January", monday.Locale(locale))))
	}

	return monthnames
}

// / This function returns an array of weekday names already in the
// right locale.
func getLocalizedWeekdayNames(locale string, cutoff int) (wdnames [8]string) {
	for i := 0; i <= 6; i++ {
		// Some arbitrary date, that allows us to pickup Weekday-Strings.
		t := time.Date(2013, 1, 5+i, 0, 0, 0, 0, time.UTC)
		wdnames[i] = convertCP(monday.Format(t, "Monday", monday.Locale(locale)))
		if cutoff > 0 {
			wdnames[i] = wdnames[i][0:cutoff]
		}
	}
	return wdnames
}

const codepageCP1252 = `!00 U+0000 .notdef
!01 U+0001 .notdef
!02 U+0002 .notdef
!03 U+0003 .notdef
!04 U+0004 .notdef
!05 U+0005 .notdef
!06 U+0006 .notdef
!07 U+0007 .notdef
!08 U+0008 .notdef
!09 U+0009 .notdef
!0A U+000A .notdef
!0B U+000B .notdef
!0C U+000C .notdef
!0D U+000D .notdef
!0E U+000E .notdef
!0F U+000F .notdef
!10 U+0010 .notdef
!11 U+0011 .notdef
!12 U+0012 .notdef
!13 U+0013 .notdef
!14 U+0014 .notdef
!15 U+0015 .notdef
!16 U+0016 .notdef
!17 U+0017 .notdef
!18 U+0018 .notdef
!19 U+0019 .notdef
!1A U+001A .notdef
!1B U+001B .notdef
!1C U+001C .notdef
!1D U+001D .notdef
!1E U+001E .notdef
!1F U+001F .notdef
!20 U+0020 space
!21 U+0021 exclam
!22 U+0022 quotedbl
!23 U+0023 numbersign
!24 U+0024 dollar
!25 U+0025 percent
!26 U+0026 ampersand
!27 U+0027 quotesingle
!28 U+0028 parenleft
!29 U+0029 parenright
!2A U+002A asterisk
!2B U+002B plus
!2C U+002C comma
!2D U+002D hyphen
!2E U+002E period
!2F U+002F slash
!30 U+0030 zero
!31 U+0031 one
!32 U+0032 two
!33 U+0033 three
!34 U+0034 four
!35 U+0035 five
!36 U+0036 six
!37 U+0037 seven
!38 U+0038 eight
!39 U+0039 nine
!3A U+003A colon
!3B U+003B semicolon
!3C U+003C less
!3D U+003D equal
!3E U+003E greater
!3F U+003F question
!40 U+0040 at
!41 U+0041 A
!42 U+0042 B
!43 U+0043 C
!44 U+0044 D
!45 U+0045 E
!46 U+0046 F
!47 U+0047 G
!48 U+0048 H
!49 U+0049 I
!4A U+004A J
!4B U+004B K
!4C U+004C L
!4D U+004D M
!4E U+004E N
!4F U+004F O
!50 U+0050 P
!51 U+0051 Q
!52 U+0052 R
!53 U+0053 S
!54 U+0054 T
!55 U+0055 U
!56 U+0056 V
!57 U+0057 W
!58 U+0058 X
!59 U+0059 Y
!5A U+005A Z
!5B U+005B bracketleft
!5C U+005C backslash
!5D U+005D bracketright
!5E U+005E asciicircum
!5F U+005F underscore
!60 U+0060 grave
!61 U+0061 a
!62 U+0062 b
!63 U+0063 c
!64 U+0064 d
!65 U+0065 e
!66 U+0066 f
!67 U+0067 g
!68 U+0068 h
!69 U+0069 i
!6A U+006A j
!6B U+006B k
!6C U+006C l
!6D U+006D m
!6E U+006E n
!6F U+006F o
!70 U+0070 p
!71 U+0071 q
!72 U+0072 r
!73 U+0073 s
!74 U+0074 t
!75 U+0075 u
!76 U+0076 v
!77 U+0077 w
!78 U+0078 x
!79 U+0079 y
!7A U+007A z
!7B U+007B braceleft
!7C U+007C bar
!7D U+007D braceright
!7E U+007E asciitilde
!7F U+007F .notdef
!80 U+20AC Euro
!82 U+201A quotesinglbase
!83 U+0192 florin
!84 U+201E quotedblbase
!85 U+2026 ellipsis
!86 U+2020 dagger
!87 U+2021 daggerdbl
!88 U+02C6 circumflex
!89 U+2030 perthousand
!8A U+0160 Scaron
!8B U+2039 guilsinglleft
!8C U+0152 OE
!8E U+017D Zcaron
!91 U+2018 quoteleft
!92 U+2019 quoteright
!93 U+201C quotedblleft
!94 U+201D quotedblright
!95 U+2022 bullet
!96 U+2013 endash
!97 U+2014 emdash
!98 U+02DC tilde
!99 U+2122 trademark
!9A U+0161 scaron
!9B U+203A guilsinglright
!9C U+0153 oe
!9E U+017E zcaron
!9F U+0178 Ydieresis
!A0 U+00A0 space
!A1 U+00A1 exclamdown
!A2 U+00A2 cent
!A3 U+00A3 sterling
!A4 U+00A4 currency
!A5 U+00A5 yen
!A6 U+00A6 brokenbar
!A7 U+00A7 section
!A8 U+00A8 dieresis
!A9 U+00A9 copyright
!AA U+00AA ordfeminine
!AB U+00AB guillemotleft
!AC U+00AC logicalnot
!AD U+00AD hyphen
!AE U+00AE registered
!AF U+00AF macron
!B0 U+00B0 degree
!B1 U+00B1 plusminus
!B2 U+00B2 twosuperior
!B3 U+00B3 threesuperior
!B4 U+00B4 acute
!B5 U+00B5 mu
!B6 U+00B6 paragraph
!B7 U+00B7 periodcentered
!B8 U+00B8 cedilla
!B9 U+00B9 onesuperior
!BA U+00BA ordmasculine
!BB U+00BB guillemotright
!BC U+00BC onequarter
!BD U+00BD onehalf
!BE U+00BE threequarters
!BF U+00BF questiondown
!C0 U+00C0 Agrave
!C1 U+00C1 Aacute
!C2 U+00C2 Acircumflex
!C3 U+00C3 Atilde
!C4 U+00C4 Adieresis
!C5 U+00C5 Aring
!C6 U+00C6 AE
!C7 U+00C7 Ccedilla
!C8 U+00C8 Egrave
!C9 U+00C9 Eacute
!CA U+00CA Ecircumflex
!CB U+00CB Edieresis
!CC U+00CC Igrave
!CD U+00CD Iacute
!CE U+00CE Icircumflex
!CF U+00CF Idieresis
!D0 U+00D0 Eth
!D1 U+00D1 Ntilde
!D2 U+00D2 Ograve
!D3 U+00D3 Oacute
!D4 U+00D4 Ocircumflex
!D5 U+00D5 Otilde
!D6 U+00D6 Odieresis
!D7 U+00D7 multiply
!D8 U+00D8 Oslash
!D9 U+00D9 Ugrave
!DA U+00DA Uacute
!DB U+00DB Ucircumflex
!DC U+00DC Udieresis
!DD U+00DD Yacute
!DE U+00DE Thorn
!DF U+00DF germandbls
!E0 U+00E0 agrave
!E1 U+00E1 aacute
!E2 U+00E2 acircumflex
!E3 U+00E3 atilde
!E4 U+00E4 adieresis
!E5 U+00E5 aring
!E6 U+00E6 ae
!E7 U+00E7 ccedilla
!E8 U+00E8 egrave
!E9 U+00E9 eacute
!EA U+00EA ecircumflex
!EB U+00EB edieresis
!EC U+00EC igrave
!ED U+00ED iacute
!EE U+00EE icircumflex
!EF U+00EF idieresis
!F0 U+00F0 eth
!F1 U+00F1 ntilde
!F2 U+00F2 ograve
!F3 U+00F3 oacute
!F4 U+00F4 ocircumflex
!F5 U+00F5 otilde
!F6 U+00F6 odieresis
!F7 U+00F7 divide
!F8 U+00F8 oslash
!F9 U+00F9 ugrave
!FA U+00FA uacute
!FB U+00FB ucircumflex
!FC U+00FC udieresis
!FD U+00FD yacute
!FE U+00FE thorn
!FF U+00FF ydieresis
`
