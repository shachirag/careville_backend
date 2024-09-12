package zegoCloudChatHandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CreateUserInfo struct {
	UserId     string `json:"UserId"`
	UserName   string `json:"UserName,omitempty"`
	UserAvatar string `json:"UserAvatar,omitempty"`
	Extra      string `json:"Extra,omitempty"`
}

type CreateUserInfoRequest struct {
	UserInfo []CreateUserInfo `json:"UserInfo"`
}

type ZegoUserRegisterResponse struct {
	Code      int    `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	ErrorList []struct {
		UserId     string `json:"UserId"`
		SubCode    int    `json:"SubCode"`
		SubMessage string `json:"SubMessage"`
	} `json:"ErrorList"`
}

var ZegoCloudAPIUrl = "https://zim-api.zego.im/?Action=UserRegister"

func RegisterZegoCloudUsers(userInfos []CreateUserInfo) (*ZegoUserRegisterResponse, error) {

	reqBody := CreateUserInfoRequest{
		UserInfo: userInfos,
	}

	requestBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("POST", ZegoCloudAPIUrl, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ZEGO_APP_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var zegoUserRegisterResponse ZegoUserRegisterResponse
	err = json.Unmarshal(body, &zegoUserRegisterResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	if zegoUserRegisterResponse.Code != 0 {
		return nil, fmt.Errorf("user registration failed: %s", zegoUserRegisterResponse.Message)
	}

	return &zegoUserRegisterResponse, nil
}
