package products

import (
	"projectgolang/users"
	"time"
)

type Product struct {
	ID           string
	Title        string
	UserID       int
	Description  string
	Price        int
	Stock        int
	ProductImage string
	User         users.User
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
