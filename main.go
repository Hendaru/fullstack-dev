package main

import (
	"log"
	"net/http"
	"projectgolang/auth"
	"projectgolang/handler"
	"projectgolang/helper"
	"projectgolang/products"
	"projectgolang/users"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "hore:hore@tcp(host.docker.internal:3306)/bwstartup?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//REPOSITORY
	userRepository := users.NewUsersRepository(db)
	productsRepository := products.NewProductRepository(db)

	//SERVICE
	userService := users.NewUsersService(userRepository)
	authService := auth.NewAuthService()
	productsService := products.NewProductService(productsRepository)

	//HANDLER
	userHandler := handler.NewUsersHandler(userService, authService)
	productHandler := handler.NewProductsHandler(productsService, productsRepository)

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")
	//api version ex v1, v2 dll
	api := router.Group("/api/v1")

	// API
	//USER
	api.POST("/user", userHandler.RegisterUsersHandler)
	api.POST("/email_check", userHandler.CheckEmailAvailabilityService)
	api.POST("/login", userHandler.LoginUsersHandler)
	api.PUT("/user/:id", authMidleware(authService, userService), userHandler.UpdateUserHandler)
	api.POST("/avatar-user", authMidleware(authService, userService), userHandler.UploadAvatarHandler)

	//PRODUCT
	api.GET("/product", authMidleware(authService, userService), productHandler.GetProductByUserIDHandler)
	api.POST("/product", authMidleware(authService, userService), productHandler.CreateProductHandler)
	api.PUT("/product/:id", authMidleware(authService, userService), productHandler.UpdateProductHandler)
	api.PUT("/product-image/:id", authMidleware(authService, userService), productHandler.UpdateProductImageHandler)

	router.Run(":8080")
}

func authMidleware(authService auth.AuthService, userService users.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHandler := c.GetHeader("Authorization")

		if !strings.Contains(authHandler, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHandler, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByIDService(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)
	}
}
