FROM debian:latest
# Para evitar que salgan mensajes de la configuración de Wireshark
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y firefox-esr wireshark iputils-ping iproute2 sudo nano
RUN groupadd -g 1000 usuario && useradd -m -u 1000 -g usuario usuario && echo usuario:usuario | chpasswd && echo "usuario ALL=(ALL) ALL" >> /etc/sudoers && echo "Defaults lecture = never" >> /etc/sudoers
USER usuario
WORKDIR /home/usuario
CMD bash