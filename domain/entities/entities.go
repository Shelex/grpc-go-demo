package entities

import (
	"github.com/Shelex/grpc-go-demo/proto"
)

type Vacation struct {
	ID            string  `json:"id" bson:"id"`
	UserID        string  `json:"userId" bson:"userId"`
	StartDate     int64   `json:"startDate" bson:"startDate"`
	DurationHours float32 `json:"durationHours" bson:"durationHours"`
	Approved      bool    `json:"approved" bson:"approved"`
	Cancelled     bool    `json:"cancelled" bson:"cancelled"`
}

type Employee struct {
	ID                  string   `json:"id" bson:"_id"`
	BadgeNumber         int32    `json:"badgeNumber" bson:"badgeNumber"`
	FirstName           string   `json:"firstName" bson:"firstName"`
	LastName            string   `json:"lastName" bson:"lastName"`
	CountryCode         string   `json:"countryCode" bson:"countryCode"`
	VacationAccrualRate float32  `json:"vacationAccrualRate" bson:"vacationAccrualRate"`
	VacationAccrued     float32  `json:"vacationAccrued" bson:"vacationAccrued"`
	Vacations           []string `json:"vacations" bson:"vacations"`
	Documents           []string `json:"documents" bson:"documents"`
}

type Document struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	FileName  string `json:"fileName"`
	Data      []byte
	CreatedAt int64 `json:"createdAt"` // Unix timestamp
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
		Vacations:           e.Vacations,
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
	protoVacations := make([]*proto.Vacation, 0, len(vacations))
	for i, v := range vacations {
		protoVacations[i] = VacationFromStorageToProto(v)
	}
	return protoVacations
}

func VacationFromStorageToProto(v Vacation) *proto.Vacation {
	return &proto.Vacation{
		ID:            v.ID,
		UserID:        v.UserID,
		StartDate:     v.StartDate,
		DurationHours: v.DurationHours,
		Approved:      v.Approved,
		Cancelled:     v.Cancelled,
	}
}

func VacationFromProtoToStorage(v *proto.Vacation) Vacation {
	return Vacation{
		ID:            v.ID,
		UserID:        v.UserID,
		StartDate:     v.StartDate,
		DurationHours: v.DurationHours,
		Approved:      v.Approved,
		Cancelled:     v.Cancelled,
	}
}
