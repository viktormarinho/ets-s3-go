package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	sessionId, err := ctx.Cookie("sessionId")

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"err": "No cookie named sessionId"})
		return
	}

	postBody, err := json.Marshal(map[string]string{"sessionId": sessionId})

	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "error converting sessionId to json", "err": err.Error()})
		return
	}

	resp, err := http.Post("http://localhost:5000/ms/me", "application/json", bytes.NewBuffer(postBody))

	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"message": "error when requested session", "err": err.Error()})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"message": "error reading auth service response body", "err": err.Error()})
		return
	}

	var authResponse AuthResponse

	json.Unmarshal(body, &authResponse)

	if authResponse.User.ID == nil {
		ctx.AbortWithStatusJSON(401, gin.H{"message": "User session not found"})
		return
	}

	ctx.Set("currentUser", authResponse.User)

	ctx.Next()
}
