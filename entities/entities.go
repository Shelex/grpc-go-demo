package entities

import (
	"github.com/Shelex/grpc-go-demo/proto"
)

type Vacation struct {
	Id          int32
	StartDate   int64
	Duration    float32
	IsCancelled bool
}

type Employee struct {
	Id                  int32
	BadgeNumber         int32
	FirstName           string
	LastName            string
	VacationAccrualRate float32
	VacationAccrued     float32
	Vacations           []Vacation
}

func EmployeeFromStorageToProto(e Employee) *proto.Employee {
	return &proto.Employee{
		Id:                  e.Id,
		BadgeNumber:         e.BadgeNumber,
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		VacationAccrualRate: e.VacationAccrualRate,
		VacationAccrued:     e.VacationAccrued,
		Vacations:           VacationsFromStorageToProto(e.Vacations),
	}
}

func EmployeeFromProtoToStorage(e *proto.Employee) Employee {
	return Employee{
		Id:                  e.Id,
		BadgeNumber:         e.BadgeNumber,
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		VacationAccrualRate: e.VacationAccrualRate,
		VacationAccrued:     e.VacationAccrued,
	}
}

func VacationsFromStorageToProto(vacations []Vacation) []*proto.Vacation {
	protoVacations := make([]*proto.Vacation, len(vacations))
	for _, v := range vacations {
		protoVacations = append(protoVacations, &proto.Vacation{
			Id:          v.Id,
			Duration:    v.Duration,
			IsCancelled: v.IsCancelled,
			StartDate:   v.StartDate,
		})
	}
	return protoVacations

}
