package products

import "errors"

type ProductServiceInterface interface {
	GetProductByUserIDService(userID int) ([]Product, error)
	CreateProductService(input CreatedProductsInput, pathImage string, productId string) (Product, error)
	UpdateProductService(input CreatedProductsInput, getProductsIdInput GetProductsIdInput) (Product, error)
	UpdateProductImageService(userID int, pathImage string, getProductsIdInput GetProductsIdInput) (Product, error)
}

type productServiceInterface struct {
	productRepositoryInterface ProductRepositoryInterface
}

func NewProductService(productRepositoryInterface ProductRepositoryInterface) *productServiceInterface {
	return &productServiceInterface{productRepositoryInterface}
}

func (s *productServiceInterface) GetProductByUserIDService(userID int) ([]Product, error) {

	products, err := s.productRepositoryInterface.GetProductByUserIDRepository(userID)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (s *productServiceInterface) CreateProductService(input CreatedProductsInput, pathImage string, productId string) (Product, error) {
	product := Product{}
	product.ID = productId
	product.Title = input.Title
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.UserID = input.User.ID
	product.ProductImage = pathImage

	newProducts, err := s.productRepositoryInterface.CreateProductRepository(product)
	if err != nil {
		return newProducts, err
	}

	return newProducts, nil
}

func (s *productServiceInterface) UpdateProductService(input CreatedProductsInput, getProductsIdInput GetProductsIdInput) (Product, error) {
	product, err := s.productRepositoryInterface.FindProductByIdRepository(getProductsIdInput.ID)
	if err != nil {
		return product, err
	}
	if product.UserID != input.User.ID {
		return product, errors.New("You are not the owner of this product")
	}

	product.ID = getProductsIdInput.ID
	product.Title = input.Title
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.UserID = input.User.ID

	updateProduct, err := s.productRepositoryInterface.UpdateProductRepository(product)
	if err != nil {
		return updateProduct, err
	}

	return updateProduct, nil
}

func (s *productServiceInterface) UpdateProductImageService(userID int, pathImage string, getProductsIdInput GetProductsIdInput) (Product, error) {
	product, err := s.productRepositoryInterface.FindProductByIdRepository(getProductsIdInput.ID)
	if err != nil {
		return product, err
	}
	if product.UserID != userID {
		return product, errors.New("You are not the owner of this product")
	}

	product.ProductImage = pathImage

	updateProduct, err := s.productRepositoryInterface.UpdateProductRepository(product)
	if err != nil {
		return updateProduct, err
	}

	return updateProduct, nil
}
