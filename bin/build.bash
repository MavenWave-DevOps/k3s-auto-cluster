env GOOS=linux GOARCH=arm go build .

echo "raspberry" | sudo scp autok3s pi@192.168.80.20:/home/pi/k3s-auto-cluster/autok3s
echo "raspberry" | sudo scp autok3s pi@192.168.80.21:/home/pi/autok3s