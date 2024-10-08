package session

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *data) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Data":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
			if z.Data == nil {
				z.Data = make(map[string]interface{}, zb0002)
			} else if len(z.Data) > 0 {
				for key := range z.Data {
					delete(z.Data, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 interface{}
				za0001, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Data")
					return
				}
				za0002, err = dc.ReadIntf()
				if err != nil {
					err = msgp.WrapError(err, "Data", za0001)
					return
				}
				z.Data[za0001] = za0002
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *data) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Data"
	err = en.Append(0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Data)))
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	for za0001, za0002 := range z.Data {
		err = en.WriteString(za0001)
		if err != nil {
			err = msgp.WrapError(err, "Data")
			return
		}
		err = en.WriteIntf(za0002)
		if err != nil {
			err = msgp.WrapError(err, "Data", za0001)
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *data) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Data"
	o = append(o, 0x81, 0xa4, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for za0001, za0002 := range z.Data {
		o = msgp.AppendString(o, za0001)
		o, err = msgp.AppendIntf(o, za0002)
		if err != nil {
			err = msgp.WrapError(err, "Data", za0001)
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *data) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Data":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
			if z.Data == nil {
				z.Data = make(map[string]interface{}, zb0002)
			} else if len(z.Data) > 0 {
				for key := range z.Data {
					delete(z.Data, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 interface{}
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Data")
					return
				}
				za0002, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Data", za0001)
					return
				}
				z.Data[za0001] = za0002
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *data) Msgsize() (s int) {
	s = 1 + 5 + msgp.MapHeaderSize
	if z.Data != nil {
		for za0001, za0002 := range z.Data {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + msgp.GuessSize(za0002)
		}
	}
	return
}
