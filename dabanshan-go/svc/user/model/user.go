package model

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	ErrNoCustomerInResponse = errors.New("Response has no matching customer")
	ErrMissingField         = "Error missing %v"
)


type User struct {
	FirstName string
	LastName string
	Email string
	UserName string
	Password string
	UserID string
	Salt string
	Authority string
}

// New ..
func New() User {
	u := User{}
	u.NewSalt()
	return u
}

// NewSalt ..
func (u *User) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}

// Validate ..
func (u *User) Validate() error {
	if u.FirstName == "" {
		return fmt.Errorf(ErrMissingField, "FirstName")
	}
	if u.LastName == "" {
		return fmt.Errorf(ErrMissingField, "LastName")
	}
	if u.UserName == "" {
		return fmt.Errorf(ErrMissingField, "UserName")
	}
	if u.Password == "" {
		return fmt.Errorf(ErrMissingField, "Password")
	}
	return nil
}
type Falter interface {
	Failed() error
}

// GetUserRequest collects the request parameters for the GetProducts method.
type GetUserRequest struct {
	A string
}

// GetUserResponse collects the response values for the GetProducts method.
type GetUserResponse struct {
	V   User  `json:"v"`
	Err error `json:"err,omitempty"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r GetUserResponse) Failed() error {
	return r.Err
}

// RegisterRequest struct
type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// RegisterUserResponse ...
type RegisterUserResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// LoginRequest ..
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse ..
type LoginResponse struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}