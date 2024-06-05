package database

import (
	"database/sql"
	"fmt"
	"time"
)

// TestNew is the function used to create an instance of the database package. It returns the type
// Model, which embeds all the types we want to be available to our application. This
// is only used when running tests.
func TestNew(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: &UserTest{},
		Plan: &PlanTest{},
	}
}

// UserTest is the structure which holds one user from the database,
// and is used for testing.
type UserTest struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    bool
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

func (u *UserTest) GetAll() ([]*User, error) {
	var users []*User

	user := User{
		ID:        "00000000-0000-0000-0000-000000000001",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "abc1234",
		Active:    true,
		IsAdmin:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users = append(users, &user)

	return users, nil
}

func (u *UserTest) UserExists(email string) bool {
	return true
}

func (u *UserTest) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        "00000000-0000-0000-0000-000000000001",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "abc1234",
		Active:    true,
		IsAdmin:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Plan: &Plan{
			ID:                  "00000000-0000-0000-0000-000000000005",
			PlanName:            "SILVER",
			PlanAmount:          1000,
			PlanAmountFormatted: "$20.00",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	return &user, nil
}

func (u *UserTest) GetOne(id string) (*User, error) {
	user := User{
		ID:        "00000000-0000-0000-0000-000000000001",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "abc1234",
		Active:    true,
		IsAdmin:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Plan: &Plan{
			ID:                  "00000000-0000-0000-0000-000000000005",
			PlanName:            "SILVER",
			PlanAmount:          1000,
			PlanAmountFormatted: "$20.00",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	return &user, nil
}

func (u *UserTest) Update() error {
	return nil
}

func (u *UserTest) Delete() error {
	return nil
}

func (u *UserTest) DeleteByID(id string) error {
	return nil
}

func (u *UserTest) Insert() (string, error) {
	return "00000000-0000-0000-0000-000000000002", nil
}

func (u *UserTest) ResetPassword(password string) error {
	return nil
}

func (u *UserTest) PasswordMatches(email string, password string) (bool, error) {
	return true, nil
}

type PlanTest struct {
	ID                  string
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p *PlanTest) GetAll() ([]*Plan, error) {
	var plans []*Plan

	plan := Plan{
		ID:                  "00000000-0000-0000-0000-000000000005",
		PlanName:            "SILVER",
		PlanAmount:          1000,
		PlanAmountFormatted: "$20.00",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	plans = append(plans, &plan)

	return plans, nil
}

func (p *PlanTest) GetOne(id string) (*Plan, error) {
	plan := Plan{
		ID:                  "00000000-0000-0000-0000-000000000005",
		PlanName:            "SILVER",
		PlanAmount:          1000,
		PlanAmountFormatted: "$20.00",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return &plan, nil
}

func (p *PlanTest) SubscribeUserToPlan(userID string, planID string) error {
	return nil
}

func (p *PlanTest) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.0
	return fmt.Sprintf("$%.2f", amount)
}
