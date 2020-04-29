package main

import "github.com/Shelex/grpc-go-demo/proto"

var employees = []*proto.Employee{
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
}
