package tree

import "github.com/AsynkronIT/protoactor-go/actor"

type Add struct {
	instructor *actor.PID
	key        int
	value      string
}

type Find struct {
	instructor *actor.PID
	key        int
	value      string
	delete     bool
}

type Delete struct{}

type Traverse struct {
	instructor *actor.PID
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
