// Package model model User
package model

import "time"

// User model info
// @Description User account information
type User struct {
	ID       string    `json:"id" bson:"_id"`
	Login    string    `json,bson,query,header:"login" validate:"required,alphanum,gte=5,lte=20"`
	Email    string    `json,bson:"email" validate:"required,email" format:"email"`
	Password string    `json,bson,query,header:"password" validate:"required"`
	Name     string    `json,bson:"name" validate:"required,alpha,gte=2,lte=25"`
	Age      int       `json,bson:"age" validate:"required,gte=0,lte=100"`
	Token    string    `json,bson:"token"`
	Role     string    `json,bson:"role"`
	Created  time.Time `json,bson:"created" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	Updated  time.Time `json,bson:"updated" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
}
