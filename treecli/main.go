package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/messages"
	"strconv"
	"sync"
	"time"
)

type CliActor struct {
}

func (state *CliActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Success:
		fmt.Println("Success")
		wg.Done()
	case *messages.Failure:
		fmt.Printf("Failure: " + msg.Cause + "\n")
		wg.Done()
	case *messages.TraverseResult:
		for _, v := range msg.Values {
			fmt.Printf("Node[id=%d, value=%s]\n", v.Key, v.Value)
		}
		wg.Done()
	case *messages.TreeData:
		fmt.Printf("New Tree[id: %d, token: %s]", msg.Id, msg.Token)
		wg.Done()
	}
}

var (
	flagBind   = flag.String("bind", "localhost:8090", "local adress:port")
	flagRemote = flag.String("remote", "localhost:8091", "remote adress:port")

	id          = flag.Int("id", 0, "tree id")
	token       = flag.String("token", "", "tree token")
	forceDelete = flag.Bool("force-delete", false, "delete tree")

	context    *actor.RootContext
	pid        *actor.PID
	spawnedPID *actor.PID
	wg         sync.WaitGroup
)

func main() {
	wg.Add(1)
	flag.Parse()
	remote.SetLogLevel(log.ErrorLevel)
	remote.Start(*flagBind)

	context = actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor { return &CliActor{} })
	pid = context.Spawn(props)

	spawnResponse, _ := remote.SpawnNamed(*flagRemote, "serviceactor", "treeservice", time.Second)
	spawnedPID = spawnResponse.Pid

	switch flag.Args()[0] {
	case "newtree":
		newTree()
	case "deletetree":
		if !*forceDelete {
			fmt.Println("ERROR: Please add force-delete flag for tree deletion")
			return
		}
		deleteTree()
	case "insertnode":
		insertNode()
	case "deletenode":
		deletenode()
	case "existsnode":
		existsnode()
	case "traversetree":
		traverseTree()
	default:
		printDocumentation()
		return
	}

	wg.Wait()
}

func newTree() {
	if len(flag.Args()) < 2 {
		fmt.Printf("ERROR: Please specify the max leaf size as second argument.")
		wg.Done()
		return
	}
	leafSize, _ := strconv.Atoi(flag.Args()[1])
	message := &messages.NewTree{MaxLeafSize: int32(leafSize)}
	context.RequestWithCustomSender(spawnedPID, message, pid)
}

func deleteTree() {
	message := &messages.DeleteTree{Token: *token, Id: int32(*id)}
	context.RequestWithCustomSender(spawnedPID, message, pid)
}

func insertNode() {
	if len(flag.Args()) < 3 {
		fmt.Printf("ERROR: Please specify the key/value as second/third argument.")
		wg.Done()
		return
	}
	key, _ := strconv.Atoi(flag.Args()[1])
	message := &messages.AddPair{Token: *token, Id: int32(*id), Key: int32(key), Value: flag.Args()[2]}
	context.RequestWithCustomSender(spawnedPID, message, pid)
}

func deletenode() {
	if len(flag.Args()) < 3 {
		fmt.Printf("ERROR: Please specify the key/value as second/third argument.")
		wg.Done()
		return
	}
	key, _ := strconv.Atoi(flag.Args()[1])
	message := &messages.RemovePair{Token: *token, Id: int32(*id), Key: int32(key), Value: flag.Args()[2]}
	context.RequestWithCustomSender(spawnedPID, message, pid)
}

func existsnode() {
	if len(flag.Args()) < 3 {
		fmt.Printf("ERROR: Please specify the key/value as second/third argument.")
		wg.Done()
		return
	}
	key, _ := strconv.Atoi(flag.Args()[1])
	message := &messages.FindPair{Token: *token, Id: int32(*id), Key: int32(key), Value: flag.Args()[2]}
	context.RequestWithCustomSender(spawnedPID, message, pid)
}

func traverseTree() {
	message := &messages.TraverseTree{Token: *token, Id: int32(*id)}
	context.RequestWithCustomSender(spawnedPID, message, pid)
}

func printDocumentation() {
	fmt.Println("ERROR: Please use the following commands: ")
	fmt.Printf("newtree\ndeletetree\ninsertnode\ndeletenode\nexistsnode\ntraversetree\n\n")
}
