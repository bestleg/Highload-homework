package models

import "otus-homework/internal/validator"

type InputUserSearch struct {
	FirstName  string              `json:"first_name"`
	SecondName string              `json:"last_name"`
	Validator  validator.Validator `json:"-"`
}

type InputUser struct {
	Password   string              `json:"password"`
	FirstName  string              `json:"first_name"`
	SecondName string              `json:"second_name"`
	Birthdate  JsonBirthDate       `json:"birthdate"`
	Sex        bool                `json:"sex"`
	Biography  string              `json:"biography"`
	City       string              `json:"city"`
	Validator  validator.Validator `json:"-"`
}

type InputAuthToken struct {
	UserID    string              `json:"id"`
	Password  string              `json:"password"`
	Validator validator.Validator `json:"-"`
}

type OutputUser struct {
	FirstName  string        `json:"first_name"`
	SecondName string        `json:"second_name"`
	Birthdate  JsonBirthDate `json:"birthdate"`
	Sex        bool          `json:"sex"`
	Biography  string        `json:"biography"`
	City       string        `json:"city"`
	Age        int           `json:"age"`
}
