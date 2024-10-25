package models

import (
	"errors"
	"time"

	"docflow-backend/db"
	"docflow-backend/utils"
)

type User struct {
	ID          int64
	FirstName   string
	LastName    string
	DateOfBirth time.Time
	Email       string `binding:"required"`
	Password    string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users (firstName, lastName, dateOfBirth, email, password) VALUES (?, ?, ?, ?, ?);"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(u.FirstName, u.LastName, u.DateOfBirth, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?;"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}

func GeUserByID(id int64) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
