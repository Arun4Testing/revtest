package main

import (
	"log"
	"net/http"
	"os"
	"time"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

// Read from environment variables
var appID = "62aeb77b6c704c79919ac9852bc4e24b"
var appCertificate = "b92d623a025a499eab5e30f5396a01e2"

func main() {

	if appID == "" || appCertificate == "" {
		log.Fatal("AGORA_APP_ID or AGORA_APP_CERTIFICATE not set in environment")
	}

	r := gin.Default()
	r.Use(corsMiddleware())

	// RTC token using USER ACCOUNT (string UID)
	r.GET("/rtc/:channel/:role/:uid", getRtcToken)

	// RTM token
	r.GET("/rtm/:uid", getRtmToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// func corsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type")
// 		c.Header("Access-Control-Allow-Methods", "GET")
// 		c.Next()
// 	}
// }

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// ---------------- RTC TOKEN ----------------

func getRtcToken(c *gin.Context) {

	channel := c.Param("channel")
	roleStr := c.Param("role")
	userAccount := c.Param("uid")

	expire := uint32(time.Now().Unix() + 3600)

	var role rtctokenbuilder.Role
	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		role = rtctokenbuilder.RoleSubscriber
	}

	//  IMPORTANT: Use BuildTokenWithUserAccount
token, err := rtctokenbuilder.BuildTokenWithAccount(
	appID,
	appCertificate,
	channel,
	userAccount,
	role,
	expire,
)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtcToken": token,
	})
}

// ---------------- RTM TOKEN ----------------

func getRtmToken(c *gin.Context) {

	userAccount := c.Param("uid")
	expire := uint32(time.Now().Unix() + 3600)

	token, err := rtmtokenbuilder.BuildToken(
		appID,
		appCertificate,
		userAccount,
		expire,
		"",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtmToken": token,
	})
}
