package middleware

import (
	"MYAPI/helper"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET")) // Change cette clé pour une valeur sécurisée

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("token") // Récupérer le token depuis le cookie
		if err != nil {
			helper.HandleError(c, 401, "Token manquant ou invalide",err)
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
			helper.HandleError(c, 401, "Token invalide",err)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok { // Extraire les claims (données) du token
			userID := claims["id"].(string)
			userEmail := claims["email"].(string)
			role := claims["role"].(string)

			c.Set("userID", userID)
			c.Set("userEmail", userEmail)
			c.Set("role", role)

			// Vérifier si le rôle correspond à celui requis
			if requiredRole != "" && requiredRole != role {
				helper.HandleError(c, 403, "Accès interdit",err)
				c.Abort()
				return
			}
		} else {
			helper.HandleError(c, 401, "Token invalide",err)
			c.Abort()
			return
		}
		c.Next()
	}
}
