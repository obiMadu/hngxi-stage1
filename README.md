# HNG XI STAGE 1 PROJECT

A dead-simple API that works with the below Request and Response formats;

## Request format
Call to Endpoint: GET `{endpoint}/api/hello?visitor_name={visitor_name}`
 - `{endpoint}` refers to the servers base url
 - `{visitor_name}` refers to a string query

## Response format
```json
{
    "client_ip": "{client_ip}",
    "location": "{client IP location_city}",
    "greeting": "Hello, {visitor_name}!, the temperature is {current_temperature_of_client_ip_location} degrees Celcius in {client IP location_city}"
}
```