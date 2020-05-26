package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Shelex/grpc-go-demo/domain/entities"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Repository struct {
	session   *mgo.Session
	name      string // db name
	employees string // employees collection name
	err       error
}

func NewMongoStorage(url string) (Storage, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	mdb := Repository{
		session:   session,
		name:      "employee-service",
		employees: "employees",
		err:       nil,
	}

	return &mdb, nil
}

func (r *Repository) Stop(ctx context.Context) error {
	if r.err != nil {
		return r.err
	}
	log.Println("closing mongodb session...")
	r.session.Close()
	log.Println("mongo session is closed")
	return nil
}

func (r *Repository) GetEmployee(ID string) (entities.Employee, error) {
	session := r.session.Copy()
	defer session.Close()
	var e entities.Employee
	err := session.DB(r.name).C(r.employees).FindId(ID).One(&e)
	if err == mgo.ErrNotFound {
		return e, fmt.Errorf("not found employee with id: %s", ID)
	}
	return e, nil
}

func (r *Repository) GetAll() ([]entities.Employee, error) {
	session := r.session.Copy()
	defer session.Close()
	var emps []entities.Employee
	err := session.DB(r.name).C(r.employees).Find(nil).Sort("_id").All(&emps)
	return emps, err
}

func (r *Repository) AddEmployee(e entities.Employee) (entities.Employee, error) {
	var empty entities.Employee
	session := r.session.Copy()
	defer session.Close()
	if err := session.DB(r.name).C(r.employees).Insert(e); err != nil {
		return empty, err
	}
	return e, nil
}

func (r *Repository) Count() (int, error) {
	session := r.session.Copy()
	defer session.Close()
	return session.DB(r.name).C(r.employees).Count()
}

func (r *Repository) UpdateEmployee(ID string, e entities.Employee) (entities.Employee, error) {
	var empty entities.Employee
	session := r.session.Copy()
	defer session.Close()

	update := make(bson.M)
	log.Printf("incoming employee updates: %v", e)

	if e.BadgeNumber != 0 {
		update["badgeNumber"] = e.BadgeNumber
	}
	if e.FirstName != "" {
		update["firstName"] = e.FirstName
	}
	if e.LastName != "" {
		update["lastName"] = e.LastName
	}
	if e.CountryCode != "" {
		update["countryCode"] = e.CountryCode
	}
	if e.VacationAccrualRate != 0 {
		update["vacationAccrualRate"] = e.VacationAccrualRate
	}
	if e.VacationAccrued != 0 {
		update["vacationAccrued"] = e.VacationAccrued
	}

	log.Printf("result of employee updates: %v", update)

	if err := session.DB(r.name).C(r.employees).UpdateId(ID, bson.M{"$set": update}); err != nil {
		return empty, err
	}

	employee, err := r.GetEmployee(ID)

	log.Printf("now employee is: %v", employee)

	if err != nil {
		return empty, err
	}
	return employee, nil
}

func (r *Repository) DeleteEmployee(ID string) (entities.Employee, error) {
	var empty entities.Employee
	session := r.session.Copy()
	defer session.Close()
	e, err := r.GetEmployee(ID)
	if err != nil {
		return empty, err
	}
	if err := session.DB(r.name).C(r.employees).RemoveId(ID); err != nil {
		return empty, err
	}
	return e, nil
}

func (r *Repository) AddDocument(userID string, ID string) error {
	session := r.session.Copy()
	defer session.Close()
	employee, err := r.GetEmployee(userID)
	if err != nil {
		return fmt.Errorf("failed to find user with id: %s", userID)
	}

	documentIndex := "documents." + strconv.Itoa(len(employee.Documents))
	update := make(bson.M)
	update[documentIndex] = ID

	if err := session.DB(r.name).C(r.employees).UpdateId(userID, bson.M{"$set": update}); err != nil {
		return fmt.Errorf("failed to update documents for user with id %s: %s", userID, err)
	}

	return nil
}

func (r *Repository) AddVacation(ID string, userID string, startDate int64, durationHours float32) (entities.Vacation, error) {
	session := r.session.Copy()
	defer session.Close()
	var v entities.Vacation
	employee, err := r.GetEmployee(userID)
	if err != nil {
		return v, err
	}

	v.ID = ID
	v.StartDate = startDate
	v.DurationHours = durationHours

	vacationIndex := "vacations." + strconv.Itoa(len(employee.Vacations))
	update := make(bson.M)
	update[vacationIndex] = v

	if err := session.DB(r.name).C(r.employees).UpdateId(userID, bson.M{"$set": update}); err != nil {
		return v, err
	}

	return v, nil
}