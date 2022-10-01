all:
	CGO_ENABLED=0 go build -o labsoncontainers ./cmd/LabsOnContainers/*

install:
	cp ./labsoncontainers /usr/local/bin
	chown root:root /usr/local/bin/labsoncontainers
	chmod 755 /usr/local/bin/labsoncontainers
	chmod u+s /usr/local/bin/labsoncontainers

uninstall:
	rm /usr/local/bin/labsoncontainers

clean:
	rm ./labsoncontainers