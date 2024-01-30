# Create image based on ubuntu 20.04
FROM synectiks/awsx-api-base:latest
SHELL ["/bin/bash", "-c"] 
WORKDIR /app
ARG ARTIFACT_NAME=awsx-api
ARG CONF_FILE=conf/config.yaml
COPY ${ARTIFACT_NAME} /app/
RUN MKDIR /app/conf
COPY ${CONF_FILE} /app/conf
EXPOSE 7000
ENTRYPOINT [ "./awsx-api start" ]