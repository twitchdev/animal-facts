/**
 *    Copyright 2018 Amazon.com, Inc. or its affiliates
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// ContextKeyType ...
type ContextKeyType string

// ChannelIDKey is the key that stores a request's Channel ID in Context
const ChannelIDKey ContextKeyType = "channelID"

// JWTClaims is the payload of a JWT
type JWTClaims struct {
	OpaqueUserID string               `json:"opaque_user_id,omitempty"`
	UserID       string               `json:"user_id"`
	ChannelID    string               `json:"channel_id,omitempty"`
	Role         string               `json:"role"`
	Permissions  JWTPubSubPermissions `json:"pubsub_perms"`
	jwt.StandardClaims
}

// JWTPubSubPermissions are PubSub permissions in JWTClaims
type JWTPubSubPermissions struct {
	Send   []string `json:"send,omitempty"`
	Listen []string `json:"listen,omitempty"`
}

// VerifyJWT is middleware to confirm the validity of incoming requests
func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			log.Println("JWT missing in request header")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			key, err := base64.StdEncoding.DecodeString(os.Getenv("EXT_SECRET"))

			if err != nil {
				return nil, err
			}

			return key, nil
		})

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if claims, ok := parsedToken.Claims.(*JWTClaims); ok && parsedToken.Valid {
			ctx := context.WithValue(r.Context(), ChannelIDKey, claims.ChannelID)

			if claims.Role != "broadcaster" {
				log.Println("Invalid JWT role")
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Println("Could not parse JWT")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	})
}

// NewJWT creates an EBS-signed JWT
func NewJWT(channelID string) string {
	var expiration = time.Now().Add(time.Minute*3).UnixNano() / int64(time.Millisecond)

	claims := JWTClaims{
		UserID:    os.Getenv("EXT_OWNER_ID"),
		ChannelID: channelID,
		Role:      "external",
		Permissions: JWTPubSubPermissions{
			Send: []string{"broadcast"},
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := base64.StdEncoding.DecodeString(os.Getenv("EXT_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}
