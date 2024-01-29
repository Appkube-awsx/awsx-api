# Create image based on ubuntu 20.04
FROM ubuntu:jammy-20240111
SHELL ["/bin/bash", "-c"] 
RUN apt update && apt upgrade -y \
 && apt install sudo curl wget zip unzip -y \
 && wget https://go.dev/dl/go1.19.13.linux-amd64.tar.gz \
 && tar -xvf go1.19.13.linux-amd64.tar.gz -C /usr/local \
 && echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc

