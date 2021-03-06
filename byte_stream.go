package goao

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ByteStream struct {
	realWrite  bool
	iOffset    int
	byPackage  []byte
	dwLength   int
	bGood      bool
	iBufLength int
}

func NewByteStream() *ByteStream {
	bs := &ByteStream{}
	bs.bGood = true
	bs.realWrite = true
	bs.iOffset = 0
	return bs
}

func (b *ByteStream) SetRealWrite(realWrite bool) {
	b.realWrite = realWrite
}

func (b *ByteStream) GetWrittenLength() int {
	return b.iOffset
}

func (b *ByteStream) PushByte(data byte) {

	if !b.realWrite {
		b.iOffset++
		return
	}

	b.byPackage = append(b.byPackage, data)

	b.iOffset++

}

func (b *ByteStream) PopByte() (byte, error) {

	if !b.bGood || (b.iOffset+1 > b.iBufLength) {
		b.bGood = false
		return 0, errors.New("popByte ERROR")
	}

	data := b.byPackage[b.iOffset]

	b.iOffset++

	return data, nil
}

func (b *ByteStream) PushUint64(data uint64) {

	if !b.realWrite {
		b.iOffset += 17
		return
	}
	byteData := fmt.Sprintf("%016x", data)

	b.PushBytes([]byte(byteData), 16)

	b.PushByte(0)
}

func (b *ByteStream) PopUint64() (uint64, error) {

	byteData, err := b.PopBytes(16)

	if err != nil {
		return 0, err
	}
	b.PopByte()
	//fmt.Println("uint:", byteData, ",offset", b.iOffset)

	strData := strings.TrimLeft(string(byteData), "0")

	if strData == "" {
		return 0, nil
	}

	data, err := strconv.ParseUint(strData, 16, 64)

	return data, err
}

func (b *ByteStream) PopInt64() (int64, error) {

	byteData, err := b.PopBytes(16)

	if err != nil {
		return 0, err
	}
	b.PopByte()

	strData := strings.TrimLeft(string(byteData), "0")

	if strData == "" {
		return 0, nil
	}

	data, errParse := strconv.ParseInt(strData, 16, 64)

	return data, errParse

}

func (b *ByteStream) PushUint64Set(data []uint64) {
	length := len(data)

	b.PushUint32(uint32(length))

	for item := range data {
		b.PushUint64(data[item])
	}
}

func (b *ByteStream) PopUint64Set() ([]uint64, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	var data []uint64

	for i := uint32(0); i < size; i++ {
		var d uint64
		d, err = b.PopUint64()
		if err != nil {
			return data, err
		}
		data = append(data, d)
	}

	return data, nil
}

func (b *ByteStream) PushUint32(data uint32) {

	if !b.realWrite {
		b.iOffset += 4
		return
	}

	b.byPackage = append(b.byPackage, byte(data>>24))

	b.byPackage = append(b.byPackage, byte(data>>16))

	b.byPackage = append(b.byPackage, byte(data>>8))

	b.byPackage = append(b.byPackage, byte(data))

	b.iOffset += 4

}

func (b *ByteStream) PushInt32(data int32) {

	if !b.realWrite {
		b.iOffset += 4
		return
	}

	b.byPackage = append(b.byPackage, byte(data>>24))

	b.byPackage = append(b.byPackage, byte(data>>16))

	b.byPackage = append(b.byPackage, byte(data>>8))

	b.byPackage = append(b.byPackage, byte(data))

	b.iOffset += 4

}

/*
func (b *ByteStream) PushUint32FromInt64(data64 int64) {

	if !b.realWrite {
		b.iOffset += 4
		return
	}
	data := uint32(data64)

	b.byPackage = append(b.byPackage, byte(data>>24))

	b.byPackage = append(b.byPackage, byte(data>>16))

	b.byPackage = append(b.byPackage, byte(data>>8))

	b.byPackage = append(b.byPackage, byte(data))

	b.iOffset += 4
}*/

func (b *ByteStream) PopUint() (int, error) {
	if !b.bGood || (b.iOffset+4 > b.iBufLength) {
		b.bGood = false
		return 0, errors.New("PopUint ERROR")
	}

	dataBytes := b.byPackage[b.iOffset : b.iOffset+4]

	b.iOffset += 4

	strData := strings.TrimLeft(string(dataBytes), "0")

	if strData == "" {
		return 0, nil
	}
	data, errParse := strconv.ParseInt(strData, 16, 10)

	if errParse != nil {
		return 0, errParse
	}

	return int(data), nil
}
func (b *ByteStream) PopUint32() (uint32, error) {

	//tData := NewUint16t(data)

	if !b.bGood || (b.iOffset+4 > b.iBufLength) {
		b.bGood = false
		return 0, errors.New("PopUint32 ERROR")
	}

	dataBytes := b.byPackage[b.iOffset : b.iOffset+4]

	b.iOffset += 4

	data := binary.BigEndian.Uint32(dataBytes)

	return data, nil
}

func (b *ByteStream) PopInt32() (int32, error) {

	//tData := NewUint16t(data)

	if !b.bGood || (b.iOffset+4 > b.iBufLength) {
		b.bGood = false
		return 0, errors.New("PopUint32 ERROR")
	}

	dataBytes := b.byPackage[b.iOffset : b.iOffset+4]

	b.iOffset += 4

	var data int32

	b_buf := bytes.NewBuffer(dataBytes)

	binary.Read(b_buf, binary.BigEndian, &data)

	return data, nil
}

func (b *ByteStream) PushUint16(data uint16) {

	//tData := NewUint16t(data)

	if !b.realWrite {
		b.iOffset += 2 //tData.GetSize()
		return
	}

	b.byPackage = append(b.byPackage, byte(data>>8))

	b.byPackage = append(b.byPackage, byte(data))

	b.iOffset += 2

}

func (b *ByteStream) PopUint16() (uint16, error) {

	//tData := NewUint16t(data)

	if !b.bGood || (b.iOffset+2 > b.iBufLength) {
		b.bGood = false
		return 0, errors.New("PopUint16 ERROR")
	}

	dataBytes := b.byPackage[b.iOffset : b.iOffset+2]

	b.iOffset += 2

	data := binary.BigEndian.Uint16(dataBytes)

	return data, nil
}

func (b *ByteStream) PushUint8(data uint8) {

	//tData := NewUint16t(data)

	b.PushByte(byte(data))
}

func (b *ByteStream) PopUint8() (uint8, error) {

	//tData := NewUint16t(data)

	byteData, err := b.PopByte()
	if err != nil {
		return 0, err
	}

	return uint8(byteData), nil
}

func (b *ByteStream) PushBytes(data []byte, length int) {

	lgh := len(data)

	if lgh == 0 {
		b.PushUint32(0)
		return
	}

	if length == 0 {
		b.PushUint32(uint32(lgh))
	}

	if !b.realWrite {

		b.iOffset += lgh
		return
	}

	b.byPackage = append(b.byPackage, data...)

	b.iOffset += lgh

}

func (b *ByteStream) PopBytes(length int) ([]byte, error) {

	if length < 0 || length > 2097152 {
		return nil, errors.New("popBytes Error: len must >=0 and <=2097152")
	}

	if !b.bGood || (b.iOffset+length > b.iBufLength) {
		b.bGood = false
		return nil, errors.New("PushBytes ERROR")
	}

	if length > 0 {

		data := b.byPackage[b.iOffset : b.iOffset+length]
		b.iOffset += length

		return data, nil
	}

	return nil, nil
}

func (b *ByteStream) PushString(data string) {

	/*if b.decodeCharset != "" && len(b.decodeCharset) > 2 {
		byteData := binary.BigEndian binary.BigEndian.Uint16(dataBytes)
	}
	*/
	b.PushBytes([]byte(data), 0)

}
func (b *ByteStream) PopString() (string, error) {
	if !b.bGood || (b.iOffset+4 > b.iBufLength) {
		b.bGood = false
		return "", errors.New(fmt.Sprintf("pop int ERROR: !m_bGood || (%d + 4) > %d", b.iOffset,b.iBufLength))
	}

	intLen, err := b.PopUint32()

	if err != nil {
		return "", err
	}

	if !b.bGood || (b.iOffset+int(intLen) > b.iBufLength) {
		b.bGood = false
		return "", errors.New(fmt.Sprintf("pop string ERROR: !m_bGood || (%d + %d) > %d", b.iOffset, intLen, b.iBufLength))
	}

	if intLen > 2097152 {
		b.bGood = false
		return "", errors.New("len > MAX_LENGTH")
	}

	if intLen > 0 {
		buf := b.byPackage[b.iOffset : b.iOffset+int(intLen)]
		b.iOffset += int(intLen)

		return string(buf), nil
	}
	return "", nil
}

func (b *ByteStream) PushStringSet(data []string) {
	length := len(data)

	b.PushUint32(uint32(length))

	for item := range data {
		b.PushString(data[item])
	}
}

func (b *ByteStream) PopStringSet() ([]string, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	var data []string

	for i := uint32(0); i < size; i++ {
		var d string
		d, err = b.PopString()
		if err != nil {
			return data, err
		}
		data = append(data, d)
	}

	return data, nil
}

func (b *ByteStream) PushMapXXOO(data map[interface{}]IAoXXOO) {

	length := len(data)

	b.PushUint32(uint32(length))

	for key, value := range data {
		b.PushObject(key)
		value.Serialize(b)
	}
}
func (b *ByteStream) PushMap(obj map[interface{}]interface{}) {

	b.PushUint32(uint32(len(obj)))

	for k, v := range obj {
		b.PushObject(k)
		b.PushObject(v)
	}
}

func (b *ByteStream) PopMap(keyType reflect.Type, valueType reflect.Type) (map[interface{}]interface{}, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	obj := make(map[interface{}]interface{})
	for i := uint32(0); i < size; i++ {

		key, errKey := b.PopObject(keyType)

		value, errValue := b.PopObject(valueType)

		if errKey != nil {
			return nil, errKey
		} else if errValue != nil {
			return nil, errValue
		}
		obj[key] = value
	}
	return obj, nil
}

func (b *ByteStream) PushMapUInt32String(obj map[uint32]string) {

	b.PushUint32(uint32(len(obj)))

	for k, v := range obj {
		b.PushUint32(k)
		b.PushString(v)
	}
}

func (b *ByteStream) PopMapUInt32String() (map[uint32]string, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	obj := make(map[uint32]string)
	for i := uint32(0); i < size; i++ {

		key, errKey := b.PopUint32()

		value, errValue := b.PopString()

		if errKey != nil {
			return nil, errKey
		} else if errValue != nil {
			return nil, errValue
		}
		obj[key] = value
	}
	return obj, nil
}

func (b *ByteStream) PushMapUInt32UInt32(obj map[uint32]uint32) {

	b.PushUint32(uint32(len(obj)))

	for k, v := range obj {
		b.PushUint32(k)
		b.PushUint32(v)
	}
}

func (b *ByteStream) PopMapUInt32UInt32() (map[uint32]uint32, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	obj := make(map[uint32]uint32)
	for i := uint32(0); i < size; i++ {

		key, errKey := b.PopUint32()

		value, errValue := b.PopUint32()

		if errKey != nil {
			return nil, errKey
		} else if errValue != nil {
			return nil, errValue
		}
		obj[key] = value
	}
	return obj, nil
}

func (b *ByteStream) PushMapStringString(obj map[string]string) {

	b.PushUint32(uint32(len(obj)))

	for k, v := range obj {
		b.PushString(k)
		b.PushString(v)
	}
}

func (b *ByteStream) PopMapStringString() (map[string]string, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	obj := make(map[string]string)
	for i := uint32(0); i < size; i++ {

		key, errKey := b.PopString()

		value, errValue := b.PopString()

		if errKey != nil {
			return nil, errKey
		} else if errValue != nil {
			return nil, errValue
		}
		obj[key] = value
	}
	return obj, nil
}

func (b *ByteStream) PushSet(obj []IAoXXOO) {
	b.PushUint32(uint32(len(obj)))
	for _, v := range obj {
		v.Serialize(b)
	}
}

func (b *ByteStream) PushBitSet(data []uint8) {
	length := len(data)

	b.PushUint32(uint32(length))

	for item := range data {
		b.PushUint8(data[item])
	}
}
func (b *ByteStream) PopBitSet() ([]uint8, error) {

	size, err := b.PopUint32()

	if err != nil {
		return nil, err
	}
	var data []uint8

	for i := uint32(0); i < size; i++ {
		var d uint8
		d, err = b.PopUint8()
		if err != nil {
			return data, err
		}
		data = append(data, d)
	}
	return data, nil
}

func (b *ByteStream) PushBool(data bool) {

	if !b.realWrite {
		b.iOffset += 1 //tData.GetSize()
		return
	}
	var bt byte
	if data {
		bt = 1
	} else {
		bt = 0
	}
	b.byPackage = append(b.byPackage, bt)
}

func (b *ByteStream) PopBool() (bool, error) {
	byteData, err := b.PopByte()
	if err != nil {
		return false, err
	}

	if byteData == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (b *ByteStream) PushObject(data interface{}) {
	switch data.(type) {
	case uint8:
		b.PushUint8(data.(uint8))
		break
	case uint16:
		b.PushUint16(data.(uint16))
		break
	case uint32:
		b.PushUint32(data.(uint32))
		break
	case uint64:
		b.PushUint64(data.(uint64))
		break
	case string:
		b.PushString(data.(string))
		break
	default:
		break
	}
}
func (b *ByteStream) PopObject(dataType reflect.Type) (data interface{}, err error) {

	switch dataType.Kind() {
	case reflect.Uint8:
		data, err = b.PopUint8()
		break
	case reflect.Uint16:
		data, err = b.PopUint16()
		break
	case reflect.Uint32:
		data, err = b.PopUint32()
		break
	case reflect.Uint64:
		data, err = b.PopUint64()
		break
	case reflect.String:
		data, err = b.PopString()
		break
	default:
		data = nil
		err = errors.New("unsupport type")
		break
	}
	return nil, err
}

func (b *ByteStream) IsGood() bool {
	return b.bGood
}

func (b *ByteStream) GetWriteBuffer() []byte {
	return b.byPackage
}

func (b *ByteStream) GetReadLength() int {
	return b.iOffset
}
func (b *ByteStream) Reset(data []byte, length int) {
	b.byPackage = data
	b.iBufLength = length
	b.iOffset = 0
	b.bGood = true
	b.realWrite = true
}
