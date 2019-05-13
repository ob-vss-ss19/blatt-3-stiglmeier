package tree

import "github.com/AsynkronIT/protoactor-go/actor"

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
