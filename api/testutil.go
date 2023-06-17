package api

import (
	"github.com/jaswdr/faker"
	db "github.com/kerok-kristoffer/backendStub/db/sqlc"
	"github.com/kerok-kristoffer/backendStub/util"
)

func randomUserWithPassword() (db.User, string, error) {
	f := faker.New()
	password := f.Internet().Password()
	return db.User{
		ID:       0,
		UserName: f.Person().FirstName(),
		Email:    f.Internet().Email(),
		FullName: f.Person().Name(),
		Hash:     "",
	}, password, nil
}

func randomUser() (db.User, error) {
	f := faker.New()
	password := f.Internet().Password()
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return db.User{}, err
	}

	return db.User{
		ID:       f.Int64(),
		UserName: f.Person().FirstName(),
		Email:    f.Internet().Email(),
		FullName: f.Person().Name(),
		Hash:     hashedPassword,
	}, nil
}
