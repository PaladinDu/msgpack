package msgpack

import (
	"fmt"
)

type Base struct {
	StrID  string    `msgpack:"index:2"`
	ID     int       `msgpack:"index:0"`
	Parent *DataNode `msgpack:"index:1"`
}
type DataNode struct {
	Base
	V1 string      `msgpack:"index:3"`
	V2 int         `msgpack:"index:4"`
	V3 map[int]int `msgpack:"index:5"`
	V4 bool        `msgpack:"index:6"`
}

type SimpleDataNode struct {
	Base
	Value2 int         `msgpack:"index:4"`
	Value3 map[int]int `msgpack:"index:5"`
}

func ExampleMarshal() {
	node := &DataNode{
		Base: Base{
			StrID: "Node1",
			ID:    1,
		},
		V1: "V1",
		V2: 2,
		V3: map[int]int{
			3: 4,
		},
		V4: true,
	}
	bytes, err := Marshal(node)
	if err != nil {
		panic(err)
	}
	var showObj any
	err = Unmarshal(bytes, &showObj)
	if err != nil {
		panic(err)
	}
	fmt.Println(showObj)
	// Output: [1 <nil> Node1 V1 2 map[3:4] true]
}

func ExampleUnmarshal() {
	node := &DataNode{
		Base: Base{
			ID:    1,
			StrID: "Node1",
		},
		V1: "V1",
		V2: 2,
		V3: map[int]int{
			3: 4,
		},
		V4: true,
	}
	bytes, err := Marshal(node)
	if err != nil {
		panic(err)
	}
	simpleNode := &SimpleDataNode{}
	err = Unmarshal(bytes, simpleNode)
	if err != nil {
		panic(err)
	}
	fmt.Println(simpleNode.ID, simpleNode.StrID, simpleNode.Value2, simpleNode.Value3)
	bytes, err = Marshal(simpleNode)
	if err != nil {
		panic(err)
	}
	node = &DataNode{}
	err = Unmarshal(bytes, node)
	fmt.Println(node.ID, node.StrID, node.V1, node.V2, node.V3, node.V4)
	var showObj any
	err = Unmarshal(bytes, &showObj)
	if err != nil {
		panic(err)
	}
	fmt.Println(showObj)
	// Output:
	// 1 Node1 2 map[3:4]
	// 1 Node1  2 map[3:4] false
	// [1 <nil> Node1 <nil> 2 map[3:4]]
}
