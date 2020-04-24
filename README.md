# Basic demo of golang grpc 
 Pluralsight course [Enhancing Application Communication with gRPC](https://app.pluralsight.com/library/courses/grpc-enhancing-application-communication/)  
 Server and Client parts of application which communicate via grpc  
 Represents employees, which could be identified by badgeNumber  

# Prerequisites
 - [Go](https://golang.org/)  
 - [GRPC](https://grpc.io/docs/quickstart/go/#grpc)  
 - [Protocol buffers](https://grpc.io/docs/quickstart/go/#protocol-buffers)  

# Install
 - clone this repository  
 - `cd grpc-go-demo`  
 - `make cert` - certificate for TLS, host should be `localhost`  
 - `make build` - build binaries (you can find it in `./cmd/` folder)  
 - `make server` - run server  

# Client options
 - GetByBadgeNumber(badgeNumber): `cmd/client -o 1 -b 1234`  
 get employee from map by badgeNumber  
 - GetAll(): `cmd/client -o 2`  
 get all list of employees  
 - AddPhoto(): `cmd/client -o 3 -b 1234`  
 add photo for employee  
 - SaveAll(): `cmd/client -o 4`  
 bulk insert new employees  
 - Save(): `cmd/client -o 5`  
 save new employee  

 # TODO
  - REST/GraphQL API layer for client
  - persistence layer for server (Mongo/Postgres)
  - expand functionality with calculating and managing employees vacations