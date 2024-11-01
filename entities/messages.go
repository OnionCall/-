package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/OnionCall/squa/common"	
	"github.com/OnionCall/squa/services"
)

type LoginResponse struct {
	IsValid bool `json:"isvalid"`
}

type sendMessageDto struct {
	DisplayName string `json:"displayname"`
	GroupId int `json:"groupid"`
	MessageContents string `json:"messagecontents"`
}

type messageResponse struct {
	MessageId int `json:"messageid"`
}

type DisplayMessage struct {
	MessageId int `json:"messageid"`
	DisplayName string 	`json:"displayname"`
	MessageContents string `json:"messagecontents"`
}

var latestMessageId int

func setLatestMessageId(messageId int) int {
	latestMessageId = messageId
	return latestMessageId
}

var UnrecievedMessages []DisplayMessage

func setMessages(m DisplayMessage) []DisplayMessage {
	UnrecievedMessages = append(UnrecievedMessages, m)
	return UnrecievedMessages
}

func MessagesService() {
	ticker := time.Tick(5 * time.Second) // target is twice a second for production

	for range ticker {
		url := fmt.Sprintf("%v/admin/messages/?messageid=%v&groupid=%v", common.Environment, latestMessageId, Group.GroupId)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Failed to create message request: %v",err) 
			common.AddError(err)
		}
		
		resp, err := services.Authorize(req, "")
		if err != nil || resp.StatusCode != 200 {
			log.Printf("%v Failed to get messages from api: %v", resp.StatusCode, err)
			common.AddError(err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("Failed to get messages from api: %v",resp.Status)
			common.AddError(err)
			return
		}
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			common.AddError(err)
			return
		}
		
		var messages []DisplayMessage
		err = json.Unmarshal(body, &messages)
		if err != nil {
			log.Printf("Failed to convert from json: %v",err)
			common.AddError(err)
			return
		}
		
		if len(messages) > 0 {
			setLatestMessageId(messages[len(messages)-1].MessageId)
			
			//We only want to render messages not sent from user to the terminal, since we already have those message in the terminal
			for _, message := range messages {
				if message.DisplayName != User.DisplayName {
					setMessages(message)
				}
			}
		}
	}
} 

func (m DisplayMessage) SendMessage() {
	group := sendMessageDto {
		DisplayName: User.DisplayName,
		GroupId: User.GroupId,
		MessageContents: m.MessageContents,
	}

	url := fmt.Sprintf("%v/admin/messages/", common.Environment)
	contentType := "application/json"

	jsonData, err := json.Marshal(group)
	if err != nil {
		common.AddError(err)
		return 
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create message request: %v",err)
		common.AddError(err)
		return
	}

	resp, err := services.Authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to send message: %v", resp.StatusCode, err)
		common.AddError(err)
		return
	}

	defer resp.Body.Close()
	var response messageResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.AddError(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		common.AddError(err)
	}

	setLatestMessageId(response.MessageId)
}
