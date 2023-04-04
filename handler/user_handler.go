package handler

import (
	"fmt"
	"net/http"
	"projectgolang/auth"
	"projectgolang/helper"
	"projectgolang/users"
	"strconv"

	"github.com/gin-gonic/gin"
)

type usersHandler struct {
	usersService users.UsersService
	authService  auth.AuthService
}

func NewUsersHandler(usersService users.UsersService, authService auth.AuthService) *usersHandler {
	return &usersHandler{usersService, authService}
}

func (h *usersHandler) RegisterUsersHandler(c *gin.Context) {
	var input users.RegisterUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.usersService.RegisterUserService(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formater := users.FormatUser(newUser, token)
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)

}

func (h *usersHandler) CheckEmailAvailabilityService(c *gin.Context) {
	var input users.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.usersService.CheckEmailAvailabilityService(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *usersHandler) LoginUsersHandler(c *gin.Context) {
	var input users.LoginInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.usersService.LoginUserService(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formater := users.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Successfuly logged in", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)
}

func (h *usersHandler) UpdateUserHandler(c *gin.Context) {
	var inputID users.GetIDUser
	err := c.ShouldBindUri(&inputID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Update User Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	if currentUser.ID != inputID.ID {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var inputData users.UpdateUserInput
	err = c.ShouldBindJSON(&inputData)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Update User Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedUser, err := h.usersService.UpdateUserByIDService(inputID.ID, inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Update User Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formater := users.FormatUpdateUser(updatedUser)
	response := helper.APIResponse("Account has been updated", http.StatusOK, "success", formater)
	c.JSON(http.StatusOK, response)
}

func (h *usersHandler) UploadAvatarHandler(c *gin.Context) {

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload Avatar Failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	userID := currentUser.ID

	path := fmt.Sprintf("avatar/%d-%s", strconv.Itoa(userID), file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload Avatar Failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.usersService.UpdateAvatarByIDService(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload Avatar Failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar Successfuly Uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
