package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/trenchesdeveloper/go-ecom/config"
	"github.com/trenchesdeveloper/go-ecom/types"
	"github.com/trenchesdeveloper/go-ecom/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GenerateToken(secret string, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := config.Envs.JWTSecret
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("unexpected signing method: %v", token.Header["alg"])
				utils.ErrorJSON(w, fmt.Errorf("permission denied"), http.StatusUnauthorized)
				return nil, fmt.Errorf("permission denied")
			}

			return []byte(secret), nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			log.Println("invalid token")
			utils.ErrorJSON(w, fmt.Errorf("permission denied"), http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(claims["userID"].(string))

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := store.GetUserByID(userID)

		if err != nil {
			utils.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", user.ID))

		handlerFunc(w, r)
	}
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return -1
	}
	return userID
}
