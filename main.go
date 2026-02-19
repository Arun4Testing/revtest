package main

import (
	"log"
	"os"
	"strconv"
	"time"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var appID string
var appCertificate string

func init() {
	// Load .env file for local development only
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (this is fine in production)")
	}
}

func main() {
	appID = os.Getenv("APP_ID")
	appCertificate = os.Getenv("APP_CERTIFICATE")

	if appID == "" || appCertificate == "" {
		log.Fatal("APP_ID or APP_CERTIFICATE missing")
	}

	r := gin.Default()
	r.Use(nocache())

	r.GET("/rtc/:channel/:role/:uid", getRtcToken)
	r.GET("/rtm/:uid", getRtmToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local fallback
	}

	log.Println("Server running on port:", port)
	r.Run(":" + port)
}

func nocache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func getRtcToken(c *gin.Context) {
	channel := c.Param("channel")
	roleStr := c.Param("role")
	uidStr := c.Param("uid")

	uid64, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}
	uid := uint32(uid64)

	expire := uint32(time.Now().Unix() + 3600)

	var role rtctokenbuilder.Role
	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		role = rtctokenbuilder.RoleSubscriber
	}

	token, err := rtctokenbuilder.BuildTokenWithUid(
		appID,
		appCertificate,
		channel,
		uid,
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

func getRtmToken(c *gin.Context) {
	uid := c.Param("uid")

	expire := uint32(time.Now().Unix() + 3600)

	token, err := rtmtokenbuilder.BuildToken(
		appID,
		appCertificate,
		uid,
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
