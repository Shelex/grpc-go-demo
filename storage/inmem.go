package storage

import (
	"fmt"
	"log"

	"github.com/Shelex/grpc-go-demo/domain/entities"
)

type Storage interface {
	GetEmployee(ID string) (entities.Employee, error)
	GetAll() ([]entities.Employee, error)
	AddEmployee(employee entities.Employee) (entities.Employee, error)
	Count() (int, error)
	UpdateEmployee(ID string, e entities.Employee) (entities.Employee, error)
	DeleteEmployee(ID string) (entities.Employee, error)
	AddDocument(userID string, ID string) error
	AddVacation(ID string, userID string, startDate int64, durationHours float32) (entities.Vacation, error)
}

type InMem struct {
	employees map[string]entities.Employee
}

func NewInMemStorage() (Storage, error) {
	return &InMem{
		employees: map[string]entities.Employee{
			"1": {
				ID:                  "1",
				BadgeNumber:         7975,
				FirstName:           "John",
				LastName:            "Doe",
				VacationAccrualRate: 2,
				VacationAccrued:     30,
			},
			"2": {
				ID:                  "2",
				BadgeNumber:         7294,
				FirstName:           "Mark",
				LastName:            "Murphy",
				VacationAccrualRate: 2.3,
				VacationAccrued:     21.4,
			},
			"3": {
				ID:                  "3",
				BadgeNumber:         5193,
				FirstName:           "Donna",
				LastName:            "Cortez",
				VacationAccrualRate: 3,
				VacationAccrued:     23.2,
			},
			"4": {
				ID:                  "4",
				BadgeNumber:         8480,
				FirstName:           "Micheal",
				LastName:            "Wood",
				VacationAccrualRate: 3.4,
				VacationAccrued:     45.2,
			},
			"5": {
				ID:                  "5",
				BadgeNumber:         6238,
				FirstName:           "Louis",
				LastName:            "Alvarez",
				VacationAccrualRate: 0.485,
				VacationAccrued:     2.5,
			},
		},
	}, nil
}

func (i *InMem) GetEmployee(ID string) (entities.Employee, error) {
	e, ok := i.employees[ID]
	if !ok {
		return e, fmt.Errorf("employee with id %s not found", ID)
	}
	return e, nil
}

func (i *InMem) GetAll() ([]entities.Employee, error) {
	count, err := i.Count()
	if err != nil {
		return nil, err
	}
	employees := make([]entities.Employee, 0, count)
	for _, e := range i.employees {
		employees = append(employees, e)
	}
	return employees, nil
}

func (i *InMem) AddEmployee(ID string, e entities.Employee) (entities.Employee, error) {
	var empty entities.Employee
	for _, employee := range i.employees {
		if e.BadgeNumber == employee.BadgeNumber {
			return empty, fmt.Errorf("badge number %d is duplicated", e.BadgeNumber)
		}
	}

	log.Printf("saving user id %s", ID)
	e.ID = ID

	i.employees[ID] = ID
	return e, nil
}

func (i *InMem) UpdateEmployee(ID string, e entities.Employee) (entities.Employee, error) {
	var empty entities.Employee
	emp, ok := i.employees[ID]
	if !ok {
		return empty, fmt.Errorf("employee with id %s not found", ID)
	}
	if e.BadgeNumber != 0 {
		emp.BadgeNumber = e.BadgeNumber
	}
	if e.FirstName != "" {
		emp.FirstName = e.FirstName
	}
	if e.LastName != "" {
		emp.LastName = e.LastName
	}
	if e.VacationAccrualRate != 0 {
		emp.VacationAccrualRate = e.VacationAccrualRate
	}
	if e.VacationAccrued != 0 {
		emp.VacationAccrualRate = e.VacationAccrualRate
	}
	i.employees[ID] = emp
	return emp, nil
}

func (i *InMem) AddDocument(userID string, ID string) error {
	emp, ok := i.employees[userID]
	if !ok {
		return fmt.Errorf("employee with id %s not found", userID)
	}
	emp.Documents = append(emp.Documents, ID)
	i.employees[userID] = emp
	return nil
}

func (i *InMem) DeleteEmployee(userID string) (entities.Employee, error) {
	var empty entities.Employee
	emp, ok := i.employees[userID]
	if !ok {
		return empty, fmt.Errorf("employee with id %s not found", userID)
	}
	delete(i.employees, userID)
	return emp, nil
}

func (i *InMem) Count() (int, error) {
	return len(i.employees), nil
}

func (i *InMem) AddVacation(ID string, userID string, startDate int64, durationHours float32) (entities.Vacation, error) {
	var v entities.Vacation
	employee, ok := i.employees[userID]
	if !ok {
		return v, fmt.Errorf("employee with id %s not found", userID)
	}

	v.ID = ID
	v.StartDate = startDate
	v.DurationHours = durationHours
	employee.Vacations = append(employee.Vacations, v)
	i.employees[userID] = employee
	return v, nil
}
