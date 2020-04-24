package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Shelex/grpc-go-demo/proto"
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
	proto.RegisterEmployeeServiceServer(s, new(employeeService))
	log.Printf("starting server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type employeeService struct{}

func (e *employeeService) GetByBadgeNumber(ctx context.Context, req *proto.GetByBadgeNumberRequest) (*proto.EmployeeResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("Metadata received: %v\n", md)
	}
	log.Printf("requested badge num: %d\n", req.BadgeNumber)
	for _, e := range employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &proto.EmployeeResponse{
				Employee: &e,
			}, nil
		}
	}
	return nil, errors.New("employee not found")
}

func (e *employeeService) GetAll(req *proto.GetAllRequest, stream proto.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		if err := stream.Send(&proto.EmployeeResponse{Employee: &e}); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (e *employeeService) AddPhoto(stream proto.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	var employeeBadge string
	if ok {
		badgeFromContext := md["badgenumber"][0]
		for _, e := range employees {
			if strconv.Itoa(int(e.BadgeNumber)) == badgeFromContext {
				employeeBadge = badgeFromContext
				break
			}
		}
		if employeeBadge == "" {
			return errors.New("employee with provided badge was not found")
		}
		log.Printf("Receiving photo for badge num: %s\n", badgeFromContext)
	}
	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			log.Printf("File received with length: %d\n", len(imgData))
			if err := os.MkdirAll(assets, os.ModePerm); err != nil {
				return err
			}
			filename := fmt.Sprintf("%s/%s-%s.png", assets, employeeBadge, strconv.FormatInt(time.Now().Unix(), 10))
			if err := ioutil.WriteFile(filename, imgData, os.ModePerm); err != nil {
				return err
			}
			return stream.SendAndClose(&proto.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		log.Printf("Received %d bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}

func (e *employeeService) Save(ctx context.Context, req *proto.EmployeeRequest) (*proto.EmployeeResponse, error) {
	log.Printf("before saving employee count is %d\n", len(employees))
	employees = append(employees, *req.Employee)
	log.Printf("now employees count is %d\n", len(employees))
	return &proto.EmployeeResponse{Employee: req.Employee}, nil
}

func (e *employeeService) SaveAll(stream proto.EmployeeService_SaveAllServer) error {
	log.Printf("before saving employees count is %d\n", len(employees))
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, *emp.Employee)
		if err := stream.Send(&proto.EmployeeResponse{
			Employee: emp.Employee,
		}); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("now employees count is %d\n", len(employees))
	return nil
}
