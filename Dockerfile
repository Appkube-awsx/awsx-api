# Create image based on ubuntu 20.04
FROM synectiks/awsx-api-base:latest
SHELL ["/bin/bash", "-c"] 
WORKDIR /awsx-api
EXPOSE 7000
ENTRYPOINT [ "go run ./main.go start" ]