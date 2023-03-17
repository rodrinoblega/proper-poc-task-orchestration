package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

type StatusInfo struct {
	executingTimePtr    *time.Time
	isExecutingATaskPtr *bool
}

func main() {
	//I have to use Mutex because I'll be reading/updating StatusInfo struct from different routines
	var mu sync.Mutex
	statusInfo := StatusInfo{}

	//Subgoroutine - loop to get task and process it
	go func() {
		executingTime := time.Now()
		isExecutingATask := false

		mu.Lock()
		statusInfo.executingTimePtr = &executingTime
		statusInfo.isExecutingATaskPtr = &isExecutingATask
		mu.Unlock()

		for {
			executingTime = time.Now()
			isExecutingATask = true

			fmt.Println("Changing status info to notify that I'm executing a task")
			mu.Lock()
			statusInfo.executingTimePtr = &executingTime
			statusInfo.isExecutingATaskPtr = &isExecutingATask
			mu.Unlock()

			//This represents the execution of a test
			fmt.Println("I'm executing a task...")
			time.Sleep(10 * time.Second)

			isExecutingATask = false
			executingTime = time.Now()

			fmt.Println("Changing status info to notify that I'm not executing any task")
			mu.Lock()
			statusInfo.executingTimePtr = &executingTime
			statusInfo.isExecutingATaskPtr = &isExecutingATask
			mu.Unlock()

			time.Sleep(2 * time.Second)
		}
	}()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		mu.Lock()
		executingTime := *statusInfo.executingTimePtr
		isExecutingATask := *statusInfo.isExecutingATaskPtr
		mu.Unlock()
		var message string
		if !isExecutingATask {
			message = fmt.Sprintf("There aren't tasks running")
		} else {
			message = fmt.Sprintf("Executing time is taking: %s", time.Since(executingTime))
		}

		c.JSON(http.StatusOK, gin.H{
			"isExecuting": isExecutingATask,
			"message":     message,
		})
	})
	r.Run()
}
