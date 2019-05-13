package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
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
	switch context.Message().(type) {
	case Add:
		if node.LeftNode != nil { // no leaf
			//if(msg.)
			//context.Send(node)
		}
	}

}
