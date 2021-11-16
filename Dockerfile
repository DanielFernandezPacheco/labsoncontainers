FROM debian:latest
RUN apt-get update && apt-get install -y firefox-esr wireshark iputils-ping iproute2
CMD bash