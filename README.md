
Gocal
=====

Gocal is a simple clone of pcal. It's a standalone tool to create monthly calendars in PDF with a no-nonsense attitude.
By default it creates a 12-page PDF with one month per page for the current year. 

Alternatively the following arguments are supported:

  gocal YEAR => Create a 12-page calendar for YEAR, e.g. 2014

  gocal MONTH YEAR => Create a 1-page calendar for the MONTH in YEAR, e.g. 5 2014

  gocal BEGIN-MONTH END-MONTH YEAR => Create a sequence from BEGIN-MONTH to END-MONTH in YEAR, e.g. 5 7 2014


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

Why not using pcal? There are several reasons. One, I couldn't build it on 
Windows, next, I don't care about Postscript, third, a calendar seems to 
be exactly a kind of project that I am able to handle from a complexity perspective 
as a single developer.

The main idea of gocal is simplicity. While it is absolutely possible to create
an application where every single stroke is configurable, I believe that most of
you are too busy to play around with lots of options and want a pleasant calendar
out-of-the-box. That's what gocal provides.

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screenshot.png?raw=true)

The power of gocal is based on the cool libraries that it includes. This implies that
several of the options are actually options of the libraries, e.g. the paper format (gofpdf) or
the supported languages (goodsign/monday).

It is suggested to hide some of the optional fields or the cells will look crowded.


Build instructions
==================

Run 

	go get github.com/StefanSchroeder/gocal

Gocal has a quite a few dependencies that go should resolve automatically.

Build with

	go build 


License
=======

The BSD type license is in the LICENSE file.


Options
=======

		-h  Help: Summarizes the options.

		-f="Gocal": Footer note

Change the string at the bottom of the page.

		-font="": font

		-font serif    (Times lookalike)
		-font mono     (Courier lookalike)
		-font sans     (Arial lookalike)
		-font path/to/font.ttf    (your favorite font)

Set the font. Only truetype fonts are supported. Look into c:\Windows\Fonts on Windows
and /usr/share/fonts to see what fonts you have. Gocal comes with three fonts built-in:
The Gnu FreeMonoBold, FreeSerifBold and FreeSansBold. They look pretty similar to 
(in that order) Courier, Times and Arial and should meet all your standard font needs.
These fonts are licensed under the Gnu FreeFont License which accompanies this README.txt.
Read more about them at https://www.gnu.org/software/freefont. To use fonts gocal, 
auxiliary files are created in a temporary directory.

In addition you can provide your own TTF on the commandline if you prefer something fancy.

		-lang="": Language

Gocal reads the LANG environment variable. If it matches one of 

		en_US en_GB nl_BE nl_NL fi_FI fr_FR fr_CA de_DE

the library goodsign/monday is used to translate the weekday names and month names.
Although this library supports a few other languages, I found that many of the 
languages do not work with the fonts I tried. The language from the environment can be 
overridden with this parameter. If your LANG is not recognized, we default to en_US.

		-nodoy: Hide day of year

		-noevents: Hide events from config file

Gocal adds events from an XML file in the current directory to dates. This option
disables this feature.

		-nomoon: Hide moon phases

		-noweek: Hide week number

The week number according to ISO-8601 is added on every Monday by default.

		-o="output.pdf": Output filename

		-p="L": Orientation (L)andscape/(P)ortrait

Typically you want landscape for calendars without image and portrait for calendars with image.

		-paper="A4": Paper format (A3 A4 A5 Letter Legal)

		-photo=filename: Show single photo (single image in PNG JPG GIF)

This option will add this image to every month.

		-photos=directory: Show multiple photos (directory with PNG JPG GIF)

e.g. gocal -photos images/

This option will add the twelve first files as images to the twelve month.
If less than twelve files are found, the sequence will re-start after the last image.
This will not work if there are non-image files in the directory (among the first twelve).

		-wall=filename: Show wallpaper PNG JPG GIF

