test:
	go test -v ./...
build:
	GOOS=linux GOARCH=arm go build -o hemepi ./...

deploy:
	rsync -r ./hemepi pi@raspberrypi.local:~/hemepi/
	rsync -r ./assets pi@raspberrypi.local:~/hemepi/
