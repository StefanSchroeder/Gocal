#!/bin/bash	
C="convert -alpha Opaque "
./Gocal.exe -o example01.pdf -p P -photos pics 1 2014
$C  example01.pdf  example01.png 
./Gocal.exe -o example02.pdf -lang fr_FR -font sans 2015
$C  example02.pdf  example02.png 
./Gocal.exe -o example03.pdf -wall golang-gopher.png -lang de_DE -font c:/windows/Fonts/cabalett.ttf 
$C  example03.pdf  example03.png 
./Gocal.exe -o example04.pdf -lang de_DE -font mono 2 2014
$C  example04.pdf  example04.png 
./Gocal.exe -o example05.pdf -lang nl_NL -plain 3 2014
$C  example05.pdf  example05.png 
./Gocal.exe -o example06.pdf -font c:/windows/Fonts/cabalett.ttf -lang en_US 4 2014
$C  example06.pdf  example06.png 
./Gocal.exe -o example07.pdf -p P -lang fr_FR -photo pics\taxi.JPG  4 2014
$C  example07.pdf  example07.png 
./Gocal.exe -o example08.pdf -lang fr_FR -photo golang-gopher.png  4 2014
$C  example08.pdf  example08.png 
./Gocal.exe -o example09.pdf -lang fi_FI -font serif -p L  4 2014
$C  example09.pdf  example09.png 
./Gocal.exe -o example10.pdf -lang fi_FI -font serif -p L 12 2013
$C  example10.pdf  example10.png 
./Gocal.exe -o example11.pdf -lang de_DE -font sans -p L 6 9 2014
$C  example11.pdf  example11.png 
./Gocal.exe -o example12.pdf -p P -photo http://golang.org/doc/gopher/frontpage.png 7 2014
$C  example12.pdf  example12.png 
./Gocal.exe -o example13.pdf -font sans -noother 7 2014
$C  example13.pdf  example13.png 
./Gocal.exe -o example14.pdf -small 2 2014
$C  example14.pdf  example14.png 

