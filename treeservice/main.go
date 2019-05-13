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

		props := actor.PropsFromProducer(func() actor.Actor {
			return &tree.NodeActor{LeafSize: int(msg.MaxLeafSize), LeftNode: nil, RightNode: nil, Values: make(map[int]string)}
		})
		currentActor.rootNode = context.Spawn(props)
		context.Respond(&messages.TreeData{Token: currentActor.token, Id: int32(currentActor.id)})
		fmt.Printf("responded to new tree with size %d...\n", msg.MaxLeafSize)
	case *messages.AddPair:
		if int32(currentActor.id) != msg.Id || currentActor.token != msg.Token {
			context.Respond(&messages.Failure{"Id or Token Mismatch!"})
			return
		}
		fmt.Printf("target pid: %s, %s\n", currentActor.rootNode.Id, currentActor.rootNode.Address)
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
		props := actor.PropsFromProducer(func() actor.Actor {
			return &TraverseActor{Instructor: context.Sender(), Values: make(map[int]string), OpenNodes: 0}
		})
		pid := context.Spawn(props)
		context.Send(pid, &StartTraverse{RootNode: currentActor.rootNode})
		fmt.Printf("instructed tree to traverse\n")
	}
}

func newToken() string {
	b := make([]byte, 2)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type StartTraverse struct {
	RootNode *actor.PID
}

type TraverseActor struct {
	Instructor *actor.PID
	Values     map[int]string
	OpenNodes  int
}

func (currentActor *TraverseActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *StartTraverse:
		fmt.Printf("Started Traversing from Actor\n")
		context.Send(msg.RootNode, &tree.Traverse{Instructor: context.Self()})
		currentActor.OpenNodes += 1
	case *tree.Data:
		fmt.Printf("Traverse Actor got Data...\n")
		currentActor.OpenNodes -= 1
		for k, v := range msg.Values {
			currentActor.Values[k] = v
		}
		if msg.RightNode != nil {
			fmt.Printf("Traverse Actor resending to right node...\n")
			context.Send(msg.RightNode, &tree.Traverse{Instructor: context.Self()})
			currentActor.OpenNodes += 1
		}
		if msg.LeftNode != nil {
			fmt.Printf("Traverse Actor resending to left node...\n")
			context.Send(msg.LeftNode, &tree.Traverse{Instructor: context.Self()})
			currentActor.OpenNodes += currentActor.OpenNodes
		}
		if currentActor.OpenNodes == 0 {
			fmt.Printf("Traverse Actor responding full data to cli...\n")
			result := make([]*messages.NodeData, 0)
			for k, v := range currentActor.Values {
				result = append(result, &messages.NodeData{Key: int32(k), Value: v})
			}
			context.Send(currentActor.Instructor, &messages.TraverseResult{Values: result})
		}
	}

}
