package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	numberOfComputers = 10
	numberOfClient = 25
)

func workingSimulation() time.Duration {
	// lets assume 3 seconds = 30 minutes
	minOnlineDuration := 3*time.Second.Seconds()
	// lets assue 6 Seconds = 60 minutes
	maxOnlineDuration := 6*time.Second.Seconds()
	ActualWorkingTime := int(minOnlineDuration)+rand.Intn(int(maxOnlineDuration-minOnlineDuration))	
	//lets start working 
	time.Sleep(time.Duration(ActualWorkingTime)*time.Second)
	//let return the duration used for working 
	return time.Duration(ActualWorkingTime)*time.Second
}


type client struct{
	token int 
}

func computerSimulation(clients chan client, close chan struct{},file *os.File){
	for computerId:=0; computerId < numberOfComputers; computerId++{
		go func(cId int){
			for{
				select {
				case <-close:
					return
				case letInClient := <-clients:
					file.WriteString("<p style=\"color:blue;\">Client with token "+strconv.Itoa(letInClient.token)+" is Online</p>")
					workingTime := workingSimulation()
					file.WriteString("<p style=\"color:green;\">Client with token "+strconv.Itoa(letInClient.token)+" is done, having spent "+workingTime.String()+ " time</p>")					
				}
			}
		}(computerId)
	}
}

func clientQueueSimulation(queue chan client,file *os.File){
	for tokengen:=1;tokengen<=numberOfClient;tokengen++{
		client_with_token := client{token:tokengen}
		go func(client_token int){
			select {
			case queue <- client_with_token:
			default:
				file.WriteString("<p style=\"color:red;\">Client with token "+strconv.Itoa(client_token)+" is waiting for turn.</p>")
				queue <-client_with_token
			}
		}(tokengen)
	}
}


func runcafe() {
	os.Mkdir("static",0755)
	file,err := os.Create("static/index.html")
	if err != nil {
		log.Fatalf("failed creating file: %s",err)
	}
	file.WriteString("<html>")
	defer file.Close()
	defer file.WriteString("</html>")
	// start := 
	// this channel represent the queue of client
	clientChannel := make(chan client)
	// this channel will be used to close the shop and shut all the processes.
	closeChannel := make(chan struct{})

	computerSimulation(clientChannel,closeChannel,file)
	clientQueueSimulation(clientChannel,file)

	time.Sleep(10*time.Second)
	file.WriteString("<p>Please Empty the Area, We are closing now.</p>")
	close(closeChannel)
}


func main() {
	runcafe()
	time.Sleep(10*time.Second)
	fmt.Println("open the link http://localhost:8080")
	http.Handle("/",http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080",nil))
}