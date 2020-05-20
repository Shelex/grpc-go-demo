package storage

import (
	"errors"
	"fmt"
	"github.com/Shelex/grpc-go-demo/entities"
)

type Storage interface {
	GetByBadge(badgeNum int32) (entities.Employee, error)
	GetAll() ([]entities.Employee, error)
	AddEmployee(entities.Employee) error
	Count() (int, error)
}

type InMem struct {
	employees []entities.Employee
}

func NewInMemStorage() Storage {
	return &InMem{
		employees: []entities.Employee{
			{
				Id:                  1,
				BadgeNumber:         7975,
				FirstName:           "John",
				LastName:            "Doe",
				VacationAccrualRate: 2,
				VacationAccrued:     30,
			},
			{
				Id:                  2,
				BadgeNumber:         7294,
				FirstName:           "Mark",
				LastName:            "Murphy",
				VacationAccrualRate: 2.3,
				VacationAccrued:     21.4,
			},
			{
				Id:                  3,
				BadgeNumber:         5193,
				FirstName:           "Donna",
				LastName:            "Cortez",
				VacationAccrualRate: 3,
				VacationAccrued:     23.2,
			},
			{
				Id:                  4,
				BadgeNumber:         8480,
				FirstName:           "Micheal",
				LastName:            "Wood",
				VacationAccrualRate: 3.4,
				VacationAccrued:     45.2,
			},
			{
				Id:                  5,
				BadgeNumber:         6238,
				FirstName:           "Louis",
				LastName:            "Alvarez",
				VacationAccrualRate: 0.485,
				VacationAccrued:     2.5,
			},
		},
	}
}

func (i *InMem) GetByBadge(badgeNum int32) (entities.Employee, error) {
	for _, e := range i.employees {
		if badgeNum == e.BadgeNumber {
			return e, nil
		}
	}
	return entities.Employee{}, errors.New("employee not found")
}

func (i *InMem) GetAll() ([]entities.Employee, error) {
	return i.employees, nil
}

func (i *InMem) AddEmployee(e entities.Employee) error {
	for _, employee := range i.employees {
		if e.BadgeNumber == employee.BadgeNumber {
			return fmt.Errorf("badge number %d is duplicated", e.BadgeNumber)
		}
	}
	i.employees = append(i.employees, e)
	return nil
}

func (i *InMem) Count() (int, error) {
	return len(i.employees), nil
}
