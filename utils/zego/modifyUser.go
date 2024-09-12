package zegoCloudChatHandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ModifyUserInfoReq struct {
	UserId     string `json:"UserId"`
	UserName   string `json:"UserName,omitempty"`   
	UserAvatar string `json:"UserAvatar,omitempty"`
	Extra      string `json:"Extra,omitempty"`     
}

type ModifyUserInfoRequest struct {
	UserInfo []ModifyUserInfoReq `json:"UserInfo"`
}

type ModifyUserInfoResponse struct {
	Code      int        `json:"Code"`
	Message   string     `json:"Message"`
	RequestId string     `json:"RequestId"`
	ErrorList []ErrorObj `json:"ErrorList"`
}

type ErrorObj struct {
	UserId     string `json:"UserId"`
	SubCode    int    `json:"SubCode"`
	SubMessage string `json:"SubMessage"`
}

func ModifyUserInfo(users []ModifyUserInfoReq) (*ModifyUserInfoResponse, error) {
	apiUrl := "https://zim-api.zego.im/?Action=ModifyUserInfo"

	reqBody := ModifyUserInfoRequest{
		UserInfo: users,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var modifyUserInfoResp ModifyUserInfoResponse
	err = json.Unmarshal(body, &modifyUserInfoResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if modifyUserInfoResp.Code != 0 {
		return &modifyUserInfoResp, errors.New(modifyUserInfoResp.Message)
	}

	return &modifyUserInfoResp, nil
}

