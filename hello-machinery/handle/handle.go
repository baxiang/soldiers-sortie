package handle

import (
	"fmt"
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

//Add test task method add
func Add(c *gin.Context, s *machinery.Server) {
	var (
		uid = uuid.New().String()
	)

	signature := &tasks.Signature{
		IgnoreWhenTaskNotRegistered: true,
		UUID: uid,
		Name: "add",
		Args: []tasks.Arg{
			{
				Type:  "int64",
				Value: 1,
			},
			{
				Type:  "int64",
				Value: 1,
			},
		},
	}
	asyncResult, err := s.SendTask(signature)
	res, err := asyncResult.Get(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": tasks.HumanReadableResults(res), "uuid": uid})
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks.HumanReadableResults(res), "uuid": uid})


}
//Add test task method longRunningTask
func LongRunningTask(c *gin.Context, s *machinery.Server) {
	var (
		uid = uuid.New().String()
	)

	signature := &tasks.Signature{
		Name: "long_running_task",
	}
	asyncResult, err := s.SendTask(signature)
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"longRunningTask": err, "uuid": uid})
	fmt.Println(asyncResult)

}