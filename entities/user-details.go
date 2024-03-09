package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io"
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
	DisplayName string
	GroupId  int
}

var User UserDetails

func(u UserDetails) setUser() UserDetails {
	User = UserDetails{
		GroupId: u.GroupId, 
		DisplayName: u.DisplayName,
	}
	return User
}

func(u UserDetails) CreateUser() {
	u.setUser()
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
