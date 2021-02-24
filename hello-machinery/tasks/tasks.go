package tasks

import (
	"log"
	"time"
)

func Add(args ...int64)(int64, error){
	var sum int64
	for _,arg :=range args{
		sum +=arg
	}
	//log.Println("add result:",sum)
	return sum,nil
}
func LongRunningTask() error {
	log.Println("Long running task started")
	for i := 0; i < 100; i++ {
		log.Println(10 - i)
		time.Sleep(1 * time.Second)
	}
	log.Println("Long running task finished")
	return nil
}