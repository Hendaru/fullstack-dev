package products

type ProductFormatter struct {
	ID           string `json:"id"`
	UserID       int    `json:"user_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ProductImage string `json:"product_image"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
}

func FormatProduct(product Product) ProductFormatter {
	formatter := ProductFormatter{}
	formatter.ID = product.ID
	formatter.UserID = product.UserID
	formatter.Title = product.Title
	formatter.Description = product.Description
	formatter.Price = product.Price
	formatter.Stock = product.Stock

	formatter.ProductImage = product.ProductImage
	return formatter
}

func FormatProducts(product []Product) []ProductFormatter {

	productFormatter := []ProductFormatter{}
	for _, value := range product {

		productFormatter = append(productFormatter, FormatProduct(value))
	}

	return productFormatter
}
