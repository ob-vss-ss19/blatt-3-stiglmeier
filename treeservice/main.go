package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/messages"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/tree"
	"sync"
)

type ServiceActor struct {
	nextId   int
	rootNode *actor.PID
	token    string
	id       int
}

var (
	flagBind = flag.String("bind", "localhost:8091", "local adress:port")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	remote.Start(*flagBind)

	remote.Register("treeservice", actor.PropsFromProducer(func() actor.Actor { return &ServiceActor{1, nil, "", 0} }))
	wg.Wait()
}

func (currentActor *ServiceActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.NewTree:
		if msg.MaxLeafSize < 1 {
			context.Respond(&messages.Failure{"Minimum leaf size is 1!"})
			return
		} else if currentActor.token != "" {
			context.Respond(&messages.Failure{"Delete existing tree first!"})
			return
		}
		currentActor.id = currentActor.nextId
		currentActor.nextId++
		currentActor.token = newToken()

		currentActor.rootNode = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return &tree.NodeActor{LeafSize: int(msg.MaxLeafSize)}
		}))

		context.Respond(&messages.TreeData{Token: currentActor.token, Id: int32(currentActor.id)})
		fmt.Printf("responded to new tree with size %d...\n", msg.MaxLeafSize)
	case *messages.AddPair:
		if int32(currentActor.id) != msg.Id || currentActor.token != msg.Token {
			context.Respond(&messages.Failure{"Id or Token Mismatch!"})
			return
		}
		context.Send(currentActor.rootNode, &tree.Add{Instructor: context.Sender(), Key: int(msg.Key), Value: msg.Value})
		fmt.Printf("instructed tree to add\n")
	case *messages.RemovePair:
		if int32(currentActor.id) != msg.Id || currentActor.token != msg.Token {
			context.Respond(&messages.Failure{"Id or Token Mismatch!"})
			return
		}
		context.Send(currentActor.rootNode, &tree.Find{Instructor: context.Sender(), Key: int(msg.Key), Value: msg.Value, Delete: true})
		fmt.Printf("instructed tree to delete\n")
	case *messages.FindPair:
		if int32(currentActor.id) != msg.Id || currentActor.token != msg.Token {
			context.Respond(&messages.Failure{"Id or Token Mismatch!"})
			return
		}
		context.Send(currentActor.rootNode, &tree.Find{Instructor: context.Sender(), Key: int(msg.Key), Value: msg.Value, Delete: false})
		fmt.Printf("instructed tree to find\n")
	case *messages.DeleteTree:
		if int32(currentActor.id) != msg.Id || currentActor.token != msg.Token {
			context.Respond(&messages.Failure{"Id or Token Mismatch!"})
			return
		}
		context.Send(currentActor.rootNode, &tree.Delete{})
		currentActor.rootNode = nil
		currentActor.token = ""
		currentActor.id = 0
		context.Respond(&messages.Success{})
		fmt.Printf("instructed tree to delete itself\n")
	case *messages.TraverseTree:
		if int32(currentActor.id) != msg.Id || currentActor.token != msg.Token {
			context.Respond(&messages.Failure{"Id or Token Mismatch!"})
			return
		}
		context.Send(currentActor.rootNode, &tree.Traverse{Instructor: context.Sender()})
		fmt.Printf("instructed tree to traverse\n")
	}
}

func newToken() string {
	b := make([]byte, 2)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
