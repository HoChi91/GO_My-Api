package routes

import (
	"MYAPI/controllers"
	"MYAPI/middleware"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	router := gin.Default()

	// Configuration CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Origine spécifique
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Important pour les cookies et `credentials: 'include'`
	}))

	router.GET("", func(c *gin.Context) {
		c.JSON(200, "HOME")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	apiGroup := router.Group("/api")
	AuthGroup := apiGroup.Group("/auth")              // Routes pour authentification
	AuthGroup.POST("/login", controllers.Login)       // Connexion (pas sécurisé)
	AuthGroup.POST("/signup", controllers.CreateUser) // Inscription (pas sécurisé)

	UsersGroup := apiGroup.Group("/users") // Routes sécurisées pour les utilisateurs
	{
		UsersGroup.Use(middleware.AuthMiddleware(""))   // Accessible à tous les utilisateurs authentifiés
		UsersGroup.GET("/:id", controllers.GetUserByID) // Récupération d'un utilisateur spécifique (si non admin renvoie uniquement les details de l'user connecté)

		adminUsersGroup := UsersGroup.Group("") // Routes accessibles uniquement aux admins
		adminUsersGroup.Use(middleware.AuthMiddleware("admin"))
		adminUsersGroup.GET("", controllers.GetUsers)          // Récupération des utilisateurs
		adminUsersGroup.PUT("/:id", controllers.ModifyUser)    // Modification d'un utilisateur
		adminUsersGroup.DELETE("/:id", controllers.DeleteUser) // Suppression d'un utilisateur
	}
	AuthorsGroup := apiGroup.Group("/authors") //Routes pour les auteurs
	{
		AuthorsGroup.Use(middleware.AuthMiddleware("")) // Accessible à tous les utilisateurs authentifiés

		AuthorsGroup.GET("", controllers.GetAuthors)        // Route GET pour recuperer la liste des auteurs
		AuthorsGroup.GET("/:id", controllers.GetAuthorByID) // Route GET pour recuperer un auteur spécifique

		adminAuthorsGroup := AuthorsGroup.Group("")
		adminAuthorsGroup.Use(middleware.AuthMiddleware("admin")) // Accessible uniquement aux admins

		adminAuthorsGroup.POST("", controllers.CreateAuthor)       // Route POST pour ajouter un auteur
		adminAuthorsGroup.PUT("/:id", controllers.ModifyAuthor)    // Route PUT pour modifier un auteur spécifique
		adminAuthorsGroup.DELETE("/:id", controllers.DeleteAuthor) // Route DELETE pour supprimer un auteur spécifique
	}
	BooksGroup := apiGroup.Group("/books") // Routes pour les livres
	{
		BooksGroup.Use(middleware.AuthMiddleware("")) // Accessible à tous les utilisateurs authentifiés

		BooksGroup.GET("", controllers.GetBooks)        // Route GET pour recuperer la liste des livres
		BooksGroup.GET("/:id", controllers.GetBookByID) // Route GET pour recuperer un livre spécifique

		adminBooksGroup := BooksGroup.Group("")
		adminBooksGroup.Use(middleware.AuthMiddleware("admin")) // Accessible à tous les admins

		adminBooksGroup.POST("", controllers.CreateBook)       // Route POST pour ajouter un livre
		adminBooksGroup.PUT("/:id", controllers.ModifyBook)    // Route PUT pour modifier un livre spécifique
		adminBooksGroup.DELETE("/:id", controllers.DeleteBook) // Route DELETE pour supprimer un livre spécifique
	}
	PurchaseGroup := apiGroup.Group("/purchase") // Routes pour les achats
	{
		PurchaseGroup.Use(middleware.AuthMiddleware("")) // Accessible à tous les utilisateurs authentifiés
		PurchaseGroup.POST("/order", controllers.PurchaseBook)
		PurchaseGroup.POST("/cart", controllers.AddToCart)
		PurchaseGroup.DELETE("/cart", controllers.RemoveCart)
		PurchaseGroup.GET("/cart", controllers.GetCart)
		PurchaseGroup.POST("/cart/finalize", controllers.FinalizeCart)
		PurchaseGroup.GET("/history/:id", controllers.GetPurchaseHistoryByID)

		adminPurchaseGroup := PurchaseGroup.Group("")
		adminPurchaseGroup.Use(middleware.AuthMiddleware("admin")) // Accessible à tous les admins

		adminPurchaseGroup.GET("/history", controllers.GetAllPurchaseHistory)
	}
	return router
}
