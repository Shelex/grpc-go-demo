# Basic demo of golang grpc 
 Based on Pluralsight course [Enhancing Application Communication with gRPC](https://app.pluralsight.com/library/courses/grpc-enhancing-application-communication/)  
 Server and Client parts of application which communicate via grpc  
 Client side exposes graphql API  
 Represents employees, which could be identified by badgeNumber  

# Prerequisites
 - [Go](https://golang.org/)  
 - [GRPC](https://grpc.io/docs/quickstart/go/#grpc)  
 - [Protocol buffers](https://grpc.io/docs/quickstart/go/#protocol-buffers)  
 - For testing file upload use [Altair](https://altair.sirmuel.design/) graphql client  

# Install
 - clone this repository  
 - `cd grpc-go-demo`  
 - `make cert` - certificate for TLS, host should be `localhost`  
 - `make gql` - generate graphql client  
 - `make build` - build binaries (you can find it in `./cmd/` folder)  
 - `make server` - run server  
 - `make client` - run client  

# Client options
 Graphql playground will be exposed and such queries and mutations are available:
 - query `getAll`
```
query {
  getAll {
    Id
    BadgeNumber
    FirstName
    LastName
    Documents
  }
}

```
  - query `getByBadge`, accepts badgeNumber as integer
```
query {
  getByBadge(badgeNumber: 6238) {
    Id
    FirstName
    LastName
    BadgeNumber
    Documents
  }
}
```
 - mutation `Save` to add new Employee
 ```
 mutation {
  Save(
    employee: {
      Id: 25
      BadgeNumber: 4040
      FirstName: "Druzhochek"
      LastName: "Pirazhochek"
      VacationAccrualRate: 0.245
    }
  ) {
    Id
    FirstName
    LastName
    BadgeNumber
    VacationAccrualRate
    VacationAccrued
    Vacations {
      Id
      IsCancelled
      StartDate
      __typename
    }
  }
}
```
 - mutation `saveAll` - bulk insert several employees
 ```
 mutation {
  SaveAll(
    employees: [
      {
        Id: 28
        BadgeNumber: 105
        FirstName: "Druzhochek"
        LastName: "Pirazhochek"
        VacationAccrualRate: 0.245
      }
      {
        Id: 26
        BadgeNumber: 104
        FirstName: "Podruzhka"
        LastName: "Vatrushka"
        VacationAccrualRate: 1.5
      }
    ]
  ) {
    savedEmployees {
      Id
      BadgeNumber
      FirstName
      LastName
      VacationAccrualRate
      VacationAccrued
    }
    error
  }
}

```
 - mutation `AddEmployeeAttachment` - to add some photo for employee
 ```
 mutation ($image: Upload!, $num: Int!) {
  AddEmployeeAttachment (file: $image, badgeNumber: $num)
}

variables: {
  "num": 7975
}

// to upload file Altair graphql client could be used
 ```

 # TODO
  - persistence layer for server (Mongo/Postgres)
  - add graceful shutdown
  - logging, better error handling, add wrapping where needed
  - dockerize
  - deploy (to Heroku, for example)