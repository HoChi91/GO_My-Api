package middleware

import (
	"MYAPI/helper"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET")) // Assurez-vous de définir JWT_SECRET dans vos variables d'environnement

// AuthMiddleware : Vérifie le token JWT, extrait l'utilisateur et vérifie son rôle.
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le token depuis le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helper.HandleError(c, http.StatusUnauthorized, "Token manquant ou invalide", nil)
			c.Abort()
			return
		}

		// Extraire le token après "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parser et valider le token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature invalide")
			}
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			helper.HandleError(c, http.StatusUnauthorized, "Token invalide", err)
			c.Abort()
			return
		}

		// Extraire les claims (données) du token
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userEmail := claims["email"].(string)
			userID := claims["id"].(string)
			role := claims["role"].(string)

			// Stocker les informations dans le contexte de la requête
			c.Set("userID", userID)
			c.Set("userEmail", userEmail)
			c.Set("role", role)

			// Vérifier si le rôle correspond à celui requis
			if requiredRole != "" && requiredRole != role {
				helper.HandleError(c, http.StatusForbidden, "Accès interdit", nil)
				c.Abort()
				return
			}
		} else {
			helper.HandleError(c, http.StatusUnauthorized, "Token invalide", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// CheckAuth : Vérifie si l'utilisateur est authentifié.
func CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"isAuthenticated": false})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature invalide")
		}
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"isAuthenticated": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isAuthenticated": true})
}