e.g. gocal -wall gopher.png

This option will add this image to every month as a background. You should only 
use images with a bright tone so that you do not obstruct the usefulness of the calendar.

		-plain This will hide everything that can be hidden.


Event File
==================

Choosing a format for the configuration file was tough.
You might think that XML is overkill and JSON or CSV or one
of the many configuration file libraries would have been 
more adequate. Perhaps you are right. The advantage of the
XML package is, that I knew how to use it and also, the 
extensibility should I ever choose to allow more complex
configuration. The defaultname is gocal.xml.

	<Gocal>
    <Gocaldate date="12/24" text="Heilig Abend" />
    <Gocaldate date="12/25" text="Weihnachten" />
    <Gocaldate date="1/6"  text="Allerheiligen" />
    <Gocaldate date="2/14" text="Enno" />
    <Gocaldate date="5/23" text="Stefan\nGründung der BRD" />
    <Gocaldate date="5/20" text="\nBirgit" />
    <Gocaldate date="9/18" text="Tomke" />
    <Gocaldate date="*/20" text="Miete" />
	</Gocal>

This is a sample of the really primitive configuration file
for gocal. It has all the supported features. date is 
in MONTH/DAY format. The text may contain a literal \n newline.
For the month a * is permitted and it obviously means 'every month'.
You can use a leading newline symbol to make the text wrap to the next
line in case of overlap.

I was considering to allow to configure all the options from the command line
also as parameters in the XML, but I think it's not really that important.


Examples
========

Run the samples.bat batch file to create some example calendars.
 
![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example01.png?raw=true)

    Gocal -o example01.pdf -p P -photos pics 1 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example02.png?raw=true)

    Gocal -o example02.pdf -lang fr_FR -font sans 2015

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example03.png?raw=true)

    Gocal -o example03.pdf -wall golang-gopher.png -lang de_DE -font c:\windows\Fonts\cabalett.ttf 

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example04.png?raw=true)

    Gocal -o example04.pdf -lang de_DE -font mono 2 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example05.png?raw=true)

    Gocal -o example05.pdf -lang nl_NL -plain 3 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example06.png?raw=true)

    Gocal -o example06.pdf -font c:\windows\Fonts\cabalett.ttf -lang en_US 4 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example07.png?raw=true)

    Gocal -o example07.pdf -p P -lang fr_FR -photo pics\taxi.JPG  4 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example08.png?raw=true)

    Gocal -o example08.pdf -lang fr_FR -photo golang-gopher.png  4 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example09.png?raw=true)

    Gocal -o example09.pdf -lang fi_FI -font serif -p L  4 2014

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example10.png?raw=true)

    Gocal -o example10.pdf -lang fi_FI -font serif -p L 12 2013

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/example11.png?raw=true)

    Gocal -o example11.pdf -lang de_DE -font sans -p L 6 9 2014

    
Roadmap
=======

* It would be really cool to allow gocal to be a drop-in replacement for pcal,
but the configuration file for pcal is really complex.
* Also, I'd love to setup a website, where you generate PDFs on the web.
* Your feature here!
* I have no plans to put multiple month on a single page (yearly calendar). If
you need that, create a 12-page calendar and use the plethora of pdf-tools
to rearrange the PDF that gocal produced. E.g. pdf2ps, psbook, psnup, pdftk, etc. 

Known bugs
==========

* The event file must be encoded in UTF-8.
* When you have multiple events on the same date, they are colliding. I
  don't intend to fix that.
* Not all text will fit into the cells with some settings, because the font size is
  not adapted dynamically to the paper format. It's a feature, not a bug.
* When using the A5 paper size, the last row of a page wraps to the next page.
* Some warnings in libraries might irritate the user.


Acknowledgments
================

I'd like to thank the developers who wrote the great libraries that **gocal** is 
relying on, esp. Sonia Keys and Kurt Jung.


