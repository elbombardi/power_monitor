install: 
	make clean
	make build
	make install_service
	make clean

clean:
	rm -fr build

build: 
	mkdir build
	go build -o ./build/

install_service:
	sudo systemctl stop power-monitor-service.service
	sudo cp ./build/power_monitor /usr/local/bin/power_monitor
	sudo cp power-monitor-service.service /etc/systemd/system/power-monitor-service.service
	sudo systemctl daemon-reload
	sudo systemctl enable power-monitor-service.service
	sudo systemctl start power-monitor-service.service
	sudo systemctl status power-monitor-service.service

read_logs:
	journalctl -u power-monitor-service.service -f

.PHONY: clean build install_service read_logs