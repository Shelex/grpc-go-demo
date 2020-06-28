// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Document struct {
	ID        string  `json:"id"`
	UserID    *string `json:"userID"`
	FileName  string  `json:"fileName"`
	Data      *string `json:"data"`
	CreatedAt int     `json:"createdAt"`
}

type Employee struct {
	ID                  string    `json:"id"`
	BadgeNumber         int       `json:"badgeNumber"`
	FirstName           string    `json:"firstName"`
	LastName            string    `json:"lastName"`
	CountryCode         string    `json:"countryCode"`
	VacationAccrualRate float64   `json:"vacationAccrualRate"`
	VacationAccrued     float64   `json:"vacationAccrued"`
	Vacations           []*string `json:"vacations"`
	Documents           []*string `json:"documents"`
}

type EmployeeInput struct {
	BadgeNumber         int     `json:"badgeNumber"`
	FirstName           string  `json:"firstName"`
	LastName            string  `json:"lastName"`
	CountryCode         string  `json:"countryCode"`
	VacationAccrualRate float64 `json:"vacationAccrualRate"`
	VacationAccrued     float64 `json:"vacationAccrued"`
}

type EmployeeSaveResult struct {
	SavedEmployees []*Employee `json:"savedEmployees"`
	Errors         *string     `json:"errors"`
}

type Vacation struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userId"`
	StartDate     int     `json:"startDate"`
	DurationHours float64 `json:"durationHours"`
	Approved      bool    `json:"approved"`
	Cancelled     bool    `json:"cancelled"`
}

type VacationRequest struct {
	UserID        string  `json:"userID"`
	StartDate     int     `json:"startDate"`
	DurationHours float64 `json:"durationHours"`
}
