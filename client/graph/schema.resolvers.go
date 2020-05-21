package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"io"
	"log"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Shelex/grpc-go-demo/client/graph/factory"
	"github.com/Shelex/grpc-go-demo/client/graph/generated"
	"github.com/Shelex/grpc-go-demo/client/graph/model"
	"github.com/Shelex/grpc-go-demo/proto"
	"google.golang.org/grpc/metadata"
)

func (r *mutationResolver) AddEmployeeAttachment(ctx context.Context, file graphql.Upload, badgeNumber int) (bool, error) {
	log.Printf("got file: %s with size: %d, and CT:%s", file.Filename, file.Size, file.ContentType)
	md := metadata.New(map[string]string{"badgenumber": strconv.Itoa(badgeNumber)})
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := r.employeeServiceClient.AddEmployeeAttachment(ctx)
	if err != nil {
		return false, err
	}
	for {
		chunk := make([]byte, 64*1024) // 64kb chunk
		n, err := file.File.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		if err := stream.Send(&proto.AddAttachmentRequest{Data: chunk}); err != nil {
			return false, err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return false, err
	}
	return res.IsOk, nil
}

func (r *mutationResolver) Save(ctx context.Context, employee model.EmployeeInput) (*model.Employee, error) {
	newEmployee := factory.EmployeeFromAPIToProto(employee)
	res, err := r.employeeServiceClient.Save(context.Background(), &proto.EmployeeRequest{Employee: newEmployee})
	if err != nil {
		return nil, err
	}
	return factory.EmployeeFromProtoToApi(res.Employee), nil
}

func (r *mutationResolver) SaveAll(ctx context.Context, employees []*model.EmployeeInput) (*model.EmployeeSaveResult, error) {
	stream, err := r.employeeServiceClient.SaveAll(context.Background())
	if err != nil {
		return nil, err
	}
	done := make(chan error)
	result := &model.EmployeeSaveResult{}
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				done <- nil
				break
			}
			if err != nil {
				done <- err
				continue
			}
			result.SavedEmployees = append(result.SavedEmployees, factory.EmployeeFromProtoToApi(res.Employee))
		}
	}()
	for _, e := range employees {
		if err := stream.Send(&proto.EmployeeRequest{
			Employee: factory.EmployeeFromAPIToProto(*e),
		}); err != nil {
			return result, err
		}
	}
	if err := stream.CloseSend(); err != nil {
		return result, err
	}
	if err := <-done; err != nil {
		result.Error = err.Error()
	}
	return result, err
}

func (r *queryResolver) GetByBadge(ctx context.Context, badgeNumber int) (*model.Employee, error) {
	md := metadata.MD{}
	md["user"] = []string{"THISISUSER"}
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, err := r.employeeServiceClient.GetByBadgeNumber(ctx, &proto.GetByBadgeNumberRequest{BadgeNumber: int32(badgeNumber)})
	if err != nil {
		return nil, err
	}
	return factory.EmployeeFromProtoToApi(res.Employee), nil
}

func (r *queryResolver) GetAll(ctx context.Context) ([]*model.Employee, error) {
	stream, err := r.employeeServiceClient.GetAll(context.Background(), &proto.GetAllRequest{})
	if err != nil {
		return nil, err
	}
	employees := make([]*model.Employee, 0, 10)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		employees = append(employees, factory.EmployeeFromProtoToApi(res.Employee))
	}
	return employees, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
