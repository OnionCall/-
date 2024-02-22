package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// "net/url"

	// "strings"

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
	url := fmt.Sprintf("%v/messagegroup/", common.Env)
	contentType := "application/json"

	jsonData, err := json.Marshal(group)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	r, err := http.Post(url, contentType, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return 0
	}

	defer r.Body.Close()

	var response groupResponse

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}

	common.CreateGroup(response.GroupId, uuid, groupKey)

	return response.GroupId
}

func GetGroupByLogin(uuid uuid.UUID, groupKey string) int {
	url := fmt.Sprintf("%v/messagegroup/?groupuuid=%v&groupkey=%v", common.Env, uuid.String(), groupKey)
	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	var group groupResponse

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(body, &group)
	if err != nil {
		log.Println(err)
	}

	common.CreateGroup(group.GroupId, uuid, groupKey)
	return group.GroupId
}
