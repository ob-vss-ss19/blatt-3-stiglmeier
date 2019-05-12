package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/messages"
	"sync"
	"time"
)

type MyActor struct {
	count int
}

func (state *MyActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		fmt.Println(msg.SomeValue + string(state.count))
	default:
		fmt.Printf("unknown message type\n")
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	remote.Start("localhost:8090")

	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor { return &MyActor{} })
	pid := context.Spawn(props)
	message := &messages.Echo{Message: "hej", Sender: 55}

	//this is to spawn remote actor we want to communicate with
	spawnResponse, _ := remote.SpawnNamed("localhost:8091", "myactor", "hello", time.Second)

	// get spawned PID
	spawnedPID := spawnResponse.Pid
	context.RequestWithCustomSender(spawnedPID, message, pid)
	//context.Send(spawnedPID, message)

	wg.Wait()
}
