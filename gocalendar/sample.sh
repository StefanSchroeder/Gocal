#!/bin/bash	
C="pdftoppm "
C2="convert -antialias -bordercolor SkyBlue -border 4x4 "
E="go run gocalendar.go"  
D=../examples

$E -o example01.pdf -p P -photos pics 1 2014
$C  example01.pdf > example01.ppm 
$C2  example01.ppm  $D/example01.png 

$E -o example02.pdf -lang fr_FR -font sans 1 2015
$C  example02.pdf >  example02.ppm 
$C2   example02.ppm  $D/example02.png 

$E -o example03.pdf -wall golang-gopher.png -lang de_DE -font c:/windows/Fonts/cabalett.ttf  1 2015
$C  example03.pdf >  example03.ppm 
$C2   example03.ppm  $D/example03.png 

$E -o example04.pdf -fontscale 0.6 -lang de_DE -font mono 2 2014
$C  example04.pdf >  example04.ppm 
$C2   example04.ppm  $D/example04.png 

$E -o example05.pdf -lang nl_NL -plain 3 2014
$C  example05.pdf >  example05.ppm 
$C2   example05.ppm  $D/example05.png 

$E -o example06.pdf -font c:/windows/Fonts/cabalett.ttf -lang en_US 4 2006
$C  example06.pdf >  example06.ppm 
$C2   example06.ppm  $D/example06.png 

$E -o example07.pdf -p P -lang fr_FR -photo pics/taxi.JPG  4 2007
$C  example07.pdf >  example07.ppm 
$C2   example07.ppm  $D/example07.png 

$E -o example08.pdf -lang fr_FR -photo golang-gopher.png  4 2008
$C  example08.pdf >  example08.ppm 
$C2   example08.ppm  $D/example08.png 

$E -o example09.pdf -lang fi_FI -font serif -p L  4 2009
$C  example09.pdf >  example09.ppm 
$C2   example09.ppm  $D/example09.png 

$E -o example10.pdf -lang fi_FI -font mono -p L 12 2010
$C  example10.pdf >  example10.ppm 
$C2   example10.ppm  $D/example10.png 

$E -o example11.pdf -lang de_DE -font sans -p L 6 2011
$C  example11.pdf >  example11.ppm 
$C2   example11.ppm  $D/example11.png 

$E -o example12.pdf -p P -photo http://golang.org/doc/gopher/frontpage.png 7 2012
$C  example12.pdf >  example12.ppm 
$C2   example12.ppm  $D/example12.png 

$E -o example13.pdf -fontscale 0.9 -font sans -noother 7 2013
$C  example13.pdf >  example13.ppm 
$C2   example13.ppm  $D/example13.png 

$E -o example14.pdf -small 2 2014
$C  example14.pdf >  example14.ppm 
$C2   example14.ppm  $D/example14.png 

$E -o example15.pdf -yearA 2015
$C  example15.pdf >  example15.ppm 
$C2   example15.ppm  $D/example15.png 

$E -o example16.pdf -yearB -p L 2016
$C  example16.pdf >  example16.ppm 
$C2   example16.ppm  $D/example16.png 

$E -o example17.pdf -yearA -fill c 2017
$C  example17.pdf >  example17.ppm 
$C2   example17.ppm  $D/example17.png 

$E -o example18.pdf -yearB -fill sS 2018
$C  example18.pdf >  example18.ppm 
$C2   example18.ppm  $D/example18.png 

rm example*.ppm example*.pdf

