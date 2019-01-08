all:
	go run gocalendar/gocalendar.go -spread 1 -yearA -o o1.pdf -lang de_DE 2019
	go run gocalendar/gocalendar.go -spread 2 -yearA -o o2.pdf -lang de_DE 2019
	go run gocalendar/gocalendar.go -spread 3 -yearA -o o3.pdf -lang de_DE 2019 
	go run gocalendar/gocalendar.go -spread 4 -yearA -o o4.pdf -lang de_DE 2019 
	go run gocalendar/gocalendar.go -spread 6 -yearA -o o6.pdf -lang de_DE 2019
	go run gocalendar/gocalendar.go -spread 12 -yearA -o o12.pdf -lang de_DE 2019

b:
	go run gocalendar/gocalendar.go -spread 1 -yearB -o o1.pdf -lang de_DE 2019
	go run gocalendar/gocalendar.go -spread 2 -yearB -o o2.pdf -lang de_DE 2019
	go run gocalendar/gocalendar.go -spread 3 -yearB -o o3.pdf -lang de_DE 2019 
	go run gocalendar/gocalendar.go -spread 4 -yearB -o o4.pdf -lang de_DE 2019 
	go run gocalendar/gocalendar.go -spread 6 -yearB -o o6.pdf -lang de_DE 2019
	go run gocalendar/gocalendar.go -spread 12 -yearB -o o12.pdf -lang de_DE 2019
	
