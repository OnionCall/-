package common

// import "github.com/google/uuid"

// type UserDetails struct {
// 	UserId      int
// 	DisplayName string
// 	GroupId  int
// }

// var User UserDetails

// func CreateUser (userId int, groupId int, displayName string) UserDetails {
// 	User = UserDetails{
// 		UserId: userId,
// 		GroupId: groupId, 
// 		DisplayName: displayName,
// 	}
// 	return User
// }

// type MessageGroup struct {
// 	GroupId   int       `json:"groupid"`
// 	GroupUuid uuid.UUID `json:"groupuuid"`
// 	GroupKey string `json:"groupkey"`
// }

// var Group MessageGroup

// func CreateGroup (groupId int, uuid uuid.UUID, groupKey string) MessageGroup {
// 	Group = MessageGroup {
// 		GroupId: groupId, 
// 		GroupUuid: uuid, 
// 		GroupKey: groupKey,
// 	}
// 	return Group
// }

// type DisplayMessage struct {
// 	MessageId int `json:"messageid"`
// 	DisplayName string 	`json:"displayname"`
// 	MessageContents string `json:"messagecontents"`
// }

// var LatestMessageId int

// func SetLatestMessageId(messageId int) int {
// 	LatestMessageId = messageId
// 	return LatestMessageId
// }

// var GroupId int

// var UnrecievedMessages []DisplayMessage

// func AddMessages(msg DisplayMessage) []DisplayMessage {
// 	UnrecievedMessages = append(UnrecievedMessages, msg)
// 	return UnrecievedMessages
// }

var Env string

type Environment struct {
	Develop string
	Production string
}

func GetEnvironment() Environment {
	environment := Environment {
		Develop: "http://localhost:8040",
		Production: "https://clisqua.com",
	}

	return environment
}

func SetEnvironment(env string) string {
	Env = env
	return Env
}

// // func (u *UnrecievedMessages) AddMessages(msg DisplayMessage) {
// // 	*u = append(*u, msg)
// // }
