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

type UserDetails struct {
	UserId      int
	DisplayName string
	GroupId  int
}

var User UserDetails

func(u UserDetails) setUser() UserDetails {
	User = UserDetails{
		UserId: u.UserId,
		GroupId: u.GroupId, 
		DisplayName: u.DisplayName,
	}
	return User
}

func(u UserDetails) CreateUser() UserDetails {
	user := createUserDto{GroupId: u.GroupId, DisplayName: u.DisplayName}
	url := fmt.Sprintf("%v/admin/userdetails/", common.Env)
	contentType := "application/json"
	
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return UserDetails{}
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
	}

	resp, err := authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to create user: %v", resp.StatusCode, err)
		return UserDetails{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return UserDetails{}
	}

	var ur userResponse
	err = json.Unmarshal(body, &ur)
	if err != nil {
		log.Printf("Failed to convert from json: %v",err)
		return UserDetails{}
	}

	u.UserId = ur.UserId
	createdUser := u.setUser()

	return createdUser
}

func (u UserDetails) DeactivateUser() {
	url := fmt.Sprintf("%v/admin/userdetails/", common.Env)
	contentType := "application/json"
	
	jsonData, err := json.Marshal(u)
	if err != nil {
		log.Println(err)
	}
	
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
	}

	resp, err := authorize(req, contentType)
	if err != nil {
		log.Printf("%v Failed to deactivate user: %v", resp.StatusCode, err)
	} else if resp.StatusCode == 200 {
		//User Deactivated
	} else if resp.StatusCode == 204 {
		//Group Deleted
	}
}