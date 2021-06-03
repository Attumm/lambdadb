package main

import (
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
)

type Claims struct {
	jwt.Payload
	//	Id             string            `json:"username"`
	//	Name           string            `json:"name"`
	//	LegacyId       int               `json:"legacy_id"`
	//	Exp            int               `json:"exp"`
	Groups []string `json:"groups"`
	//	AdditionalInfo map[string]string `json:"additional_info"`
	//	DName          string            `json:"d_name"`
}

func handleJWT(tokenString, sharedSecret string) (Claims, error) {
	var claims Claims
	var hs = jwt.NewHS256([]byte(sharedSecret))
	hd, err := jwt.Verify([]byte(tokenString), hs, &claims)
	if err != nil {
		return claims, fmt.Errorf("can't verify jwt %v %v %v", hd, hs, claims)
	}
	return claims, nil
}
