package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/messages"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/tree"
)

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
		context.Send(msg.RootNode, &tree.Traverse{Instructor: context.Self()})
		currentActor.OpenNodes += currentActor.OpenNodes
	case *tree.Data:
		currentActor.OpenNodes -= currentActor.OpenNodes
		for k, v := range msg.Values {
			currentActor.Values[k] = v
		}
		if msg.RightNode != nil {
			context.Send(msg.RightNode, &tree.Traverse{Instructor: context.Self()})
			currentActor.OpenNodes += currentActor.OpenNodes
		}
		if msg.LeftNode != nil {
			context.Send(msg.LeftNode, &tree.Traverse{Instructor: context.Self()})
			currentActor.OpenNodes += currentActor.OpenNodes
		}
		if currentActor.OpenNodes == 0 {
			result := make([]*messages.NodeData, 0)
			for k, v := range currentActor.Values {
				result = append(result, &messages.NodeData{Key: int32(k), Value: v})
			}
			context.Send(currentActor.Instructor, &messages.TraverseResult{Values: result})
		}
	}

}
