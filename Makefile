build:
	GOOS=linux GOARCH=arm go build -o hemepi main.go

deploy:
	scp ./hemepi pi@raspberrypi.local:~/Development
