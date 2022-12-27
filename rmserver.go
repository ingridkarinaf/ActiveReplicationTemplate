package main

import (
	"fmt"
	service "github.com/ingridkarinaf/ActiveReplicationTemplate/interface"
	grpc "google.golang.org/grpc"
	"strconv"
	"os"
	"context"
	"net"
	"log"
	"time"
)

/*
	- Responsible for maintaining copies of data for FEServer. 
	- Is subject to crashing.
*/

type RMServer struct {
	service.UnimplementedServiceServer
	id              int32 //portnumber, between 5000 and 5002
	ctx             context.Context
	data   			int32 //Update
	lockChannel 	chan bool
}

func main() {
	//log to file instead of console
	f := setLogRMServer()
	defer f.Close()

	portInput, _ := strconv.ParseInt(os.Args[1], 10, 32) //Takes arguments 5000, 5001 and 5002
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rmServer := &RMServer{
		id:              int32(portInput),
		data:   		0, //update
		ctx:             ctx,
		lockChannel: 	make(chan bool, 1),
	}

	//Unlock channel
	rmServer.lockChannel <- true

	list, err := net.Listen("tcp", fmt.Sprintf(":%v", portInput))
	if err != nil {
		log.Fatalf("RM Server %v: Failed to listen on port: %v", rmServer.id, err)
	}

	grpcServer := grpc.NewServer()
	service.RegisterServiceServer(grpcServer, rmServer)
	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("RM server %v failed to serve %v", rmServer.id, err)
		}
	}()

	for {}
}

func (RM *RMServer) Update(ctx context.Context, hashUpt *service.UpdateRequest) (*service.UpdateReply, error){
	<- RM.lockChannel 
	time.Sleep(5 * time.Second)
	//update this whole function
	//update: depending on requirements, update data before saving to message
	returnMessage := &service.UpdateReply{
		CurrentValue: int32(RM.data),
		Outcome: true,
	}
	fmt.Println("RM data: ", RM.data)
	RM.data++
	fmt.Println("RM data: ", RM.data)
	RM.lockChannel <- true
	return returnMessage, nil
}


func (RM *RMServer) Retrieve(ctx context.Context, getRqst *service.RetrieveRequest) (*service.RetrieveReply, error) {
	//update this whole function
	value := RM.data
	getResp := &service.RetrieveReply{
		Value:  value,
	}
	return getResp, nil
}

func setLogRMServer() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}