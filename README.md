- [awsx-api](#awsx-api)
- [project architecture](#project-architecture)
- [api-endpoint](#api-endpoint) 
- [start server](#start-server)
- [Details of All Sub Command](#details-of-all-sub-command)





# awsx-api
awsx-api is a golang based REST api server that exposes GET, POST, DELETE and PUT endpoints that will subsequently allow us to perform the full range of operations on AWS entities.\


# project architecture
![project structure](project_structure.png "project structure")

    NOTE: To perform operation on AWS entities, awsx-api uses aws-sdk and go cli packages written specifically to deal with AWS entities

    1. main.go
        To get started we have to create a web server which can handle HTTP requests. 
        To do this we have a file called main.go.
        The main function in main.go kiks off the server.
            go run .\main.go start
    
    2. Configuration
        * config.yaml: config.yaml is the main configuration file which contain all the server configuration like IP address, port etc.

                server:
                    address: localhost
                    port: 7000
                    static_content_root_directory: /home/userTests/awsx-api-static-files
                    cors_allow_all: true
                    white_list_urls: http://localhost:3002

        * config.go: All the code of reading the configuration from config.yaml file and creating the global config reference is written in config.go  

    3. server
        * server.go: server.go contains the code to create, start and stop the web server

    4. routing
        * For http routing, awsx-api uses the gorilla/mux router 
        * router.go: It contains the logic to create new router and all the http end-points are defined in router.go
    
    5. handlers
        * All the business logic to interact with different AWS entities and perform operations on them will be written in handlers package
            1. appconfig-handler.go: It gets the resource config summary of any AWS account

    6. models
        * models.go: models.go will have all the models needed in api operations

    7. log
        * log.go: A custom log.go created for awsx-api.

# api-endpoint 
    
https://github.com/Appkube-awsx/awsx-api/blob/main/specs/allDetailsFile/allElementDetails.md

- build/run/debug/test in postman

# start server
    go run .\main.go start

# Details of All Sub Command

All the supported subcommands and there source code locations are mentiioned in 

for getElementDetails

https://github.com/Appkube-awsx/awsx-api/tree/main/specs/allElementPanel

for costDetails

https://github.com/Appkube-awsx/awsx-cost/blob/main/awsx-costData/README.md
        
    

| S.No | Sub-command           | Description                                           | Output Format                                  | Functionalities                                                                                                                                                                            | Specs Links |
|------|-----------------------|-------------------------------------------------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| 1    | getLandingZoneDetails | Collect Information about any specific landing zone  | Percentage(%)                               | 1. Get Elements Metadata, 2. Get List of every elements with their config infos, 3. Get List of products and Environments for the landing Zone, 4. Get the cost of the landing zone   |    |
| 2    | getLandingZoneCompliance | Collect Information about any specific landing zone compliances and security | Bytes                                 | 1. Get overall account Compliance, 2. Get the elementWise Compliance, 3. Run Compliance on products and environments                                                                       |  |
| 3    | getElementDetails | Collect Information about specific cloud elements -Run Queries | Percentage(%), Bytes                               | 1. EC2, 2. EKS, 3. ECS, 4. LAMBDA, 5. API Gw, 6. Load Balancer  |    [https://github.com/Appkube-awsx/awsx-api/tree/main/specs/allElementPanel](https://github.com/Appkube-awsx/awsx-api/tree/main/specs/allElementPanel)   |
| 4    | getCostDetails        | Collect Information about account and elements specific costs | Bytes                            | 1. Total Account, 2. Product and Envwise, 3. Element Wise, 4. Spikes and Trends, 5. App/Data/Nw Service wise Costs                                                                      | [https://github.com/Appkube-awsx/awsx-cost](https://github.com/Appkube-awsx/awsx-cost)|
