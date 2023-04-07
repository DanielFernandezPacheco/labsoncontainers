all:
	CGO_ENABLED=0 go build -o labsoncontainers ./cmd/LabsOnContainers/*

install:
	cp ./labsoncontainers /usr/local/bin
	chown root:root /usr/local/bin/labsoncontainers
	chmod 755 /usr/local/bin/labsoncontainers
	chmod u+s /usr/local/bin/labsoncontainers
	mkdir -p /home/usuario/.labsoncontainers
	mkdir -p /home/usuario/.labsoncontainers/recent_configs
	cp -r ./public/. /home/usuario/.labsoncontainers
	chmod -R 755 /home/usuario/.labsoncontainers


uninstall:
	rm /usr/local/bin/labsoncontainers

clean:
	rm ./labsoncontainers
