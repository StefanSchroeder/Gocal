#!/bin/bash	
C="convert -alpha Opaque "
C2="convert -antialias -bordercolor SkyBlue -border 10x10 "
E="go run gocalendar.go"  
D=../examples


$E -o example01.pdf -p P -photos pics 1 2014
$C  example01.pdf example01.png 
$C2  example01.png  $D/example01.png 

$E -o example02.pdf -lang fr_FR -font sans 1 2015
$C  example02.pdf  example02.png 
$C2   example02.png  $D/example02.png 

$E -o example03.pdf -wall golang-gopher.png -lang de_DE -font c:/windows/Fonts/cabalett.ttf  1 2015
$C  example03.pdf  example03.png 
$C2   example03.png  $D/example03.png 

$E -o example04.pdf -fontscale 0.6 -lang de_DE -font mono 2 2014
$C  example04.pdf  example04.png 
$C2   example04.png  $D/example04.png 

$E -o example05.pdf -lang nl_NL -plain 3 2014
$C  example05.pdf  example05.png 
$C2   example05.png  $D/example05.png 

$E -o example06.pdf -font c:/windows/Fonts/cabalett.ttf -lang en_US 4 2006
$C  example06.pdf  example06.png 
$C2   example06.png  $D/example06.png 

$E -o example07.pdf -p P -lang fr_FR -photo pics/taxi.JPG  4 2007
$C  example07.pdf  example07.png 
$C2   example07.png  $D/example07.png 

$E -o example08.pdf -lang fr_FR -photo golang-gopher.png  4 2008
$C  example08.pdf  example08.png 
$C2   example08.png  $D/example08.png 

$E -o example09.pdf -lang fi_FI -font serif -p L  4 2009
$C  example09.pdf  example09.png 
$C2   example09.png  $D/example09.png 

$E -o example10.pdf -lang fi_FI -font mono -p L 12 2010
$C  example10.pdf  example10.png 
$C2   example10.png  $D/example10.png 

$E -o example11.pdf -lang de_DE -font sans -p L 6 2011
$C  example11.pdf  example11.png 
$C2   example11.png  $D/example11.png 

#$E -o example12.pdf -p P -photo http://golang.org/doc/gopher/frontpage.png 7 2012
#$C  example12.pdf  example12.png 
#$C2   example12.png  $D/example12.png 

$E -o example13.pdf -fontscale 0.9 -font sans -noother 7 2013
$C  example13.pdf  example13.png 
$C2   example13.png  $D/example13.png 

$E -o example14.pdf -small 2 2014
$C  example14.pdf  example14.png 
$C2   example14.png  $D/example14.png 

$E -o example15.pdf -yearA 2015
$C  example15.pdf  example15.png 
$C2   example15.png  $D/example15.png 

$E -o example16.pdf -yearB -p L 2016
$C  example16.pdf  example16.png 
$C2   example16.png  $D/example16.png 

$E -o example17.pdf -yearA -fill c 2017
$C  example17.pdf  example17.png 
$C2   example17.png  $D/example17.png 

$E -o example18.pdf -yearB -fill sS 2018
$C  example18.pdf  example18.png 
$C2   example18.png  $D/example18.png 

