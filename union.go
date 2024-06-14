package msgpack

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	unionMap sync.Map
)

type Union struct {
	registerInt2Type sync.Map
	registerType2Int sync.Map
}

func (u *Union) addSubStruct(i int, t reflect.Type) {
	u.registerType2Int.Store(t, i)
	u.registerInt2Type.Store(i, t)
}

func (u *Union) encodeUnionValue(e *Encoder, v reflect.Value) error {
	if v.IsNil() {
		e.EncodeNil()
		return nil
	}
	err := e.EncodeArrayLen(2)
	if err != nil {
		return err
	}
	v = v.Elem()
	i, ok := u.registerType2Int.Load(v.Type())
	if !ok {
		return errors.New(fmt.Sprintf("unregister sub type:%s", v.Type().Name()))
	}
	err = e.EncodeInt(int64(i.(int)))
	if err != nil {
		return err
	}
	return e.EncodeValue(v)
}

func (u *Union) decoderUnionValue(d *Decoder, v reflect.Value) error {

	l, err := d.DecodeArrayLen()
	if err != nil {
		return err
	}
	if l == -1 || l == 0 {
		return nil
	}
	if l != 2 {
		return errors.New(fmt.Sprintf("invalid union array len:%d", l))
	}
	i, err := d.DecodeInt()
	if err != nil {
		return err
	}
	t, ok := u.registerInt2Type.Load(i)
	if !ok {
		return errors.New(fmt.Sprintf("unregister union type:%d", i))
	}
	subV := reflect.New(t.(reflect.Type))
	subV = subV.Elem()
	err = d.DecodeValue(subV)
	if err != nil {
		return err
	}
	v.Set(subV)
	return nil
}

func RegisterUnionInterface[T any]() *Union {
	var defaultValue *T
	valueType := reflect.TypeOf(defaultValue).Elem()
	u := &Union{}
	useU, loaded := unionMap.LoadOrStore(valueType, u)
	u = useU.(*Union)
	if loaded {
		return u
	}
	RegisterByTyp(valueType, u.encodeUnionValue, u.decoderUnionValue)
	return u
}

func RegisterUnionInterfaceSubStruct[T any, Sub any](subType int) {
	var defaultT T
	u, ok := unionMap.Load(reflect.TypeOf(defaultT))
	if !ok {
		u = RegisterUnionInterface[T]()
	}
	var defaultSub Sub
	u.(*Union).addSubStruct(subType, reflect.TypeOf(defaultSub))
}
