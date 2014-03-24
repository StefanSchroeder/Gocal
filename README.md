Warning: Gocal is pre-alpha and released here in good faith. Don't use it to control nuclear power plants.

Gocal
=====

Gocal is a simple clone of pcal. It's a tool to create monthly calendars in PDF with a no-nonsense attitude.
Right now it creates a 12-page PDF with one month per page. Pcal has quite a few options, that gocal will
never support, like calculating Easter. 

Features:

* PDF generation
* Week number on every Monday
* Events in XML configuration file
* No other languages than English currently
* No other 'first day of the week' than Monday

Why not using PCal? There are several reasons. One, I couldn't build it on Windows, next, I don't care about Postscript any longer, third, a calendar seems to be exactly a kind of project that I am able to handle from a complexity perspective as a single developer.


The main idea of gocal is that you print one or more month and use it to add 
handwritten notes to it.

![Logo](http://github.com/StefanSchroeder/Gocal/blob/master/screenshot.png?raw=true)

Roadmap
=======

It would be really cool to allow gocal to be a drop-in replacement for pcal, OTOH why bother?

Also, I think that gocal would be a good sample project to study i18n with Go.


Known bugs:
===========

* I dropped sunrise/sunset.
