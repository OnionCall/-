package services

import (
	"fmt"
	"os"
	"time"
)

func BeginSession() {
	sessionDuration := 5 * time.Minute
	timer := time.NewTimer(sessionDuration)
	
	<-timer.C
	Clear()
	fmt.Println("Session has expired")
	os.Exit(0)
}
