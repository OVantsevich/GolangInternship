package domain

type (
	Entity struct {
		ID        int    `sql:"primary_key"`
		Name      string `json:"name" sql:"type:varchar(50)"`
		Age       int    `json:"age" sql:"type:integer"`
		IsDeleted bool   `sql:"type:boolean"`
	}
)
