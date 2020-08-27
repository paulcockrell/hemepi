build:
	GOOS=linux GOARCH=arm go build -o hemepi ./...

deploy:
	scp ./hemepi pi@raspberrypi.local:~/Development
