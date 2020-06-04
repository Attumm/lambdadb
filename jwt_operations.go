package main

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
)


func handleJWT(tokenString, sharedSecret string ) (map[string] interface{}, bool) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sharedSecret), nil
	})
	if err != nil {
		fmt.Println("error: ", token, err)
		return nil, false
	}
	return claims, true
}

