package middleware

import (
	"MYAPI/helper"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET")) // Change cette clé pour une valeur sécurisée

// AuthMiddleware : Vérifie le token, extrait l'utilisateur et vérifie son rôle.
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token") // Récupérer le token depuis le cookie
		if err != nil {
			helper.HandleError(c, 401, "Token manquant ou invalide", err)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // Parser et valider le token
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature invalide")
			}
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			helper.HandleError(c, 401, "Token invalide", err)
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
				helper.HandleError(c, 403, "Accès interdit", err)
				c.Abort()
				return
			}
		} else {
			helper.HandleError(c, 401, "Token invalide", err)
			c.Abort()
			return
		}
		c.Next()
	}
}

// CheckAuth : Vérifie si l'utilisateur est authentifié.
func CheckAuth(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"isAuthenticated": false})
		return
	}

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
