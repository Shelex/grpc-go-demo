package factory

import (
	"github.com/Shelex/grpc-go-demo/client/graph/model"
	"github.com/Shelex/grpc-go-demo/proto"
)

func EmployeeFromAPIToProto(e model.EmployeeInput) *proto.Employee {
	return &proto.Employee{
		Id:                  int32(e.ID),
		BadgeNumber:         int32(e.BadgeNumber),
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		VacationAccrualRate: float32(e.VacationAccrualRate),
	}
}

func EmployeeFromProtoToApi(e *proto.Employee) *model.Employee {
	return &model.Employee{
		ID:                  int(e.Id),
		BadgeNumber:         int(e.BadgeNumber),
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		VacationAccrualRate: float64(e.VacationAccrualRate),
		VacationAccrued:     float64(e.VacationAccrued),
		Vacations:           VacationsFromProtoToApi(e.Vacations),
	}
}

func VacationsFromProtoToApi(vacations []*proto.Vacation) []*model.Vacation {
	apiVacations := make([]*model.Vacation, len(vacations))
	for _, v := range vacations {
		apiVacations = append(apiVacations, &model.Vacation{
			ID:          int(v.Id),
			Duration:    float64(v.Duration),
			StartDate:   int(v.StartDate),
			IsCancelled: v.IsCancelled,
		})
	}
	return apiVacations
}
