package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type StatusInfo struct {
	executingTime    time.Time
	isExecutingATask bool
}

func main() {
	statusInfo := &StatusInfo{}
	statusInfoChannel := make(chan StatusInfo, 1)

	//Subgoroutine - loop to get task and process it
	go loopToProcessTask(statusInfoChannel)
	go loopToUpdateStatusInfo(statusInfo, statusInfoChannel)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		var message string
		var mu sync.Mutex
		mu.Lock()
		currentStatusInfo := statusInfo
		mu.Unlock()

		if !currentStatusInfo.isExecutingATask {
			message = fmt.Sprintf("There aren't tasks running")
		} else {
			message = fmt.Sprintf("Executing time is taking: %s", time.Since(currentStatusInfo.executingTime))
		}

		c.JSON(http.StatusOK, gin.H{
			"isExecuting": statusInfo.isExecutingATask,
			"message":     message,
		})
	})
	r.Run()
}

func loopToUpdateStatusInfo(statusInfo *StatusInfo, statusInfoChannel chan StatusInfo) {
	var mu sync.Mutex
	for {
		select {
		case statusOfChannel := <-statusInfoChannel:
			if !statusOfChannel.isExecutingATask {
				fmt.Println("Updating the channel notifying that there's no task being executing")
			} else {

				fmt.Println("Updating the channel notifying that there's a task being executing")
			}
			mu.Lock()
			statusInfo.executingTime = statusOfChannel.executingTime
			statusInfo.isExecutingATask = statusOfChannel.isExecutingATask
			mu.Unlock()
		}
	}
}

func loopToProcessTask(statusInfoChannel chan StatusInfo) {
	for {
		fmt.Println("Changing status info to notify that I'm not executing any task")
		statusInfo := StatusInfo{executingTime: time.Now(), isExecutingATask: false}
		statusInfo.executingTime = time.Now()
		statusInfo.isExecutingATask = false
		statusInfoChannel <- statusInfo

		time.Sleep(5 * time.Second)

		fmt.Println("Changing status info to notify that I'm about to execute a task")
		statusInfo.executingTime = time.Now()
		statusInfo.isExecutingATask = true
		statusInfoChannel <- statusInfo

		//This represents the execution of a test
		fmt.Println("I'm executing a task...")
		time.Sleep(10 * time.Second)
	}
}
