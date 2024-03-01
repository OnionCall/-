package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/onioncall/squa/common"
	"github.com/onioncall/squa/services"
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
	url := fmt.Sprintf("%v/admin/userdetails/", common.Environment)
	contentType := "application/json"
	
	jsonData, err := json.Marshal(user)
	if err != nil {
		common.AddError(err)
		return UserDetails{}
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		common.AddError(err)
	}

	resp, err := services.Authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to create user: %v", resp.StatusCode, err)
		common.AddError(err)
		return UserDetails{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.AddError(err)
		return UserDetails{}
	}

	var ur userResponse
	err = json.Unmarshal(body, &ur)
	if err != nil {
		log.Printf("Failed to convert from json: %v",err)
		common.AddError(err)
		return UserDetails{}
	}

	u.UserId = ur.UserId
	createdUser := u.setUser()

	return createdUser
}

func (u UserDetails) DeactivateUser() {
	url := fmt.Sprintf("%v/admin/userdetails/", common.Environment)
	contentType := "application/json"
	
	jsonData, err := json.Marshal(u)
	if err != nil {
		common.AddError(err)
	}
	
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		common.AddError(err)
	}

	resp, err := services.Authorize(req, contentType)
	if err != nil {
		log.Printf("%v Failed to deactivate user: %v", resp.StatusCode, err)
		common.AddError(err)
	} else if resp.StatusCode == 200 {
		//User Deactivated
	} else if resp.StatusCode == 204 {
		//Group Deleted
	} else {
		log.Print("I don't know how you got here, or what it means. But you're here. And it means something. I just don't know what.")
	}
}