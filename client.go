package main

import (
	service "github.com/ingridkarinaf/ActiveReplicationTemplate/interface"
	grpc "google.golang.org/grpc"
	"log"
	"os"
	"bufio"
	"context"
	"reflect"
	"fmt"
)

/* 
Responsible for:
	1. Making connection to FE, has to redial if they lose connection
	2. 
Limitations:
	1. Can only dial to pre-determined front-ends or by incrementing 
	(then there is no guarantee that there is an FE with that port number)
	2. Assumes a failed request is due to a crashed server, redials immediately
*/

var server service.ServiceClient
var connection *grpc.ClientConn 

func main() {

	//Creating log file
	f := setLogClient()
	defer f.Close()

	FEport := ":" + os.Args[1] //dial 4000 or 4001 (available ports on the FEServers)
	conn, err := grpc.Dial(FEport, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	connection = conn

	server = service.NewServiceClient(connection) //creates a connection with an FE server
	defer connection.Close()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		
		for {
			println("Enter 'increment' to increment value, 'get' to get value. (without quotation marks)")
			scanner.Scan()
			textChoice := scanner.Text()
			if (textChoice == "increment") {
				serviceUpdate := &service.UpdateRequest{}

				result := Update(serviceUpdate)
				if result.Outcome == true {
					fmt.Printf("Value successfully incremented to %v\n", result.CurrentValue)
					log.Printf("Value successfully incremented to %v\n", result.CurrentValue)
				} else {
					log.Println("Update unsuccessful, please try again.")
				}
				
			} else if (textChoice == "get") {
				
				getReq := &service.RetrieveRequest{}

				result := Retrieve(getReq) 
				log.Printf("Client: Value: %v \n", int(result))
				fmt.Printf("Value: %v \n", int(result))
			} else {
				log.Println("Sorry, didn't catch that. ")
				fmt.Println("Sorry, didn't catch that. ")
			}
		}
	}()

	for {}
}

//If function returns an error, redial to other front-end and try again
func Update(hashUpt *service.UpdateRequest) (*service.UpdateReply) {
	result, err := server.Update(context.Background(), hashUpt) //What does the context.background do?
	if err != nil {
		log.Printf("Client %s hashUpdate failed:%s. \n Redialing and retrying. \n", connection.Target(), err)
		Redial()
		return Update(hashUpt)
	}
	return result
}

func Retrieve(getRsqt *service.RetrieveRequest) (int32) {
	result, err := server.Retrieve(context.Background(), getRsqt)
	if err != nil {
		log.Printf("Client %s get request failed: %s", connection.Target(), err)
		Redial()
		return Retrieve(getRsqt)
	}

	if reflect.ValueOf(result.Value).Kind() != reflect.ValueOf(int32(5)).Kind() {
		return 0
	}
	return result.Value
}

//In the case of losing connection - alternates between predefined front-
func Redial() {
	var port string
	if connection.Target()[len(connection.Target())-1:] == "1" {
		port =  ":4000"
	} else {
		port = ":4001"
	}

	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Client: Unable to connect to port %s: %v", port, err)
	}

	connection = conn
	server = service.NewServiceClient(connection) //creates a connection with an FE server
}


func setLogClient() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}