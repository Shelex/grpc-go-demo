package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Shelex/grpc-go-demo/entities"
	"github.com/Shelex/grpc-go-demo/proto"
	"github.com/Shelex/grpc-go-demo/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	port   = ":9000"
	assets = "imageRepository"
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

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)
	srv := &employeeService{
		repository: storage.NewInMemStorage(),
	}
	proto.RegisterEmployeeServiceServer(s, srv)
	log.Printf("starting server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type employeeService struct {
	repository storage.Storage
}

func (e *employeeService) GetByBadgeNumber(ctx context.Context, req *proto.GetByBadgeNumberRequest) (*proto.EmployeeResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("metadata received: %v\n", md)
	}
	log.Printf("requested badge num: %d\n", req.BadgeNumber)
	employee, err := e.repository.GetByBadge(req.BadgeNumber)
	if err != nil {
		return nil, err
	}
	return &proto.EmployeeResponse{
		Employee: entities.EmployeeFromStorageToProto(employee),
	}, nil
}

func (e *employeeService) GetAll(req *proto.GetAllRequest, stream proto.EmployeeService_GetAllServer) error {
	employees, err := e.repository.GetAll()
	if err != nil {
		return err
	}
	for _, e := range employees {
		if err := stream.Send(&proto.EmployeeResponse{Employee: entities.EmployeeFromStorageToProto((e))}); err != nil {
			return err
		}
	}
	return nil
}

func (e *employeeService) AddPhoto(stream proto.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	var employeeBadge int32
	if ok {
		badgeFromContext := md["badgenumber"][0]
		badgeNum, err := strconv.Atoi(badgeFromContext)
		if err != nil {
			return err
		}
		employee, err := e.repository.GetByBadge(int32(badgeNum))
		if err != nil {
			return err
		}
		employeeBadge = employee.BadgeNumber
		log.Printf("receiving photo for badge num: %d\n", employeeBadge)
	}
	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			log.Printf("file received with length: %d\n", len(imgData))
			if err := os.MkdirAll(assets, os.ModePerm); err != nil {
				return err
			}
			filename := fmt.Sprintf("%s/%d-%s.png", assets, employeeBadge, strconv.FormatInt(time.Now().Unix(), 10))
			if err := ioutil.WriteFile(filename, imgData, os.ModePerm); err != nil {
				return err
			}
			return stream.SendAndClose(&proto.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		log.Printf("received %d bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}

func (e *employeeService) Save(ctx context.Context, req *proto.EmployeeRequest) (*proto.EmployeeResponse, error) {
	if err := e.repository.AddEmployee(entities.EmployeeFromProtoToStorage(req.Employee)); err != nil {
		return nil, err
	}
	count, _ := e.repository.Count()
	log.Printf("employee with badge %d successfully saved; now have %d\n", req.Employee.BadgeNumber, count)
	return &proto.EmployeeResponse{Employee: req.Employee}, nil
}

func (e *employeeService) SaveAll(stream proto.EmployeeService_SaveAllServer) error {
	initialCount, err := e.repository.Count()
	if err != nil {
		return err
	}
	log.Printf("now have %d employees\n", initialCount)
	var savedCount int
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := e.repository.AddEmployee(entities.EmployeeFromProtoToStorage(emp.Employee)); err != nil {
			return err
		}
		if err := stream.Send(&proto.EmployeeResponse{
			Employee: emp.Employee,
		}); err != nil {
			return err
		}
		savedCount++
	}
	current, err := e.repository.Count()
	if err != nil {
		return err
	}
	log.Printf("successfully saved %d employees; now have %d\n", savedCount, current)
	return nil
}
