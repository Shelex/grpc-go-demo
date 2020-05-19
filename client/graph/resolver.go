package graph

import (
	"github.com/Shelex/grpc-go-demo/proto"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	employeeServiceClient proto.EmployeeServiceClient
}

func NewResolver(employeeServiceClient proto.EmployeeServiceClient) *Resolver {
	return &Resolver{
		employeeServiceClient: employeeServiceClient,
	}

}
