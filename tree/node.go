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

type Data struct {
	LeftNode  *actor.PID
	RightNode *actor.PID
	Values    map[int]string
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
				return &NodeActor{LeafSize: int(node.LeafSize), LeftNode: nil, RightNode: nil, Values: make(map[int]string)}
			}))
			node.RightNode = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{LeafSize: int(node.LeafSize), LeftNode: nil, RightNode: nil, Values: make(map[int]string)}
			}))
			sortedKeys := sortNode(node.Values)
			node.LeftMaxKey = sortedKeys[(len(sortedKeys)/2)-1]
			for k, v := range node.Values {
				if k <= sortedKeys[(len(sortedKeys)/2)-1] {
					context.Send(node.LeftNode, &Add{Instructor: context.Self(), Key: k, Value: v})
				} else {
					context.Send(node.RightNode, &Add{Instructor: context.Self(), Key: k, Value: v})
				}
			}
			node.Values = nil
		}
		context.Send(msg.Instructor, &messages.Success{})
	case *Find:
		fmt.Println("finding...")
		if node.LeftNode != nil {
			if node.LeftMaxKey <= msg.Key {
				context.Send(node.LeftNode, msg)
			} else {
				context.Send(node.RightNode, msg)
			}
		} else {
			for k, v := range node.Values {
				if k == msg.Key && v == msg.Value {
					if msg.Delete {
						delete(node.Values, k)
					}
					context.Send(msg.Instructor, &messages.Success{})
					return
				}
			}
			context.Send(msg.Instructor, &messages.Failure{Cause: "Node not found"})
		}
	case *Delete:
		fmt.Println("deleting...")
		node.Values = nil
		if node.LeftNode != nil {
			context.Send(node.LeftNode, msg)
			context.Send(node.RightNode, msg)
		}
		node.LeftNode = nil
		node.RightNode = nil
	case *Traverse:
		fmt.Println("traversing...")
		context.Send(msg.Instructor, &Data{LeftNode: node.LeftNode, RightNode: node.RightNode, Values: node.Values})
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
