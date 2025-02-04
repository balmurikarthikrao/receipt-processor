# Receipt Processor

## Overview
The Receipt Processor is a Go application that processes receipts and calculates points based on predefined rules. It provides a RESTful API for submitting receipts and retrieving points awarded for each receipt.

## Project Structure
```
receipt-processor
├── controllers          # Contains HTTP request handlers
│   └── receipt_controller.go
├── models               # Defines the data models
│   └── receipt.go
├── services             # Implements business logic
│   └── receipt_service.go
├── main.go              # Entry point of the application
├── go.mod               # Module definition
├── go.sum               # Dependency checksums
└── README.md            # Project documentation
```

## Setup Instructions
1. Clone the repository:
   ```
   git clone <repository-url>
   cd receipt-processor
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run main.go
   ```

## API Endpoints
### Process Receipt
- **Endpoint:** `/receipts/process`
- **Method:** `POST`
- **Payload:** JSON object representing the receipt.
- **Response:** JSON object containing the receipt ID.

### Get Points
- **Endpoint:** `/receipts/{id}/points`
- **Method:** `GET`
- **Response:** JSON object containing the number of points awarded.

## Usage Examples
### Process Receipt Example
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    }
  ],
  "total": "35.35"
}
```

### Get Points Example
- Request: `GET /receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points`
- Response: 
```json
{
  "points": 28
}
```

## How to run this application

Build the Docker image:
```shell
docker build -t receipt-processor .
```

Run the Docker container:
```shell
docker run -p 8080:8080 receipt-processor
```