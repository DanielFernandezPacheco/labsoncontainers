FROM debian:latest
RUN apt-get update && apt-get install -y firefox-esr wireshark iputils-ping iproute2
RUN groupadd -g 1000 usuario && useradd -m -u 1000 -g usuario usuario
USER usuario
WORKDIR /home/usuario
CMD bash