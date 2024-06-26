package msgpack

import (
	"fmt"
)

type Marshaler interface {
	MarshalMsgpack() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalMsgpack([]byte) error
}

type CustomEncoder interface {
	EncodeMsgpack(*Encoder) error
}

type CustomDecoder interface {
	DecodeMsgpack(*Decoder) error
}

//------------------------------------------------------------------------------

type RawMessage []byte

var (
	_ CustomEncoder = (RawMessage)(nil)
	_ CustomDecoder = (*RawMessage)(nil)
)

func (m RawMessage) EncodeMsgpack(enc *Encoder) error {
	return enc.write(m)
}

func (m *RawMessage) DecodeMsgpack(dec *Decoder) error {
	msg, err := dec.DecodeRaw()
	if err != nil {
		return err
	}
	*m = msg
	return nil
}

//------------------------------------------------------------------------------

type unexpectedCodeError struct {
	hint string
	code byte
}

func (err unexpectedCodeError) Error() string {
	return fmt.Sprintf("msgpack: unexpected code=%x decoding %s", err.code, err.hint)
}

func StructToAnyStruct(obj any) (any, error) {
	var data any
	bytes, err := Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func CopyAnyT[T any](src T) (T, error) {
	var v T
	srcData, err := Marshal(src)
	if err != nil {
		return v, err
	}
	err = Unmarshal(srcData, &v)
	if err != nil {
		return v, err
	}
	return v, nil
}
