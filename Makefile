build:
	env GOOS=linux GOARCH=arm go build . && mv autok3s bin/autok3s

RPI_IPS=192.168.80.20 192.168.80.21
scp: build
	for ip in $(RPI_IPS); do \
		sudo scp bin/autok3s pi@$$ip:/home/pi/autok3s ;\
	done