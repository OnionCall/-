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

func CreateGroup(uuid uuid.UUID, groupKey string) int {
	group := createRoomDto{GroupUuid: uuid.String(), GroupKey: groupKey}
	url := fmt.Sprintf("%v/admin/messagegroup/", common.Env)
	contentType := "application/json"

	jsonData, err := json.Marshal(group)
	if err != nil {
		log.Println(err)
		return 0
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
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

	common.CreateGroup(response.GroupId, uuid, groupKey)
	return response.GroupId
}

func GetGroupByLogin(uuid uuid.UUID, groupKey string) int {
	url := fmt.Sprintf("%v/admin/messagegroup/?groupuuid=%v&groupkey=%v", common.Env, uuid.String(), groupKey)
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

	common.CreateGroup(group.GroupId, uuid, groupKey)
	return group.GroupId
}
