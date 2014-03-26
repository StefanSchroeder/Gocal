
Gocal
=====

Gocal is a simple clone of pcal. It's a standalone tool to create monthly calendars in PDF with a no-nonsense attitude.
It creates by default a 12-page PDF with one month per page for the current year. 
Alternatively the following arguments are supported:

  gocal YEAR => Create a 12-page calendar for YEAR

  gocal MONTH YEAR => Create a 1-page calendar for the MONTH in YEAR

  gocal BEGIN-MONTH END-MONTH YEAR => Create a sequence from BEGIN-MONTH to END-MONTH in YEAR

Gocal has a few gimmicks:

* PDF generation
* Optional week number on every Monday
* Optional day of year number on every day
* Optional moon phases
* Events in XML configuration file
* Several languages 
* Wallpaper option
* Photo calendar option (from single image or directory)
* Page orientation and paper size option
* Font selection (limited to Times, Helvetica/Arial, Courier)

Pcal has quite a few options, that gocal will never support, like calculating Easter. 

Why not using pcal or gcal? There are several reasons. One, I couldn't build it on 
Windows, next, I don't care about Postscript any longer, third, a calendar seems to 
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


Options
=======

  -f="Gocal": Footer note

Change the string at the bottom of the page.

  -font="Cabalett": Main font (Times/Arial/Courier)

Set the font. The fonts times, arial/helvetica and courier are built-in in PDF.
Other fonts are currently not supported, but these fonts should be sufficient for
all practical purposes. As a proof of concept I included the freeware Cabalett-font by
Harold Lohner (http://haroldsfonts.com/) which is really beautiful.

  -lang="": Language

Gocal reads the LANG environment variable. If it matches one of 

		en_US en_GB nl_BE nl_NL fi_FI fr_FR fr_CA de_DE

the library goodsign/monday is used to translate the weekday names and month names.
Although this library supports a few other languages, I found that many of the 
languages do not work. The language from the environment can be overridden with this
parameter. If your LANG is not recognized, we default to en_US.

  -nodoy=false: Hide day of year (false)

  -noevents=false: Hide events from config file (false)

Gocal adds events from an XML file in the current directory to dates. This option
disables this feature.

  -nomoon=false: Hide moon phases (false)
  -noweek=false: Hide week number (false)

The week number according to ISO-8601, added on every Monday.

  -o="output.pdf": Output filename
  -p="L": Orientation (L)andscape/(P)ortrait

Typically you want landscape for calendars without image and portrait for calendars with image.

  -paper="A4": Paper format (A3 A4 A5 Letter Legal)
  -photo="": Show photo (single image PNG JPG GIF)

This option will add this image to every month.

  -photos="": Show photos (directory PNG JPG GIF)

e.g. gocal -photos images/

This option will add the twelve first files as images to the twelve month.
If less than twelve files are found, the sequence will re-start after the last image.

  -wall="": Show wallpaper PNG JPG GIF

e.g. gocal -wall gopher.png 

This option will add this image to every month as a background. You should only 
use images with a bright tone so that you do not obstruct the usefulness of the calendar.


Examples
========

 
![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screen0.png?raw=true)

German locale with font Cabaletta and gopher wallpaper.

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screen1.png?raw=true)

Finnish locale with all gimmicks hidden.

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screen2.png?raw=true)

English locale with default settings and Cabaletta font.

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screen3.png?raw=true)

English locale in photo-mode.

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screen4.png?raw=true)

Finnish locale with times font (landscape).

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screen5.png?raw=true)

English locale with times font (portrait). 


Roadmap
=======

It would be really cool to allow gocal to be a drop-in replacement for pcal, OTOH why bother?


Known bugs:
===========

* When you have multiple events on the same date, they are overlapping.
* Not all text will fit into the cells with some settings.

