package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-stiglmeier/messages"
	"sort"
)

type Add struct {
	Instructor *actor.PID
	Value      string
	Key        int
}

type Find struct {
	Instructor *actor.PID
	Key        int
	Value      string
	Delete     bool
}

type Delete struct{}

type Traverse struct {
	Instructor *actor.PID
}

type NodeActor struct {
	LeafSize   int
	LeftNode   *actor.PID
	RightNode  *actor.PID
	LeftMaxKey int
	Values     map[int]string
}

func (node *NodeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Add:
		fmt.Println("adding...")
		if node.LeftNode != nil { // no leaf
			if msg.Key <= node.LeftMaxKey {
				context.Send(node.LeftNode, msg)
			} else {
				context.Send(node.RightNode, msg)
			}
		} else if len(node.Values) < node.LeafSize { // free map
			node.Values[msg.Key] = msg.Value
			fmt.Println("added key: %d, value: %s", msg.Key, msg.Value)
		} else { // no free map -> split
			fmt.Println("splitted")
			node.Values[msg.Key] = msg.Value

			node.LeftNode = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{LeafSize: int(node.LeafSize)}
			}))
			node.RightNode = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{LeafSize: int(node.LeafSize)}
			}))
			sortedKeys := sortNode(node.Values)
			node.LeftMaxKey = sortedKeys[(len(sortedKeys)/2)-1]
			for k, v := range node.Values {
				if k <= sortedKeys[(len(sortedKeys)/2)-1] {
					context.Send(node.LeftNode, Add{Instructor: context.Self(), Key: k, Value: v})
				} else {
					context.Send(node.LeftNode, Add{Instructor: context.Self(), Key: k, Value: v})
				}
			}
			node.Values = nil
		}
		context.Send(msg.Instructor, &messages.Success{})
	default:
		fmt.Printf("invalid message type for node: %s\n", msg)
	}

}

func sortNode(Values map[int]string) []int {
	var keys []int
	for k := range Values {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}
