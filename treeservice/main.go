package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/messages"
	"sync"
)

func main() {
	fmt.Println("Hello Tree-Service!")
	var wg sync.WaitGroup
	wg.Add(1)

	remote.Start("localhost:8091")

	remote.Register("hello", actor.PropsFromProducer(func() actor.Actor { return &MyActor{} }))
	wg.Wait()
}

type MyActor struct{}

func (*MyActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *messages.Echo:
		context.Respond(&messages.Response{SomeValue: "message received"})
		fmt.Printf("responding...")
	default:
		fmt.Printf("unknown message type\n")
	}
}
