package dto

type Customer struct {
	CustomerId     string
	PhoneNumber    string
	FullName       string
	Salt           string
	Password       string
	Iteration      int
	SecurityLength int
}
