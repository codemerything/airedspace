package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler, jwtSecret string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware hit")
		fmt.Println("all cookies:", r.Cookies())
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unathorized: missing auth cookie", http.StatusUnauthorized)
			return
		}
		fmt.Println("cookie value:", cookie.Value)

		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		fmt.Println("parse error:", err)
		fmt.Println("token valid:", token.Valid)
		if err != nil || !token.Valid {
			http.Error(w, "Unathorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("parse error:", err)
		fmt.Println("token valid:", token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
			return
		}

		username := claims["username"].(string)
		userID := claims["user_id"].(float64)
		ctx := context.WithValue(r.Context(), "username", username)
		ctx = context.WithValue(ctx, "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
