package factory

import (
	"github.com/Shelex/grpc-go-demo/client/graph/model"
	"github.com/Shelex/grpc-go-demo/proto"
)

func EmployeeFromAPIToProto(e model.EmployeeInput) *proto.Employee {
	return &proto.Employee{
		BadgeNumber:         int32(e.BadgeNumber),
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		CountryCode:         e.CountryCode,
		VacationAccrualRate: float32(e.VacationAccrualRate),
		VacationAccrued:     float32(e.VacationAccrued),
	}
}

func EmployeeFromProtoToApi(e *proto.Employee) *model.Employee {
	apiDocs := make([]*string, len(e.Documents))
	for i := range apiDocs {
		apiDocs[i] = &e.Documents[i]
	}
	vacations := make([]*string, len(e.Vacations))
	for i := range vacations {
		vacations[i] = &e.Vacations[i]
	}
	return &model.Employee{
		ID:                  e.ID,
		BadgeNumber:         int(e.BadgeNumber),
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		CountryCode:         e.CountryCode,
		VacationAccrualRate: float64(e.VacationAccrualRate),
		VacationAccrued:     float64(e.VacationAccrued),
		Vacations:           vacations,
		Documents:           apiDocs,
	}
}

func VacationsFromProtoToApi(vacations []*proto.Vacation) []*model.Vacation {
	apiVacations := make([]*model.Vacation, 0, len(vacations))
	for _, v := range vacations {
		apiVacations = append(apiVacations, VacationFromProtoToApi(v))
	}
	return apiVacations
}

func VacationFromProtoToApi(v *proto.Vacation) *model.Vacation {
	return &model.Vacation{
		ID:            v.ID,
		UserID:        v.UserID,
		DurationHours: float64(v.DurationHours),
		StartDate:     int(v.StartDate),
		Cancelled:     v.Cancelled,
		Approved:      v.Approved,
	}
}

func AttachmentFromProtoToApi(a *proto.Attachment) *model.Document {
	stringifiedData := string(a.Data)
	return &model.Document{
		ID:        a.ID,
		UserID:    &a.UserID,
		FileName:  a.Filename,
		Data:      &stringifiedData,
		CreatedAt: int(a.CreatedAt),
	}
}
