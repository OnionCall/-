package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/onioncall/cli-squa/cli/common"
)

type createUserDto struct {
	GroupId int `json:"groupid"`
	DisplayName string `json:"displayname"`
}

type userResponse struct {
	UserId int `json:"userid"`
}

func CreateUser(groupId int, displayName string) common.UserDetails {
	user := createUserDto{GroupId: groupId, DisplayName: displayName}
	url := fmt.Sprintf("%v/admin/userdetails/", common.Env)
	contentType := "application/json"
	
	
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return common.UserDetails{}
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
	}

	resp, err := authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to create user: %v", resp.StatusCode, err)
		return common.UserDetails{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return common.UserDetails{}
	}

	var ur userResponse
	err = json.Unmarshal(body, &ur)
	if err != nil {
		log.Printf("Failed to convert from json: %v",err)
		return common.UserDetails{}
	}

	createdUser := common.CreateUser(ur.UserId, groupId, displayName)
	return createdUser
}