package domain

type (
	Entity struct {
		ID      string `json:"id" bson:"_id"`
		Name    string `json,bson:"name" validate:"required"`
		Age     int    `json,bson:"age" validate:"gte=0,lte=100"`
		Deleted bool   `json,bson:"deleted"`
	}
)
