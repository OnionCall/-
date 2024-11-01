package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/OnionCall/squa/common"
	"github.com/OnionCall/squa/services"
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

func(g MessageGroup) CreateGroup() int {
	group := createRoomDto{GroupUuid: g.GroupUuid.String(), GroupKey: g.GroupKey}
	url := fmt.Sprintf("%v/admin/messagegroup/", common.Environment)
	contentType := "application/json"

	jsonData, err := json.Marshal(group)
	if err != nil {
		common.AddError(err)
		return 0
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create group request: %v", err)
		common.AddError(err)
		return 0
	}

	resp, err := services.Authorize(req, contentType)
	if err != nil || resp.StatusCode != 201 {
		log.Printf("%v Failed to create group: %v", resp.StatusCode, err)
		common.AddError(err)
		return 0
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.AddError(err)
	}

	var response groupResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		common.AddError(err)
	}

	g.GroupId = response.GroupId
	g.setGroup()

	return g.GroupId
}

func (g MessageGroup) GetGroupByLogin() int {
	url := fmt.Sprintf("%v/admin/messagegroup/?groupuuid=%v&groupkey=%v", common.Environment, g.GroupUuid.String(), g.GroupKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		common.AddError(err)
	}

	resp, err := services.Authorize(req, "")
	if err != nil || resp.StatusCode != 200 {
		log.Printf("%v Failed to get group: %v", resp.StatusCode, err)
		common.AddError(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.AddError(err)
	}

	var group groupResponse
	err = json.Unmarshal(body, &group)
	if err != nil {
		common.AddError(err)
	}

	g.GroupId = group.GroupId
	g.setGroup()

	return group.GroupId
}
