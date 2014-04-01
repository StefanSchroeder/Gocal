#!/bin/bash	
C="convert -alpha Opaque "
E=./Gocal.exe
D=examples
$E -o example01.pdf -p P -photos pics 1 2014
$C  example01.pdf  $D/example01.png 
$E -o example02.pdf -lang fr_FR -font sans 2015
$C  example02.pdf  $D/example02.png 
$E -o example03.pdf -wall golang-gopher.png -lang de_DE -font c:/windows/Fonts/cabalett.ttf 
$C  example03.pdf  $D/example03.png 
$E -o example04.pdf -lang de_DE -font mono 2 2014
$C  example04.pdf  $D/example04.png 
$E -o example05.pdf -lang nl_NL -plain 3 2014
$C  example05.pdf  $D/example05.png 
$E -o example06.pdf -font c:/windows/Fonts/cabalett.ttf -lang en_US 4 2014
$C  example06.pdf  $D/example06.png 
$E -o example07.pdf -p P -lang fr_FR -photo pics\taxi.JPG  4 2014
$C  example07.pdf  $D/example07.png 
$E -o example08.pdf -lang fr_FR -photo golang-gopher.png  4 2014
$C  example08.pdf  $D/example08.png 
$E -o example09.pdf -lang fi_FI -font serif -p L  4 2014
$C  example09.pdf  $D/example09.png 
$E -o example10.pdf -lang fi_FI -font serif -p L 12 2013
$C  example10.pdf  $D/example10.png 
$E -o example11.pdf -lang de_DE -font sans -p L 6 9 2014
$C  example11.pdf  $D/example11.png 
$E -o example12.pdf -p P -photo http://golang.org/doc/gopher/frontpage.png 7 2014
$C  example12.pdf  $D/example12.png 
$E -o example13.pdf -font sans -noother 7 2014
$C  example13.pdf  $D/example13.png 
$E -o example14.pdf -small 2 2014
$C  example14.pdf  $D/example14.png 

