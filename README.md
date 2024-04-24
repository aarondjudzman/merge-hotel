# merge-hotel 

This is a simple hotel booking system that allows users to search for hotels and book them. It uses a combination of APIs to fetch hotel data from different sources and merges it into a single data model.

## Table of Contents
1. [API Documentation](#api-documentation)
2. [Run Production Web Server Locally](#run-production-web-server-locally)
3. [Run Development Web Server Locally](#run-development-web-server-locally)
4. [Design Specification](#design-specification)

## API Documentation
### GET `/hotels`

### Query Parameters
- hotels: A comma-separated list of hotel IDs to retrieve. If not provided, all hotels are returned regardless of ID.
- destination: The ID of the destination to retrieve hotels for. If not provided, all hotels are returned regardless of destination.

### Example Request
```
GET /hotels?hotels=iJhz,SjyX&destination=5432
```

## Run production web server locally 
- Clone the repository to your local machine using the following command:
```
git clone https://github.com/ascenda/merge-hotel.git
```
- Change directory to the cloned repository:
```
cd merge-hotel
```

- If you already have Go installed, you can build and run the project using the following command so it runs on your local machine:
```
go build -o merge-hotel .
./merge-hotel
```
- If you don't have Go installed, you can find pre-built binary to run the project on your local machine based on your environment in the `dist` folder. E.g for ARM64 MacOS:
```
./dist/merge-hotel_darwin_arm64/merge-hotel
```
- The project will start on localhost:8080. You can make this example request to get a list of hotels:
```
curl http://localhost:8080/hotels
```

## Run development web server locally
- Install Go using the official instructions from the [Go website](https://go.dev/doc/install).
- Clone the repository to your local machine using the following command:
```
git clone https://github.com/ascenda/merge-hotel.git
```
- Change directory to the cloned repository:
```
cd merge-hotel
```
- Get the dependencies using the following command:
```
go get
```
- Run the project using the following command:
```
go run .
```
- The project will start on localhost:8080. You can make this example request to get a list of hotels:
```
curl http://localhost:8080/hotels
```

## Design Specification
### Data Model
The business logic uses a common data model to represent hotels. The data model is based on the default response format of the task specificiation. It is specified in the `entity/hotel.go` file.

### Data supplier
Each supplier has its own parser to convert the response format to the common data model. The design choice for each supplier is outlined as comments in each supplier file in the `supplier` folder.

### Usecase
The usecase is a simple implementation of the business logic. It uses the data model and the data supplier to fetch hotels from the APIs and merge them into a single data model. The choice of selecting which data to keep to deliver the most complete data set is outlined in the `entity/merger.go` file.

### Optimisation
- The usecase concurrently fetches data from all suppliers instead of sequentially.
- The usecase uses a cache to store the results of the previous call to avoid fetching the same data multiple times within short time intervals.

### Possible Further Improvements
- Use distributed cache.
- Use retry mechanism when fetching data from suppliers.
