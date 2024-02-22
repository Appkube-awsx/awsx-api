 - [awsx-api](#awsx-api)

- [awsx ecs cpu utilization panel api](#awsx-ec2-cpu-utilization-panel-api)

   - [overview](#overview)
   - [api endpoint](#api-endpoint)
   - [api reference](#api-reference)
   - [errors](#errors)
   - [https status code summary](#https-status-code-summary)

- [curl command](#curl-command)
- [output](#output)



# awsx-api
awsx-api is a golang based REST api server that exposes GET, POST, DELETE and PUT endpoints that will subsequently allow us to perform the full range of operations on AWS entities.


#  awsx ec2 cpu utilization panel api

## overview
This Go code defines an HTTP handler function for retrieving CPU utilization metrics for EC2 instances in the AWS cloud. The API supports both direct authentication using AWS credentials and cross-account authentication.

## api endpoint
- **Endpoint:** `/awsx-api/getQueryOutput`
- **HTTP Method:** `GET`
This markdown file contains all api document Order-wise how does flow works of EC2 CPU Utilization Panel

	baseMetricUrl:
		http://localhost:7000

## api reference
The EC2 CPU Utilization Panel is organized around REST. Our API has predictable resource-oriented URLs, accepts form-encoded request bodies, returns JSON-encoded responses, and uses standard HTTP response codes, authentication, and verbs.

## errors

EC2 CPU Utilization Panel uses conventional HTTP response codes to indicate the success or failure of an API request. In general: Codes in the 2xx range indicate success. Codes in the 4xx range indicate an error that failed given the information provided (e.g., a required parameter was omitted, a charge failed, etc.). Codes in the 5xx range indicate an error with Stripe's servers (these are rare).

Some 4xx errors that could be handled programmatically (e.g., a card is declined) include an error code that briefly explains the error reported.

 ## https status code summary

Code   | Summary
------------- | -------------
200 - OK  | Everything worked as expected.
400 - Bad Request  | The request was unacceptable, often due to missing a required parameter.
401 - Unauthorized | No valid API key provided.
402 - Request Failed | The parameters were valid but the request failed.
403 - Forbidden | The API key doesn't have permissions to perform the request.
404 - Not Found | The requested resource doesn't exist.
409 - Conflict | The request conflicts with another request (perhaps due to using the same idempotent key)
429 - Too Many Requests | Too many requests hit the API too quickly. We recommend an exponential backoff of your requests.
500, 502, 503, 504 - Server Errors | Something went wrong on Stripe's end. (These are rare.)

## request parameters
- `zone`: AWS region or availability zone.
- `cloudElementId`: ID of the AWS cloud element (optional).
- `cloudElementApiUrl`: URL for AWS cloud element API (optional).
- `instanceID`: ID of the EC2 instance.
- `elementType`: Type of the AWS element (e.g., EC2).
- `query`: Query string for metric data.
- `startTime`: Start time for the metric data retrieval (optional).
- `endTime`: End time for the metric data retrieval (optional).
- `statistic`: Statistic types for metric data (comma-separated, optional).
- `crossAccountRoleArn`: Cross-account role ARN for authentication (if using cross-account authentication).
- `externalId`: External ID for cross-account authentication.

# curl command 
```
http://localhost:7000/awsx-api/getQueryOutput?vaultUrl=<afreenxxxx1309>&elementType=EC2&elementId=900000&query=cpu_utilization_panel&responsetype=json&endTime=2023-12-02T23%3A59%3A59Z&startTime=2023-12-01T00%3A00%3A00Z
```

## output
```
{"AverageUsage":20.0125,"CurrentUsage":4,"MaxUsage":23.075}

```


## appkube-platform (grafana) documentation


#### Purpose

The `testAppkubeCputUtilization` function is a Go function designed to query CPU utilization data using the AppKube API and Infinity client.

#### Parameters

- `zone` (string): The AWS region/zone to query (e.g., "us-east-1").
- `vaultUrl` appears to be a query parameter used to specify the URL of a vault
- `elementType` (string): Type of the AWS resource (e.g., "AWS/EC2").
- `instanceID` (string): ID of the specific AWS instance.
- `query` (string): Query type, in this case, "CPUUtilization".
- `statistic` (string): Statistic type, in this case, "SampleCount".

#### Usage

```
// Example Usage
testAppkubeCputUtilization()
```
#### Dependencies
- Infinity package for creating an Infinity client.
- The pluginhost and backend packages for querying data and plugin context.

#### Error Handling

If an error occurs while creating the Infinity client, an error message will be printed to the console.

#### Notes

- Ensure that the necessary dependencies are installed before using the function.
- This function assumes a specific JSON format for the query and sends it to the AppKube API.

Example JSON Query
```
{
    "type": "appkube-api",
    "source": "url",
    "productId": 1,
    "environmentId": 2,
    "moduleId": 2,
    "serviceId": 2,
    "serviceType": "java app service",
    "zone": "us-east-1",
    "externalId": "657907747545",
    "crossAccountRoleArn": "arn:aws:iam::657907747545:role/CrossAccount",
    "elementType": "AWS/EC2",
    "instanceID": "i-05e4e6757f13da657",
    "query": "CPUUtilization",
    "statistic": "SampleCount"
}

```

#### Response

The function prints the response frames obtained from the AppKube API.

```
fmt.Println("Response: ", res.Frames)

```

# For more Details go through the below link

for ec2 check this git repo:
     
git clone https://github.com/Appkube-awsx/awsx/blob/main/ec2Document.md

for grafana check this git repo

   git clone https://github.com/Appkube-awsx/awsx/blob/main/Grafana%20Document.md
