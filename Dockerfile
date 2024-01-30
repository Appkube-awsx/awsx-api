# Create image based on ubuntu 20.04
FROM synectiks/awsx-api-base:latest
SHELL ["/bin/bash", "-c"] 

EXPOSE 7000
RUN go build
ENTRYPOINT [ "go run ./main.go start" ]