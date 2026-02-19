package main

import (
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"time"

	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

var appID = "62aeb77b6c704c79919ac9852bc4e24b"
var appCertificate = "b92d623a025a499eab5e30f5396a01e2"

// Converts string user ID to numeric UID for Agora
func uidFromString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func main() {
	r := gin.Default()
	r.Use(nocache())

	r.GET("/rtc/:channel/:role/:uid", getRtcToken)
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

	uid := uidFromString(uidStr) // Convert string user ID to uint32

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtmToken": token,
	})
}
