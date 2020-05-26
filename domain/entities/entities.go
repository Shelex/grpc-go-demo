package entities

import (
	"github.com/Shelex/grpc-go-demo/proto"
)

type Vacation struct {
	ID            string
	StartDate     int64
	DurationHours float32
	Approved      bool
	Cancelled     bool
}

type Employee struct {
	ID                  string
	BadgeNumber         int32
	FirstName           string
	LastName            string
	CountryCode         string
	VacationAccrualRate float32
	VacationAccrued     float32
	Vacations           []Vacation
	Documents           []string
}

type Document struct {
	ID        string
	UserID    string
	FileName  string
	Data      []byte
	CreatedAt int64 // Unix timestamp
}

func DocumentFromStorageToProto(d Document) *proto.Attachment {
	return &proto.Attachment{
		ID:        d.ID,
		UserID:    d.UserID,
		Filename:  d.FileName,
		CreatedAt: d.CreatedAt,
		Data:      d.Data,
	}
}

func EmployeeFromStorageToProto(e Employee) *proto.Employee {
	return &proto.Employee{
		ID:                  e.ID,
		BadgeNumber:         e.BadgeNumber,
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		CountryCode:         e.CountryCode,
		VacationAccrualRate: e.VacationAccrualRate,
		VacationAccrued:     e.VacationAccrued,
		Vacations:           VacationsFromStorageToProto(e.Vacations),
		Documents:           e.Documents,
	}
}

func EmployeeFromProtoToStorage(e *proto.Employee) Employee {
	return Employee{
		ID:                  e.ID,
		BadgeNumber:         e.BadgeNumber,
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		CountryCode:         e.CountryCode,
		VacationAccrualRate: e.VacationAccrualRate,
		VacationAccrued:     e.VacationAccrued,
	}
}

func VacationsFromStorageToProto(vacations []Vacation) []*proto.Vacation {
	protoVacations := make([]*proto.Vacation, len(vacations))
	for _, v := range vacations {
		protoVacations = append(protoVacations, VacationFromStorageToProto(v))
	}
	return protoVacations
}

func VacationFromStorageToProto(v Vacation) *proto.Vacation {
	return &proto.Vacation{
		ID:            v.ID,
		StartDate:     v.StartDate,
		DurationHours: v.DurationHours,
		Approved:      v.Approved,
		Cancelled:     v.Cancelled,
	}
}

func VacationFromProtoToStorage(v *proto.Vacation) Vacation {
	return Vacation{
		ID:            v.ID,
		StartDate:     v.StartDate,
		DurationHours: v.DurationHours,
		Approved:      v.Approved,
		Cancelled:     v.Cancelled,
	}
}
