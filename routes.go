package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	_ "io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	ERRORED  = "errored" // or: "finished", "canceled"
	FINISHED = "finished"
	CANCELED = "canceled"
	IOS      = "ios"
	Android  = "android"
)

type Artifacts struct {
	BuildUrl        string `form:"buildUrl"`
	LogsS3KeyPrefix string `form:"logsS3KeyPrefix"`
}

type Metadata struct {
	AppName            string `form:"appName"`
	Username           string `form:"username"`
	AppVersion         string `form:"appVersion"`
	AppBuildVersion    string `form:"appBuildVersion"`
	BuildProfile       string `form:"buildProfile"`
	SdkVersion         string `form:"sdkVersion"`
	RuntimeVersion     string `form:"RuntimeVersion"`
	Channel            string `form:"channel"`
	ReactNativeVersion string `form:"reactNativeVersion"`
}

type WebhookPayload struct {
	Id                  string `form:"id"`
	AccountName         string `form:"accountName" ` // "accountName": "dsokal",
	ProjectName         string `form:"projectName"`
	BuildDetailsPageUrl string `form:"buildDetailsPageUrl"`
	ParentBuildId       string `form:"parentBuildId"`
	AppId               string `form:"appId"`
	Platform            string `form:"platform"`
	Status              string `form:"status"`
	Artifacts           Artifacts
	Metadata            Metadata
	Message             string `form:"message"`
}

func handleHook(c *gin.Context) {
	var artifacts WebhookPayload
	c.ShouldBindBodyWith(&artifacts, binding.JSON)
	body, _ := c.GetRawData()
	secretKey := os.Getenv("SECRET_WEBHOOK_KEY")
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write(body)
	hash := h.Sum(nil)
	sign := c.GetHeader("expo-signature")
	sha1 := fmt.Sprintf("sha1=%s", hex.EncodeToString(hash[:]))
	log.Println(sha1, sign)
	compareResult := subtle.ConstantTimeCompare([]byte(sign), []byte(sha1))
	if compareResult == 0 {
		c.String(http.StatusForbidden, "go away!")
		return
	}
	go trySend(artifacts)
	c.JSON(http.StatusOK, "Success")
}
