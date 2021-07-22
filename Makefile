all:
	go run gocalendar/gocalendar.go -spread 1 -yearA -o test-output/test-example_ao1.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 2 -yearA -o test-output/test-example_ao2.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 3 -yearA -o test-output/test-example_ao3.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 4 -yearA -o test-output/test-example_ao4.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 6 -yearA -o test-output/test-example_ao6.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 12 -yearA -o test-output/test-example_ao12.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 1 -yearB -o test-output/test-example_bo1.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 2 -yearB -o test-output/test-example_bo2.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 3 -yearB -o test-output/test-example_bo3.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 4 -yearB -o test-output/test-example_bo4.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 6 -yearB -o test-output/test-example_bo6.pdf -lang de_DE 2021
	go run gocalendar/gocalendar.go -spread 12 -yearB -o test-output/test-example_bo12.pdf -lang de_DE 2021
	
