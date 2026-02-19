package main

import (
	"log"
	"os"
	"time"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

var appID string
var appCertificate string

func main() {
	// Hardcoded App ID and Certificate
	appID = "62aeb77b6c704c79919ac9852bc4e24b"
	appCertificate = "b92d623a025a499eab5e30f5396a01e2"

	r := gin.Default()
	r.Use(nocache())

	// RTC token for publisher/subscriber
	r.GET("/rtc/:channel/:role/:uid", getRtcToken)

	// RTM token
	r.GET("/rtm/:uid", getRtmToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	r.Run(":" + port)
}

// No cache middleware
func nocache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// RTC token generator for userAccount
func getRtcToken(c *gin.Context) {
	channel := c.Param("channel")
	roleStr := c.Param("role")
	userAccount := c.Param("uid") // string UID

	expire := uint32(time.Now().Unix() + 3600)

	var role rtctokenbuilder.Role
	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		role = rtctokenbuilder.RoleSubscriber
	}

	token, err := rtctokenbuilder.BuildTokenWithUserAccount(
		appID,
		appCertificate,
		channel,
		userAccount,
		role,
		expire,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"rtcToken": token,
	})
}

// RTM token generator
func getRtmToken(c *gin.Context) {
	userAccount := c.Param("uid")

	expire := uint32(time.Now().Unix() + 3600)

	token, err := rtmtokenbuilder.BuildToken(
		appID,
		appCertificate,
		userAccount,
		expire,
		"", // required salt
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"rtmToken": token,
	})
}
