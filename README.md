# Basic demo of golang grpc 
 Based on Pluralsight course [Enhancing Application Communication with gRPC](https://app.pluralsight.com/library/courses/grpc-enhancing-application-communication/)  
 Server and Client parts of application which communicate via grpc  
 Client side exposes graphql API  
 Represents employees, with ability to upload files and creating vacation requests  

# Prerequisites
 - [Go](https://golang.org/)  
 - [GRPC](https://grpc.io/docs/quickstart/go/#grpc)  
 - [Protocol buffers](https://grpc.io/docs/quickstart/go/#protocol-buffers)  
 - [MongoDB](https://www.mongodb.com/download-center/community)  
 - For testing file upload use [Altair](https://altair.sirmuel.design/) graphql client  

# Install
 - clone this repository  
 - `cd grpc-go-demo`  
 - `make cert` - certificate for TLS, host should be `localhost`  
 - `make build` - build binaries (you can find it in `./cmd/` folder)  
 - `make server` - run domain-server  
 - `make client` - run api-client  
 - open `http://localhost:8080/` for GraphQL playground or use Altair/Postman/Insomnia  

# Develop
 - `make gql` - regenerate graphql client  
 - `make gen` - regenerate proto files

# Client options

### For testing file uploading Altair graphql client is recommended
 - query employees: get list of all employees
```graphql
 query {
  employees {
    id
    firstName
    lastName
    badgeNumber
    countryCode
    vacationAccrualRate
    vacationAccrued
    documents
    vacations {
      id
      startDate
      durationHours
      approved
      cancelled
    }
  }
}
 ```

 - query employees: get list of all employees
```graphql
query {
  employee(id: "775bc400-9b23-4b40-b336-d591f842f16b") {
    id
    badgeNumber
    firstName
    lastName
    countryCode
    vacationAccrualRate
    vacationAccrued
    documents
    vacations {
      id
      startDate
      durationHours
      approved
      cancelled
    }
  }
}
 ```

  - query attachment(ID): get attachment by ID with bytes in property `data`
```
query {
  attachment(id: "e9463e2c-9f8e-11ea-91cf-f01898689eb3") {
    id
    userID
    fileName
    data
    createdAt
  }
}
```

- mutation addEmployee(employee): create employee
```graphql
 mutation {
  addEmployee(
    employee: {
      badgeNumber: 6062
      firstName: "Druzhochek"
      countryCode: "UA"
      lastName: "Pirazhochek"
      vacationAccrualRate: 0.245
      vacationAccrued: 0
    }
  ) {
    id
    badgeNumber
    firstName
    lastName
    countryCode
    vacationAccrualRate
    vacationAccrued
    documents
    vacations {
      id
      startDate
      durationHours
      approved
      cancelled
    }
  }
}
```

- mutation addEmployees(employees): create multiple employee
```graphql
 mutation {
  addEmployees(
    employees: [
      {
        badgeNumber: 203
        firstName: "Druzhochek"
        lastName: "Pirazhochek"
        countryCode: "LS"
        vacationAccrualRate: 0.245
        vacationAccrued: 0
      }
      {
        badgeNumber: 204
        firstName: "Podruzhka"
        lastName: "Vatrushka"
        countryCode: "AS"
        vacationAccrualRate: 1.5
        vacationAccrued: 2
      }
    ]
  ) {
    savedEmployees {
      id
      badgeNumber
      firstName
      lastName
      vacationAccrualRate
      vacationAccrued
      documents
    }
    errors
  }
}
```

- mutation addAttachment(file, userID): create attachment for user
  ```json:
  {
    "id": "775bc400-9b23-4b40-b336-d591f842f16b"
  }
  ```
  ```graphql
  mutation($image: Upload!, $id: String!) {
  addAttachment(userID: $id, file: $image) {
    id
    userID
    fileName
    data
    createdAt
    }
  }``` 


- mutation updateEmployee(employee): update employee fields
```graphql
  mutation {
  updateEmployee(
    userID: "775bc400-9b23-4b40-b336-d591f842f16b"
    employee: {
      firstName: "John"
      badgeNumber: 1010
      lastName: "Doe"
      countryCode: "US"
      vacationAccrualRate: 2.356
      vacationAccrued: 2
    }
  ) {
    id
    badgeNumber
    firstName
    lastName
    countryCode
    vacationAccrualRate
    vacationAccrued
    vacations {
      id
      startDate
      durationHours
      approved
      cancelled
    }
    documents
  }
}
  ```
  

- mutation deleteEmployee(id): delete employee
```graphql
   mutation {
  deleteEmployee(userID: "41578cef-3e2f-4b3d-8982-19fa0328f8b9") {
    id
    firstName
    lastName
    badgeNumber
  }
}
```


- mutation addVacation(vacation): create vacation request for employee
```graphql
  mutation {
  addVacation(
    vacation: {
      userID: "775bc400-9b23-4b40-b336-d591f842f16b"
      startDate: 1800187289
      durationHours: 48
    }
  ) {
    id
    startDate
    durationHours
    approved
    cancelled
  }
}
```


 # TODO
  - add checks for add, update employee queries checks for badgeNum duplication
  - persistence layer for server (Mongo/Postgres/S3 for files)
  - add graceful shutdown
  - logging, better error handling, add wrapping where needed
  - dockerize
  - deploy (to Heroku, for example)