package database

// UserInterface is the interface for the user type. In order
// to satisfy this interface, all specified methods must be implemented.
// We do this so we can test things easily. Both data.User and data.UserTest
// implement this interface.
type UserInterface interface {
	GetAll() ([]*User, error)
	UserExists(email string) bool
	GetByEmail(email string) (*User, error)
	GetOne(id string) (*User, error)
	Update(user *User) error
	Delete(user *User) error
	DeleteByID(id string) error
	Insert(user *User) (string, error)
	ResetPassword(email string, password string) error
	PasswordMatches(email string, password string) (bool, error)
}

// PlanInterface is the type for the plan type. Both data.Plan and data.PlanTest
// implement this interface.
type PlanInterface interface {
	GetAll() ([]*Plan, error)
	GetOne(id string) (*Plan, error)
	SubscribeUserToPlan(userID string, planID string) error
	AmountForDisplay() string
}
