package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Shelex/grpc-go-demo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	option := flag.Int("o", 1, "invoke grpc command")
	badgeNumber := flag.Int("b", 1, "badgeNumber")
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewEmployeeServiceClient(conn)
	switch *option {
	case 1:
		GetByBadgeNumber(client, badgeNumber)
	case 2:
		GetAll(client)
	case 3:
		AddPhoto(client, badgeNumber)
	case 4:
		SaveAll(client)
	case 5:
		Save(client)
	}
}

// save single employee
func Save(client proto.EmployeeServiceClient) {
	newEmployee := proto.Employee{
		Id:                  8,
		BadgeNumber:         3310,
		FirstName:           "Wayne",
		LastName:            "Lewter",
		VacationAccrualRate: 2.76,
		VacationAccrued:     12.1,
	}
	res, err := client.Save(context.Background(), &proto.EmployeeRequest{Employee: &newEmployee})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved employee: %v\n", res.Employee)
}

// bulk insert employee
func SaveAll(client proto.EmployeeServiceClient) {
	newEmployees := []proto.Employee{
		{
			Id:                  6,
			BadgeNumber:         4748,
			FirstName:           "Samuel",
			LastName:            "Weldon",
			VacationAccrualRate: 2.76,
			VacationAccrued:     12.1,
		},
		{
			Id:                  7,
			BadgeNumber:         2776,
			FirstName:           "Geraldine",
			LastName:            "Foster",
			VacationAccrualRate: 1.78,
			VacationAccrued:     33.2,
		},
	}
	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				done <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.Employee)
		}
	}()
	for _, e := range newEmployees {
		if err := stream.Send(&proto.EmployeeRequest{
			Employee: &e,
		}); err != nil {
			log.Fatal(err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatal(err)
	}
	<-done
}

// add file for specific employee
func AddPhoto(client proto.EmployeeServiceClient, badgeNumber *int) {
	f, err := os.Open("img.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	md := metadata.New(map[string]string{"badgenumber": strconv.Itoa(*badgeNumber)})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for {
		chunk := make([]byte, 64*1024) // 64kb chunk
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		if err := stream.Send(&proto.AddPhotoRequest{Data: chunk}); err != nil {
			log.Fatal(err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.IsOk)
}

func GetAll(client proto.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &proto.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("Got employee: %v\n", *res.Employee)
	}
}

func GetByBadgeNumber(client proto.EmployeeServiceClient, badgeNumber *int) {
	md := metadata.MD{}
	md["user"] = []string{"THISISUSER"}
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)
	res, err := client.GetByBadgeNumber(ctx, &proto.GetByBadgeNumberRequest{BadgeNumber: int32(*badgeNumber)})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Employee: %v\n", res.Employee)
}
