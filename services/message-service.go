package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/onioncall/cli-squa/cli/common"
)

type LoginResponse struct {
	IsValid bool `json:"isvalid"`
}

type sendMessageDto struct {
	UserId int `json:"userid"`
	GroupId int `json:"groupid"`
	MessageContents string `json:"messagecontents"`
}

type messageResponse struct {
	MessageId int `json:"messageid"`
}

func MessagesService() {
	ticker := time.Tick(5 * time.Second) // target is twice a second for production

	for range ticker {
		url := fmt.Sprintf("%v/admin/messages/?messageid=%v&groupid=%v", common.Env, common.LatestMessageId, common.Group.GroupId)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Failed to create message request: %v",err) 
		}
		
		resp, err := authorize(req, "")
		if err != nil || resp.StatusCode != 200 {
			log.Printf("%v Failed to get messages from api: %v", resp.StatusCode, err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("Failed to get messages from api: %v",resp.Status)
			return
		}
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		
		var messages []common.DisplayMessage
		err = json.Unmarshal(body, &messages)
		if err != nil {
			log.Printf("Failed to convert from json: %v",err)
			return
		}
		
		if len(messages) > 0 {
			common.SetLatestMessageId(messages[len(messages)-1].MessageId)
			
			//We only want to render messages not sent from user to the terminal, since we already have those message in the terminal
			for _, message := range messages {
				if message.DisplayName != common.User.DisplayName {
					common.AddMessages(message)
				}
			}
		}
	}
} 

func SendMessage(messageContents string) {
	group := sendMessageDto {
		UserId: common.User.UserId,
		GroupId: common.User.GroupId,
		MessageContents: messageContents,
	}

	url := fmt.Sprintf("%v/admin/messages/", common.Env)
	contentType := "application/json"

	jsonData, err := json.Marshal(group)
	if err != nil {
		fmt.Println(err)
		return 
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create message request: %v",err)
		return
	}

	resp, err := authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to send message: %v", resp.StatusCode, err)
		return
	}

	defer resp.Body.Close()
	var response messageResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}
}