package products

import "gorm.io/gorm"

type ProductRepositoryInterface interface {
	CreateProductRepository(product Product) (Product, error)
	GetProductByUserIDRepository(ID int) ([]Product, error)
	UpdateProductRepository(product Product) (Product, error)
	FindProductByIdRepository(ID string) (Product, error)
}

type productRepositoryInterface struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepositoryInterface {
	return &productRepositoryInterface{db}
}

func (r *productRepositoryInterface) GetProductByUserIDRepository(ID int) ([]Product, error) {
	var product []Product

	err := r.db.Where("user_id = ?", ID).Find(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *productRepositoryInterface) FindProductByIdRepository(ID string) (Product, error) {
	var product Product
	err := r.db.Where("id = ?", ID).Find(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepositoryInterface) CreateProductRepository(product Product) (Product, error) {
	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *productRepositoryInterface) UpdateProductRepository(product Product) (Product, error) {
	err := r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}
