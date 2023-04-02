![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/01-2025.png?raw=true)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/StefanSchroeder/Gocal?status.png)](https://godoc.org/github.com/StefanSchroeder/Gocal)
[![Go Report Card](https://goreportcard.com/badge/github.com/StefanSchroeder/Gocal)](https://goreportcard.com/report/github.com/StefanSchroeder/Gocal)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/StefanSchroeder/Gocal/badge)](https://api.securityscorecards.dev/projects/github.com/StefanSchroeder/Gocal)


Gocal
=====

Gocal is a simple clone of pcal. It's a standalone tool and a library to create
monthly calendars in PDF with a no-nonsense attitude. By default it creates a
12-page PDF with one month per page for the current year.  
The library is called gocal, while the standalone tool is named gocalendar.

Gocalendar can be built with 'go build' in the gocalendar folder.

The following arguments are supported:

	gocalendar 2024 # Create a 12-page calendar for YEAR, e.g. 2024

	gocalendar 5 2025 # Create a 1-page calendar for the MONTH in YEAR, e.g. 5 2025

	gocalendar 5 7 2026 # Create a sequence from BEGIN to END in YEAR

There is also a year mode that shows the entire year on one page.
Have a look at the examples below to get an idea of gocal's capabilities.

Gocal was built because of pcal's lack of support for UTF-8 and because I feel 
that Postscript is obsolete.


Features
========

* PDF generation
* Week number on every Monday
* Day of year 
* Moon phases
* Events in XML configuration file
* Several languages 
* Wallpaper option
* Photo calendar option (from single image or directory)
* Page orientation and paper size option
* Font selection
* Year calendar (two layouts)
* Import of ICS files (local file or URL)


The main design goal of gocal is simplicity. While it is absolutely possible to create
an application where every single stroke is configurable, I believe that most
of you are too busy to play around with lots of options and want a pleasant
calendar out-of-the-box. That's what gocal provides.

The power of gocal is based on the cool libraries that it includes. This
implies that several of the options are actually options of the libraries, e.g.
the paper format (gofpdf) or the supported languages (goodsign/monday).

It is suggested to hide some of the optional fields or the cells will look
crowded.


Build instructions
==================

Run 

	go install github.com/StefanSchroeder/Gocal/gocalendar@latest

This will install the *gocalendar* executable in your go-bin folder.

Or the long haul:

Run 

	go get github.com/StefanSchroeder/Gocal

Gocal has a quite a few dependencies that go should resolve automatically.

To build the standalone tool, change into the gocalendar directory and run

	go build 

You create a bunch of sample files (and in passing test gocal) by running

	go test


Example library use
===================


    package main
    import (
      "github.com/StefanSchroeder/gocal"
    )
    func main() {
      g := gocal.New(1,12,2010)
      g.CreateCalendar("test-example01.pdf")
    }

License
=======

The license is in the LICENSE file. (It's MIT.)

API-documentation
=================

For the API documentation of gocal the library [visit the auto-generated docs on
godoc.org](https://pkg.go.dev/github.com/StefanSchroeder/Gocal).

Options of gocalendar
=====================

### Help

		-h  Help: Summarizes the options.

### Footer

		-f="Gocal": Footer note

Change the string at the bottom of the page. To disable the footer, simply set
it to the empty string. 

### Font

		-font="": font

		-font serif    (Times lookalike)
		-font mono     (Courier lookalike)
		-font sans     (Arial lookalike)
		-font path/to/font.ttf    (your favorite font)

Set the font. Only truetype fonts are supported. Look into c:\Windows\Fonts on
Windows and /usr/share/fonts to see what fonts you have. Gocal comes with three
fonts built-in: The Gnu FreeMonoBold, FreeSerifBold and FreeSansBold. They look
pretty similar to (in that order) Courier, Times and Arial and should meet all
your standard font needs.  These fonts are licensed under the Gnu FreeFont
License which accompanies this README.txt.  Read more about them at
https://www.gnu.org/software/freefont.  Auxiliary files are created in a
temporary directory.

The Borel-font in the sample directory is licensed under the OPL
and contained here only for illustrative purposes. 
It was downloaded from https://fontesk.com/borel-font/.

In addition you can provide your own TTF on the commandline if you prefer something fancy.

### Font size

Font sizes relative to the default size can be set with 

    -fontscale floatingpoint number

The default is 1.0 and typically you shouldn't need to change the font size
drastically. But since weekday names in some languages might be a lot longer
than in other languages, ugly collisions may occur.  To avoid that you can
rescale the fonts a little by setting -fontscale to 0.9.  There is a shortcut
to reduce the fontsizes globally to 75% of the default size to gain more room
for manual notes.

    -small


### Language

		-lang="": Language

Gocal reads the LANG environment variable. If it matches one of 

		en_US en_GB nl_BE nl_NL fi_FI fr_FR fr_CA de_DE

the library goodsign/monday is used to translate the weekday names and month
names.  Although this library supports a few other languages, I found that many
of the languages do not work with the fonts I tried. The language from the
environment can be overridden with this parameter. If your LANG is not
recognized, we default to en_US.

### Hiding stuff

		-nodoy: Hide day of year

		-nomoon: Hide moon phases

		-noweek: Hide week number

The week number according to ISO-8601 is added on every Monday by default.

		-noother: Hide neighbormonth days

This option hides the leading and trailing days of the neighbor months.
By default these days are printed in light grey.

		-plain This will hide everything that can be hidden (but not neighbormonth days).

### Output

		-o="output.pdf": Output filename

### Paper orientation

		-p="L": Orientation (L)andscape/(P)ortrait

Typically you want landscape for calendars without image and portrait for calendars with image.

### Paper format

		-paper="A4": Paper format (A3 A4 A5 Letter Legal)


### Graying out

    -fill configure filled grid

Can be any combination of:

    X even columns filled gray

    x odd columns filled gray

    Y even rows filled gray

    y odd rows filled gray

    S Sundays filled gray

    s Saturdays filled gray
  
    c checkerboard odd

    C checkerboard even

    1 Monday

    2 Tuesday

    3 Wednesday

    4 Thursday

    5 Friday

    6 Saturday

    7 Sunday

Example:

    -fill "x3Y"

The gray boxes are not transparent; therefore it doesn't make a lot
of sense to combine gray boxes with a wallpaper image.


### Photo / Photos / Wallpaper

		-photo=filename: Show single photo (single image in PNG JPG GIF)

This option will add this image to every month.
The filename can be a URL, qualified with http:// and it must have a valid image extension.

		-photos=directory: Show multiple photos (directory with PNG JPG GIF)

e.g. gocal -photos images/

This option will add the twelve first files as images to the twelve month.  If
less than twelve files are found, the sequence will re-start after the last
image.  This will not work if there are non-image files in the directory (among
the first twelve).  The directory option does NOT support URLs.

		-wall=filename: Show wallpaper PNG JPG GIF

e.g. gocal -wall gopher.png

This option will add this image to every month as a background. You should only
use images with a bright tone so that you do not obstruct the usefulness of the
calendar.  The filename can be a URL, and must start with http:// and must have
a valid image extension.


### Year calendar

    -yearA 

    -yearB

Two different layouts are available. One with the months on the top and the
days on the left and vice versa. Obviously there is less space for the
individual day in this mode. Still, many of the options are available here.

    -spread NUMBER

In the year calendars, the entire calendar is put on one page. Using the
spread option, you can spread the year across NUMBER pages. Valid values are
1 (default), 2, 3, 4, 6 and 12. Using 12 is essentially the same as not using the
year mode, because it will put every month on its own page.

Example:

   -yearB -spread 4

This will put three month on each page.


Event File
==================

This is a sample file event configuration file. 
Image can be a URL, which must start with http://


    <Gocal>
      <Gocaldate date="1/15"  text="Alice" />
      <Gocaldate date="2/15"  text="Bob" />
      <Gocaldate date="3/15"  text="Charles" />
      <Gocaldate date="4/15"  text="Daisy" />
      <Gocaldate date="5/15"  text="Æþelbryht" />
      <Gocaldate date="6/15"  text="Frank" />
      <Gocaldate date="6/15"  text="\nGeorge" />
      <Gocaldate date="7/15"  text="Henry" />
      <Gocaldate date="8/15"  text="Isildur" />
      <Gocaldate date="9/15"  text="Ethelbert" />
      <Gocaldate date="10/15"  text="Æþelbyrht" image="golang-gopher.png" />
      <Gocaldate date="11/15"  text="Eðilberht" />
      <Gocaldate date="Monday" text="Run Marathon" />
      <Gocaldate date="*/20" text="Pay rent" />
    </Gocal>

Please note the cool Anglo-Saxon/Scandinavian letters, thanks to UTF-8 support.

This is a sample of the configuration file for gocal. It has all the supported
features. date is in MONTH/DAY format. The text may contain a literal \n
newline.  For the month a * is permitted and it obviously means 'every month'.
You can use a leading newline symbol to make the text wrap to the next line in
case of overlap. THe optional image tag will put an image into the cell.

For the day an English Weekday name is permitted. It means: Every
matching weekday.

I was considering to allow to configure all the options from the command line
also as parameters in the XML, but I think it's not really that important.

The image can also be URL, but keep in mind, that every image will be
downloaded every time, because the files are downloaded to a temporary folder
which is deleted after gocalendar is done.

ICS iCalendar files
===================

Using
	-ics filename or URL

you can provide one or more ICS calendar objects. The events in
the calendar will be added on matching dates.

Timezones are ignored. There is still no automatic linebreaking
and no prevention of overlap with other configuration event
entries.

From the ICS file, the *SUMMARY* attribute is added as text to
the calendar.

Example:
	gocalendar -ics http://www.google.com/calendar/ical/de.german%23holiday%40group.v.calendar.google.com/public/basic.ics 
	Add German holidays to your calendar.

Example:
	gocalendar -ics my.ics
	Add your own local ICS file.

Examples
========

There is more than one way to create some example calendars.

You can use `go test` to validate that the library works 
or you can run `make` after compilation to check that the standalone
tool works. Both procedures provide ample suggestions what you can do.

The blue frames are not part of the Gocal output, but have been
added for these screenshots.
 
![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example01.png?raw=true) 

	gocalendar -o example01.pdf -p P -photos pics 1 2026

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example02.png?raw=true) 

	gocalendar -o example02.pdf -lang fr_FR -font sans 2027

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example03.png?raw=true) 

	gocalendar -o example03.pdf -wall golang-gopher.png -lang de_DE -font data/Borel-Regular.ttf 2028

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example04.png?raw=true) 

	gocalendar -o example04.pdf -lang de_DE -font mono 2 2029

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example05.png?raw=true) 

	gocalendar -o example05.pdf -lang nl_NL -plain 3 2030

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example06.png?raw=true) 

	gocalendar -o example06.pdf -font data/Borel-Regular.ttf -lang en_US 4 2031

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example07.png?raw=true) 

	gocalendar -o example07.pdf -p P -lang fr_FR -photo pics\taxi.JPG 4 2032

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example08.png?raw=true) 

	gocalendar -o example08.pdf -lang fr_FR -photo golang-gopher.png  4 2033

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example09.png?raw=true) 

	gocalendar -o example09.pdf -lang fi_FI -font serif -p L 4 2034

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example10.png?raw=true) 

	gocalendar -o example10.pdf -lang fi_FI -font serif -p L 12 2035

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example11.png?raw=true) 

	gocalendar -o example11.pdf -lang de_DE -font sans -p L 6 9 2036

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example12.png?raw=true) 

	gocalendar -o example12.pdf -p P -photo http://golang.org/doc/gopher/frontpage.png 7 2037

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example13.png?raw=true) 

	gocalendar -o example13.pdf -font sans -noother 7 2038

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example14.png?raw=true) 

	gocalendar -o example14.pdf -small 2 2039

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example15.png?raw=true) 

	gocalendar -o example15.pdf -yearA 2040

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example16.png?raw=true) 

	gocalendar -o example16.pdf -yearB 2041

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example17.png?raw=true) 

	gocalendar -o example17.pdf -yearA -fill "c" 2042

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example18.png?raw=true) 

	gocalendar -o example18.pdf -yearB -fill "sS" 2043

![Logo](https://github.com/StefanSchroeder/Gocal/blob/master/examples/example19.png?raw=true) 

	gocalendar -o example19.pdf -ics data/german.ics 12 2023

    
Roadmap
=======

This is a mature application that I use regularly and that meets all
my needs. If you are looking for a fancy option, let me know by
writing an issue. I will consider adding it.

Known bugs
==========

* When you have multiple events on the same date, they are overlapping. I
  don't intend to fix that. Use the Newline to arrange your stuff.
* Not all text will fit into the cells with some settings, because the font size is
  not adapted dynamically to the paper format. It's a feature, not a bug.
* When using the A5 paper size, the last row of a page wraps to the next page.
* Some warnings in libraries might irritate the user.
* The dates and months are not validated. Nothing prevents you from trying to 
  generate a calendar for "13 2014", which will panic.


Acknowledgments
================

I'd like to thank the developers who wrote the great libraries that **gocal** is 
relying on, especially Sonia Keys and Kurt Jung and Йордан Пулов.


