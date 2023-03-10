package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	DING_REBOT_API   = "https://oapi.dingtalk.com/robot/send?access_token="
	DING_ROBOT_TOKEN = "DING_ROBOT_TOKEN"
	MARKDOWN         = "markdown"
)

type DingMessageMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingMessage struct {
	Msgtype  string              `json:"msgtype"`
	Markdown DingMessageMarkdown `json:"markdown"`
}

func NewMarkdownMsg() DingMessage {
	return DingMessage{
		Msgtype: MARKDOWN,
	}
}

func trySend(artifacts WebhookPayload) {
	// ๐release_testing_version
	title := fmt.Sprintf("๐release_%s_%s_%s", artifacts.Metadata.BuildProfile, artifacts.Platform, artifacts.Metadata.AppVersion)
	dingMessage := NewMarkdownMsg()
	dingMessage.Markdown.Title = title
	if artifacts.Status == FINISHED {
		buildDetailUrl := artifacts.Artifacts.BuildUrl
		pngByte, _ := qrcode.Encode(buildDetailUrl, qrcode.Medium, 256)
		pngBase64 := base64.StdEncoding.EncodeToString(pngByte)
		pngFile := fmt.Sprintf("%s%s", "data:image/png;base64,", pngBase64)
		content := fmt.Sprintf("### yimi \n  %s  \n > ![](%s)", title, pngFile)
		dingMessage.Markdown.Text = content
	} else if artifacts.Status == ERRORED {
		dingMessage.Markdown.Text = fmt.Sprintf("### yimi \n  %s--โ  \n  ๐๐๐๐๐๐ๅบ็ฉๅ็ซฏๅๅๅธๅคฑ่ดฅไบๅง", title)
	} else {
		// cancel
		return
	}

	token := os.Getenv(DING_ROBOT_TOKEN)
	remoteApi := fmt.Sprintf("%s%s", DING_REBOT_API, token)
	body, _ := json.Marshal(dingMessage)
	res, err := http.Post(remoteApi, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Println(err)
	}
	resByte, _ := io.ReadAll(res.Body)
	log.Println(string(resByte))
}
