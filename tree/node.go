package tree

import "github.com/AsynkronIT/protoactor-go/actor"

type Add struct {
	Instructor *actor.PID
	Key        int
	Value      string
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

func (state *NodeActor) Receive(context actor.Context) {
	//switch msg := context.Message().(type) {
	//case :

}
