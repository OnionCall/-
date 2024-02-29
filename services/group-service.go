package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/onioncall/cli-squa/cli/common"
)

type createRoomDto struct {
	GroupUuid string `json:"groupuuid"`
	GroupKey  string `json:"groupkey"`
}

type groupResponse struct {
	GroupId int `json:"groupid"`
}

type MessageGroup struct {
	GroupId   int       `json:"groupid"`
	GroupUuid uuid.UUID `json:"groupuuid"`
	GroupKey string `json:"groupkey"`
}

var Group MessageGroup

func(g MessageGroup) setGroup () MessageGroup {
	Group = MessageGroup {
		GroupId: g.GroupId, 
		GroupUuid: g.GroupUuid, 
		GroupKey: g.GroupKey,
	}

	return Group
}

type GroupService interface {

}

func(g MessageGroup) CreateGroup() int {
	group := createRoomDto{GroupUuid: g.GroupUuid.String(), GroupKey: g.GroupKey}
	url := fmt.Sprintf("%v/admin/messagegroup/", common.Env)
	contentType := "application/json"

	jsonData, err := json.Marshal(group)
	if err != nil {
		log.Println(err)
		return 0
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create group request: %v", err)
		return 0
	}

	resp, err := authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to create group: %v", resp.StatusCode, err)
		return 0
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var response groupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}

	g.GroupId = response.GroupId
	g.setGroup()

	return g.GroupId
}

func (g MessageGroup) GetGroupByLogin() int {
	url := fmt.Sprintf("%v/admin/messagegroup/?groupuuid=%v&groupkey=%v", common.Env, g.GroupUuid.String(), g.GroupKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	resp, err := authorize(req, "")
	if err != nil || resp.StatusCode != 200 {
		log.Printf("%v Failed to get group: %v", resp.StatusCode, err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var group groupResponse
	err = json.Unmarshal(body, &group)
	if err != nil {
		log.Println(err)
	}

	g.GroupId = group.GroupId
	g.setGroup()

	return group.GroupId
}
