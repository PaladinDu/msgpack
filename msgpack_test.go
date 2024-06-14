package msgpack

import (
	"fmt"
	"reflect"
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

type ITestInterface interface {
	ShowValue()
}

type AttrStructA struct {
	Value int `msgpack:"index:0"`
}

func (a AttrStructA) ShowValue() {
}

func NewAttrStructA(value int) ITestInterface {
	ret := &AttrStructA{
		Value: value,
	}
	return ret
}

type AttrStructB struct {
	Value string `msgpack:"index:0"`
}

func (a *AttrStructB) ShowValue() {
}

func NewAttrStructB(value string) ITestInterface {
	ret := &AttrStructB{
		Value: value,
	}
	return ret
}
func checkInterface() ITestInterface {
	return &AttrStructB{}
}

type StructData struct {
	Attr1 ITestInterface `msgpack:"index:0"`
	Attr2 ITestInterface `msgpack:"index:1,Discriminator"`
}

func showMsgPackObj(obj any) {
	bytes, err := Marshal(obj)
	if err != nil {
		panic(err)
	}
	var showObj any
	err = Unmarshal(bytes, &showObj)
	if err != nil {
		panic(err)
	}
	fmt.Println(showObj)
}

func ExampleUnmarshal_subStruct() {
	RegisterUnionInterface[ITestInterface]()
	RegisterUnionInterfaceSubStruct[ITestInterface, *AttrStructA](1)
	RegisterUnionInterfaceSubStruct[ITestInterface, *AttrStructB](2)
	obj := &StructData{
		Attr1: NewAttrStructA(1),
		Attr2: NewAttrStructB("2"),
	}
	bytes, err := Marshal(obj)
	if err != nil {
		panic(err)
	}
	var showObj any
	err = Unmarshal(bytes, &showObj)
	if err != nil {
		panic(err)
	}
	fmt.Println(showObj)
	decodeObj := &StructData{}
	err = Unmarshal(bytes, decodeObj)
	if err != nil {
		panic(err)
	}
	showMsgPackObj(decodeObj)
	// Output:
	// [[1 [1]] [2 [2]]]
	// [[1 [1]] [2 [2]]]
}

type Enum int
type FIntP1 int
type A struct {
	Type  Enum   `msgpack:"index:0"`
	Value FIntP1 `msgpack:"index:1,precision:1"`
}

func ExampleUnmarshal_SelfType() {
	var dv FIntP1
	RegisterByTyp(reflect.TypeOf(dv), func(encoder *Encoder, value reflect.Value) error {
		return encoder.EncodeFloat32(float32(value.Int()) / 10)
	}, func(decoder *Decoder, value reflect.Value) error {
		f, err := decoder.DecodeFloat32()
		if err != nil {
			return err
		}
		value.SetInt(int64(f*10 + 0.1))
		return nil
	})
	a := &A{
		Type:  1,
		Value: 2,
	}
	data, err := Marshal(a)
	if err != nil {
		panic(err)
	}
	var obj any
	err = Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	showMsgPackObj(obj)
	b := &A{}
	err = Unmarshal(data, b)
	fmt.Println(b.Type, b.Value)
	// Output:
	//[1 0.2]
	//1 2

}
