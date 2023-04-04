package products

import "projectgolang/users"

type GetProductsIdInput struct {
	ID string `uri:"id" binding:"required"`
}

type CreatedProductsInput struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Price       int    `form:"price"`
	Stock       int    `form:"stock"`
	User        users.User
}
