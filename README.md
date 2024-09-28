## Delivery Backend API in Go

### Installation

To run the project, ensure you have the following tools installed on your machine:

- Docker (for containerization)
- Docker Compose

## Running the project

Clone the Repository:
git clone https://github.com/themelancholyspirit/delivery-system
cd delivery-system

Start the Services: 
Make sure to have Docker installed and running. Then, you can start the backend and PostgreSQL services using Docker Compose
docker-compose up

```

## Running the tests

To run the tests, you can use the following command:

```bash
make test
```


## API Endpoints

```

1. Create Order
Method: POST
Endpoint: /orders

Request Body:

{
    "origin": ["START_LATITUDE", "START_LONGITUDE"],
    "destination": ["END_LATITUDE", "END_LONGITUDE"]
}

Example Curl Command:

curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d '{
    "origin": ["52.5200", "13.4050"],
    "destination": ["48.8566", "2.3522"]
}'


2. Get All Orders
Method: GET
Endpoint: /orders?page=:page&limit=:limit

curl -X GET "http://localhost:8080/orders?page=1&limit=10"


3. Update Order Status

Request Body:

{
    "status": "TAKEN"
}

Example Curl Command:

curl -X PATCH http://localhost:8080/orders/1 -H "Content-Type: application/json" -d '{
    "status": "TAKEN"
}'

```
