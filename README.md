- [awsx-api](#awsx-api)
- [project architecture](#project-architecture)
- [api-endpoint](#api-endpoint)
- [start server](#start-server)
- [Details of All Sub Command](#details-of-all-sub-command)
- [api-to-cli](#api-to-cli)

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
    
https://github.com/Appkube-awsx/awsx-api/blob/main/specs/allgetElementDetailsList/allElementDetails.md

- build/run/debug/test in postman



# start server
    go run .\main.go start

# Details of All Sub Command

All the supported subcommands and there source code locations are mentiioned in 

for getElementDetails

https://github.com/Appkube-awsx/awsx-api/tree/main/specs/getElementDetails

for costDetails

https://github.com/Appkube-awsx/awsx-cost/blob/main/awsx-costData/README.md
        
    

| S.No | Sub-command           | Description                                           | Output Format                                  | Functionalities                                                                                                                                                                            | Specs Links |
|------|-----------------------|-------------------------------------------------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| 1    | getLandingZoneDetails | Collect Information about any specific landing zone  | Percentage(%)                               | 1. Get Elements Metadata, 2. Get List of every elements with their config infos, 3. Get List of products and Environments for the landing Zone, 4. Get the cost of the landing zone   |    |
| 2    | getLandingZoneCompliance | Collect Information about any specific landing zone compliances and security | Bytes                                 | 1. Get overall account Compliance, 2. Get the elementWise Compliance, 3. Run Compliance on products and environments                                                                       |  |
| 3    | getElementDetails | Collect Information about specific cloud elements -Run Queries | Percentage(%), Bytes                               | 1. EC2, 2. EKS, 3. ECS, 4. LAMBDA, 5. API Gw, 6. Load Balancer  |    [https://github.com/Appkube-awsx/awsx-api/blob/main/specs/allgetElementDetailsList/allElementDetails.md](https://github.com/Appkube-awsx/awsx-api/blob/main/specs/allgetElementDetailsList/allElementDetails.md)   |
| 4    | getCostDetails        | Collect Information about account and elements specific costs | Bytes                            | 1. Total Account, 2. Product and Envwise, 3. Element Wise, 4. Spikes and Trends, 5. App/Data/Nw Service wise Costs                                                                      | [https://github.com/Appkube-awsx/awsx-cost](https://github.com/Appkube-awsx/awsx-cost)|




| S.No | Sub-command           | Description                                           | Output Format                                  | Functionalities                                                                                                                                                                            | Specs Links |
|------|-----------------------|-------------------------------------------------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| 1    | ec2 | go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="1234" --query="cpu_utilization_panel" --elementType="AWS/EC2" 
--responseType=json --startTime="" --endTime=""`  | Percentage(%)                               | 1. Get Elements Metadata, 
| 3    | getElementDetails | Collect Information about specific cloud elements -Run Queries | Percentage(%), Bytes                               | 1. EC2, 2. EKS, 3. ECS, 4. LAMBDA, 5. API Gw, 6. Load Balancer  |    [https://github.com/Appkube-awsx/awsx-api/blob/main/specs/allgetElementDetailsList/allElementDetails.md](https://github.com/Appkube-awsx/awsx-api/blob/main/specs/allgetElementDetailsList/allElementDetails.md)   |
| 4    | getCostDetails        | Collect Information about account and elements specific costs | Bytes                            | 1. Total Account, 2. Product and Envwise, 3. Element Wise, 4. Spikes and Trends, 5. App/Data/Nw Service wise Costs                                               



 # api-to-cli
| S.No | API Endpoint                                                            | Sample Command                                                                                                                                        |Status                                                                                                                                        |
|------|-------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1    | `/awsx-api/getSupportedQueries?elementType=EC2`                         | `go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="1234" --query="cpu_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""` |
| 2    | `/awsx-api/getSupportedQueries?elementType=EKS`                         | `go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="1234" --query="cpu_utilization_panel" --elementType="AWS/EKS" --responseType=json --startTime="" --endTime=""` |
| 3    | `/awsx-api/getSupportedQueries?elementType=ECS`                         | `go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="1234" --query="cpu_utilization_panel" --elementType="AWS/ECS" --responseType=json --startTime="" --endTime=""` |
| 4    | `/awsx-api/getSupportedQueries?elementType=LAMBDA`                      | `go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="1234" --query="cpu_utilization_panel" --elementType="AWS/LAMBDA" --responseType=json --startTime="" --endTime=""` |
| 5    | `/awsx-api/getQueryOutput?elementType=EC2,elementId="9321",query=cpu_utilization_panel,timeRange={},responseType=json` | `go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="9321" --query="cpu_utilization_panel" --elementType="AWS/EC2" --responseType=json --startTime="" --endTime=""` | cli - done |
| 6    | `/awsx-api/getQueryOutput?elementType=EC2,elementId="9321",query=cpu_utilization_panel,timeRange={},responseType=frame` | `go run awsx-getelementdetails.go --vaultUrl=vault.synectiks.net --elementId="9321" --query="cpu_utilization_panel" --elementType="AWS/EC2" --responseType=frame --startTime="" --endTime=""` | cli - done |
| 7    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getCloudElementsMetaData`    | `go run awsx-getLandingZoneDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getCloudElementsMetaData" --responseType=json`    |
| 8    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getCloudElementsConfigData`  | `go run awsx-getLandingZoneDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getCloudElementsConfigData" --responseType=json`  |
| 9    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getProuctsData`             | `go run awsx-getLandingZoneDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getProuctsData" --responseType=json`             |
| 10    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getLandingZoneCompliance`    | `go run aws-getLandingZoneCompliance.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getLandingZoneCompliance" --responseType=json` |
| 11    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getElementWiseCompliance`                                | `go run aws-getLandingZoneCompliance.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getElementWiseCompliance" --responseType=json` |
| 12    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getProductWiseCompliance`                                | `go run aws-getLandingZoneCompliance.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getProductWiseCompliance" --responseType=json` |
| 13    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="All",query=getTotalCostOfLandingZone`                                 | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="ALL" --query="getTotalCostOfLandingZone" --responseType=json`            |
| 14    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="123455",query=getTotalCostOfLandingZone`                               | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="123455" --query="getTotalCostOfLandingZone" --responseType=json`          |
| 15    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="123455",query=getCloudElementWiseCostOfLandingZone`                    | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="123455" --query="getCloudElementWiseCostOfLandingZone" --responseType=json` |
| 16    | `/awsx-api/getQueryOutput?elementType=product,productId="ALL",query=getTotalCostOfProduct`                                             | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --elementType=product --productId="ALL" --query="getTotalCostOfProduct" --responseType=json` |
| 17    | `/awsx-api/getQueryOutput?elementType=product,productId="123455",query=getTotalCostOfProduct`                          | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --elementType=product --productId="123455" --query="getTotalCostOfProduct" --responseType=json` |
| 18    | `/awsx-api/getQueryOutput?elementType=product,productId="123455",query=getDetailCostOfProduct`                         | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --productId="123455" --query="getDetailCostOfProduct" --responseType=json`             |
| 19    | `/awsx-api/getQueryOutput?elementType=product,productId="123455",query=getCloudElementWiseCostOfProduct`              | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --productId="123455" --query="getCloudElementWiseCostOfProduct" --responseType=json`      |
| 20   | `/awsx-api/getQueryOutput?elementType=product,productId="123455",query=getCostSpikeOfProduct, frequency=Daily/Weekly/Monthly` | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --productId="123455" --query="getCostSpikeOfProduct" --frequency=Daily/Weekly/Monthly --responseType=json` |
| 21   | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getCostSpikeOfLandingZone, frequency=Daily/Weekly/Monthly` | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getCostSpikeOfLandingZone" --frequency=Daily/Weekly/Monthly --responseType=json` |
| 22    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getAppDataNwWiseCostOfLandingZone`     | `go run aws-getCostDetails.go --vaultUrl=vault.synectiks.net --landingZoneId="12233" --query="getAppDataNwWiseCostOfLandingZone" --responseType=json`  |
| 23    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getInfraToHostedServicesOfLandinZone`  | `go run aws-getTopologyDetails.go --vaultUrl=vault.synectiks.net --elementType=landingZone --landingZoneId="12233" --query="getInfraToHostedServicesOfLandinZone" --responseType=json` |
| 24    | `/awsx-api/getQueryOutput?elementType=landingZone,landingZoneId="12233",query=getServicesToHostingInfraOfLandinZone` | `go run aws-getTopologyDetails.go --vaultUrl=vault.synectiks.net --elementType=landingZone --landingZoneId="12233" -query="getServicesToHostingInfraOfLandinZone" --responseType=json` |
| 25    | `/awsx-api/getQueryOutput?elementType=product,productId="12233",query=getModulesOfTheProduct`                         | `go run aws-getTopologyDetails.go --vaultUrl=vault.synectiks.net --elementType=product --productId="12233" --query="getModulesOfTheProduct" --responseType=json` |
| 26   | `/awsx-api/getQueryOutput?elementType=module,moduleId="12233",query=getAppDataServicesOdModule`                        | `go run aws-getTopologyDetails.go --vaultUrl=vault.synectiks.net --elementType=module --moduleId="12233" --query="getAppDataServicesOdModule" --responseType=json` |
| 27   | `/awsx-api/getQueryOutput?elementType=module,moduleId="12233",query=getAppServicesOdModule`                            | `go run aws-getTopologyDetails.go --vaultUrl=vault.synectiks.net --elementType=module --moduleId="12233" --query="getAppServicesOdModule" --responseType=json` |
| 28   | `/awsx-api/getQueryOutput?elementType=module,moduleId="12233",query=getDataServicesOdModule`                           | `go run aws-getTopologyDetails.go --vaultUrl=vault.synectiks.net --elementType=module --moduleId="12233" --query="getDataServicesOdModule" --responseType=json` |
| 29   |`/awsx-api/getQueryOutput?elementType=EC2,cloudElementId="12233",query=getSLADetailsOfEC2`                           | `go run aws-getSlaDetails.go --vaultUrl=vault.synectiks.net --elementType=module --moduleId="12233" --query="getSLADetailsOfModule" --responseType=json` |




