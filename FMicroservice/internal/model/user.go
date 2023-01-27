package model

import "time"

type User struct {
	ID       string    `json:"id" bson:"_id"`
	Login    string    `json,bson:"login" validate:"required,alphanum,gte=5,lte=20"`
	Email    string    `json,bson:"email" validate:"required,email"`
	Password string    `json,bson:"password" validate:"required"`
	Name     string    `json,bson:"name" validate:"required,alpha,gte=2,lte=25"`
	Age      int       `json,bson:"age" validate:"required,gte=0,lte=100"`
	Token    string    `json,bson:"token"`
	Deleted  bool      `json,bson:"deleted"`
	Created  time.Time `json,bson:"created"`
	Updated  time.Time `json,bson:"updated"`
}