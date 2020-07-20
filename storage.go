package main

import (
	"context"

	"github.com/volatiletech/authboss/v3"
)

// Check that current user and storage implement right authboss interfaces.
var (
	AssertUser   = &User{}
	AssertStorer = &MemStorage{}

	_ authboss.User         = AssertUser
	_ authboss.AuthableUser = AssertUser

	_ authboss.ServerStorer = AssertStorer
)

// User is a struct which implements User interface of Authboss module.
type User struct {
	email    string
	password string
}

func (u *User) GetPID() string {
	return u.email
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) PutPID(pid string) {
	u.email = pid
}

func (u *User) PutPassword(password string) {
	u.password = password
}

// MemStorage is in-memory database to store user's data and tokens.
type MemStorage struct {
	users  map[string]User
	tokens map[string][]string
}

// NewMemStorage creates new instance of MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		users: map[string]User{
			"joey@jojo.com": {
				email:    "joey@jojo.com",
				password: "joey", // Please, never store passwords like this, It's maybe the worst thing I ever done in my life.
			},
			"average@john.com": {
				email:    "average@john.com",
				password: "12345", // Please, never store passwords like this, It's maybe the worst thing I ever done in my life.
			},
		},
	}
}

func (ms *MemStorage) Load(_ context.Context, key string) (authboss.User, error) {
	v, ok := ms.users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &v, nil
}

func (ms *MemStorage) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)
	ms.users[u.email] = *u

	return nil
}
