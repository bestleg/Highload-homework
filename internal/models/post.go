package models

import "otus-homework/internal/validator"

type InputCreatePost struct {
	Text      string              `json:"text"`
	Validator validator.Validator `json:"-"`
}

type InputUpdatePost struct {
	PostID    string              `json:"id"`
	Text      string              `json:"text"`
	Validator validator.Validator `json:"-"`
}

type OutputPost struct {
	Text string `json:"text"`
}
