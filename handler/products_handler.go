package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"projectgolang/helper"
	"projectgolang/products"
	"projectgolang/users"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type productsHandler struct {
	service                    products.ProductServiceInterface
	productRepositoryInterface products.ProductRepositoryInterface
}

func NewProductsHandler(service products.ProductServiceInterface, productRepositoryInterface products.ProductRepositoryInterface) *productsHandler {
	return &productsHandler{service, productRepositoryInterface}
}

func (h *productsHandler) GetProductByUserIDHandler(c *gin.Context) {
	// userID, _ := strconv.Atoi(c.Param("user_id"))

	currentUser := c.MustGet("currentUser").(users.User)

	product, err := h.service.GetProductByUserIDService(currentUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": "Product not found"}
		response := helper.APIResponse("Get product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Get products success", http.StatusOK, "success", products.FormatProducts(product))
	c.JSON(http.StatusOK, response)
}

func (h *productsHandler) CreateProductHandler(c *gin.Context) {
	var input products.CreatedProductsInput
	err := c.ShouldBind(&input)
	if err != nil {
		error := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": error}
		response := helper.APIResponse("Create product failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser

	file, err := c.FormFile("product_image")
	if err != nil {
		errorMessage := gin.H{"errors": "Image is required"}
		response := helper.APIResponse("Create product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	productId := "product-" + ksuid.New().String()
	pathImageProduct := fmt.Sprintf("images/"+productId+"-%s", file.Filename)

	fileExtensionImage := filepath.Ext(pathImageProduct)

	if fileExtensionImage != ".jpg" && fileExtensionImage != ".png" {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err = c.SaveUploadedFile(file, pathImageProduct)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newProduct, err := h.service.CreateProductService(input, pathImageProduct, productId)
	if err != nil {
		response := helper.APIResponse("Create product failed 2", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Create product success", http.StatusOK, "success", products.FormatProduct(newProduct))
	c.JSON(http.StatusOK, response)
}

func (h *productsHandler) UpdateProductHandler(c *gin.Context) {
	var input products.CreatedProductsInput
	err := c.ShouldBind(&input)
	if err != nil {
		error := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": error}
		response := helper.APIResponse("Update product failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser

	var getProductsIdInput products.GetProductsIdInput
	err = c.ShouldBindUri(&getProductsIdInput)
	if err != nil {
		errorMessage := gin.H{"errors": "ID not found"}
		response := helper.APIResponse("Update product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updatedProduct, err := h.service.UpdateProductService(input, getProductsIdInput)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Update product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Update product success", http.StatusOK, "success", products.FormatProduct(updatedProduct))
	c.JSON(http.StatusOK, response)
}

func (h *productsHandler) UpdateProductImageHandler(c *gin.Context) {

	var getProductsIdInput products.GetProductsIdInput
	err := c.ShouldBindUri(&getProductsIdInput)
	if err != nil {
		errorMessage := gin.H{"errors": "ID not found"}
		response := helper.APIResponse("Update product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("product_image")
	if err != nil {
		errorMessage := gin.H{"errors": "Image is required"}
		response := helper.APIResponse("Create product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	productOld, err := h.productRepositoryInterface.FindProductByIdRepository(getProductsIdInput.ID)

	pathImageProduct := fmt.Sprintf("images/"+productOld.ID+"-%s", file.Filename)
	fileExtensionImage := filepath.Ext(pathImageProduct)

	if fileExtensionImage != ".jpg" && fileExtensionImage != ".png" {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = os.Remove(productOld.ProductImage)
	if err != nil {
		errorMessage := gin.H{"errors": "Failed Upload Image hhh"}
		response := helper.APIResponse("Failed Upload Image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.SaveUploadedFile(file, pathImageProduct)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)

	updatedProduct, err := h.service.UpdateProductImageService(currentUser.ID, pathImageProduct, getProductsIdInput)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Update product failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Update product success", http.StatusOK, "success", products.FormatProduct(updatedProduct))
	c.JSON(http.StatusOK, response)
}
