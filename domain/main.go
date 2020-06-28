package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/Shelex/grpc-go-demo/domain/entities"
	"github.com/Shelex/grpc-go-demo/proto"
	"github.com/Shelex/grpc-go-demo/storage"
	"github.com/Shelex/grpc-go-demo/storage/documents"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	port      = ":9000"
	MONGO_URL = "localhost:27017"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	var repo storage.Storage
	var env string
	var dbErr error
	env = os.Getenv("ENV")
	switch env {
	case "test":
		repo, dbErr = storage.NewInMemStorage()
	default:
		repo, dbErr = storage.NewMongoStorage(MONGO_URL)
	}
	if dbErr != nil {
		log.Fatalf("failed to connect to repository: %s", dbErr)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)

	srv := &employeeService{
		repository: repo,
		documents:  documents.NewLocalFS(),
	}
	proto.RegisterEmployeeServiceServer(s, srv)
	log.Printf("starting server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type employeeService struct {
	repository storage.Storage
	documents  documents.FileStorage
}

func (e *employeeService) Employees(req *proto.GetAllRequest, stream proto.EmployeeService_EmployeesServer) error {
	employees, err := e.repository.GetAll()
	if err != nil {
		return err
	}
	for _, emp := range employees {
		if err := stream.Send(&proto.EmployeeResponse{Employee: entities.EmployeeFromStorageToProto(emp)}); err != nil {
			return err
		}
	}
	return nil
}

func (e *employeeService) EmployeeByID(ctx context.Context, req *proto.ByIDRequest) (*proto.EmployeeResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("metadata received: %v\n", md)
	}
	employee, err := e.repository.GetEmployee(req.ID)
	if err != nil {
		return nil, err
	}
	return &proto.EmployeeResponse{
		Employee: entities.EmployeeFromStorageToProto(employee),
	}, nil
}

func (e *employeeService) AddEmployee(ctx context.Context, req *proto.EmployeeRequest) (*proto.EmployeeResponse, error) {
	req.Employee.ID = uuid.New().String()
	employee, err := e.repository.AddEmployee(entities.EmployeeFromProtoToStorage(req.Employee))
	if err != nil {
		return nil, err
	}
	count, _ := e.repository.Count()
	log.Printf("employee %s successfully saved; now have %d\n", employee.ID, count)
	return &proto.EmployeeResponse{Employee: entities.EmployeeFromStorageToProto(employee)}, nil
}
func (e *employeeService) AddEmployees(stream proto.EmployeeService_AddEmployeesServer) error {
	initialCount, err := e.repository.Count()
	if err != nil {
		return err
	}
	log.Printf("now have %d employees\n", initialCount)
	var savedCount int
	var errorMessage string
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		emp.Employee.ID = uuid.New().String()
		saved, err := e.repository.AddEmployee(entities.EmployeeFromProtoToStorage(emp.Employee))
		if err != nil {
			errorMessage += fmt.Sprintf("\n%s", err.Error())
			continue
		}
		if err := stream.Send(&proto.EmployeeResponse{
			Employee: entities.EmployeeFromStorageToProto(saved),
		}); err != nil {
			return err
		}
		savedCount++
	}
	if errorMessage != "" {
		return errors.New(errorMessage)
	}
	current, err := e.repository.Count()
	if err != nil {
		return err
	}
	log.Printf("successfully saved %d employees;\n now have %d\n", savedCount, current)
	return nil
}

func (e *employeeService) AddAttachment(stream proto.EmployeeService_AddAttachmentServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	var userID string
	var fileName string
	if ok {
		userID = md.Get("userID")[0]
		fileName = md.Get("filename")[0]
		emp, err := e.repository.GetEmployee(userID)
		if err != nil {
			return err
		}
		log.Printf("receiving photo for user %s (%s %s)\n", emp.ID, emp.FirstName, emp.LastName)
	}
	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			log.Printf("file received with length: %d\n", len(imgData))
			document, err := e.documents.SaveDocument(userID, fileName, imgData)
			document.UserID = userID
			if err != nil {
				return err
			}
			if err := e.repository.AddDocument(userID, document.ID); err != nil {
				return err
			}
			return stream.SendAndClose(entities.DocumentFromStorageToProto(document))
		}
		if err != nil {
			return err
		}
		log.Printf("received %d bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}

func (e *employeeService) AttachmentByID(ctx context.Context, req *proto.ByIDRequest) (*proto.Attachment, error) {
	doc, err := e.documents.GetDocument(req.ID)
	if err != nil {
		return nil, err
	}
	return entities.DocumentFromStorageToProto(doc), nil
}

func (e *employeeService) DeleteEmployee(ctx context.Context, req *proto.ByIDRequest) (*proto.EmployeeResponse, error) {
	employee, err := e.repository.DeleteEmployee(req.ID)
	if err != nil {
		return nil, err
	}
	return &proto.EmployeeResponse{
		Employee: entities.EmployeeFromStorageToProto(employee),
	}, nil
}

func (e *employeeService) UpdateEmployee(ctx context.Context, req *proto.EmployeeUpdateRequest) (*proto.EmployeeResponse, error) {
	employee, err := e.repository.UpdateEmployee(req.ID, entities.EmployeeFromProtoToStorage(req.Updates.Employee))
	if err != nil {
		return nil, err
	}
	return &proto.EmployeeResponse{
		Employee: entities.EmployeeFromStorageToProto(employee),
	}, nil
}

func (e *employeeService) AddVacation(ctx context.Context, req *proto.VacationRequest) (*proto.Vacation, error) {
	start := time.Unix(req.StartDate, 0)
	hours24 := time.Now().Add(24 * time.Hour)

	after24Hours := start.After(hours24)
	if !after24Hours {
		return nil, fmt.Errorf("vacation start date should be not less than 24 hours from now")
	}

	vacationID := uuid.New().String()

	vacation, err := e.repository.AddVacation(vacationID, req.UserID, req.StartDate, req.DurationHours)
	if err != nil {
		return nil, err
	}
	return entities.VacationFromStorageToProto(vacation), nil
}

func (e *employeeService) Vacations(req *proto.GetAllRequest, stream proto.EmployeeService_VacationsServer) error {
	vacations, err := e.repository.Vacations()
	if err != nil {
		return err
	}
	for _, vacation := range vacations {
		if err := stream.Send(entities.VacationFromStorageToProto(vacation)); err != nil {
			return err
		}
	}
	return nil
}
