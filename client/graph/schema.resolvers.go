package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Shelex/grpc-go-demo/client/graph/factory"
	"github.com/Shelex/grpc-go-demo/client/graph/generated"
	"github.com/Shelex/grpc-go-demo/client/graph/model"
	"github.com/Shelex/grpc-go-demo/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

func (r *mutationResolver) AddEmployee(ctx context.Context, employee model.EmployeeInput) (*model.Employee, error) {
	log.WithField("employee", employee).Info("requesting AddEmployee")
	newEmployee := factory.EmployeeFromAPIToProto(employee)
	res, err := r.employeeServiceClient.AddEmployee(context.Background(), &proto.EmployeeRequest{Employee: newEmployee})
	if err != nil {
		return nil, err
	}
	return factory.EmployeeFromProtoToApi(res.Employee), nil
}

func (r *mutationResolver) AddEmployees(ctx context.Context, employees []*model.EmployeeInput) (*model.EmployeeSaveResult, error) {
	stream, err := r.employeeServiceClient.AddEmployees(context.Background())
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
		stringified := err.Error()
		result.Errors = &stringified
	}
	return result, err
}

func (r *mutationResolver) AddAttachment(ctx context.Context, userID string, file graphql.Upload) (*model.Document, error) {
	log.Printf("got file: %s with size: %d, and CT:%s", file.Filename, file.Size, file.ContentType)
	md := metadata.New(map[string]string{"userID": userID, "filename": file.Filename})
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := r.employeeServiceClient.AddAttachment(ctx)
	if err != nil {
		return nil, err
	}
	for {
		chunk := make([]byte, 64*1024) // 64kb chunk
		n, err := file.File.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		if err := stream.Send(&proto.AttachmentRequest{UserID: userID, Filename: file.Filename, Data: chunk}); err != nil {
			return nil, err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return factory.AttachmentFromProtoToApi(res), nil
}

func (r *mutationResolver) UpdateEmployee(ctx context.Context, userID string, employee model.EmployeeInput) (*model.Employee, error) {
	res, err := r.employeeServiceClient.UpdateEmployee(ctx, &proto.EmployeeUpdateRequest{
		ID: userID,
		Updates: &proto.EmployeeRequest{
			Employee: factory.EmployeeFromAPIToProto(employee),
		},
	})
	if err != nil {
		return nil, err
	}
	return factory.EmployeeFromProtoToApi(res.Employee), nil
}

func (r *mutationResolver) DeleteEmployee(ctx context.Context, userID string) (*model.Employee, error) {
	res, err := r.employeeServiceClient.DeleteEmployee(ctx, &proto.ByIDRequest{
		ID: userID,
	})
	if err != nil {
		return nil, err
	}
	return factory.EmployeeFromProtoToApi(res.Employee), nil
}

func (r *mutationResolver) AddVacation(ctx context.Context, vacation model.VacationRequest) (*model.Vacation, error) {
	newVacation, err := r.employeeServiceClient.AddVacation(ctx, &proto.VacationRequest{
		UserID:        vacation.UserID,
		StartDate:     int64(vacation.StartDate),
		DurationHours: float32(vacation.DurationHours),
	})
	if err != nil {
		return nil, err
	}
	return factory.VacationFromProtoToApi(newVacation), nil
}

func (r *queryResolver) Employees(ctx context.Context) ([]*model.Employee, error) {
	log.Info("requesting Employees")
	stream, err := r.employeeServiceClient.Employees(context.Background(), &proto.GetAllRequest{})
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

func (r *queryResolver) Employee(ctx context.Context, id string) (*model.Employee, error) {
	md := metadata.New(map[string]string{"userID": id})
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, err := r.employeeServiceClient.EmployeeByID(ctx, &proto.ByIDRequest{ID: id})
	if err != nil {
		return nil, err
	}
	return factory.EmployeeFromProtoToApi(res.Employee), nil
}

func (r *queryResolver) Attachment(ctx context.Context, id string) (*model.Document, error) {
	doc, err := r.employeeServiceClient.AttachmentByID(ctx, &proto.ByIDRequest{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	return factory.AttachmentFromProtoToApi(doc), nil
}

func (r *queryResolver) Vacations(ctx context.Context) ([]*model.Vacation, error) {
	stream, err := r.employeeServiceClient.Vacations(context.Background(), &proto.GetAllRequest{})
	if err != nil {
		return nil, err
	}
	vacations := make([]*model.Vacation, 0, 10)
	for {
		vacation, err := stream.Recv()
		if err == io.EOF {
			break
		}
		vacations = append(vacations, factory.VacationFromProtoToApi(vacation))
	}
	return vacations, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
