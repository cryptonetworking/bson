package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

func writeCString(dst *bufio.Writer, str string) error {
	if len(str) >= math.MaxInt32 {
		return errors.New("bson: too large cString")
	}
	_, err := dst.WriteString(str)
	if err != nil {
		return err
	}
	return dst.WriteByte(0)
}
func WriteInt32(dst *bufio.Writer, n int32) error {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n))
	_, err := dst.Write(b)
	return err
}
func WriteInt64(dst *bufio.Writer, n int64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))
	_, err := dst.Write(b)
	return err
}

func WriteDocument(dst *bufio.Writer, doc D) error {
	b := bytes.NewBuffer(nil)
	buf := bufio.NewWriter(b)
	for _, elem := range doc {
		err := WriteElement(buf, elem.Name, elem.Value)
		if err != nil {
			return err
		}
	}
	err := buf.Flush()
	if err != nil {
		return err
	}
	if b.Len() >= math.MaxInt32 {
		return errors.New("bson: too large document")
	}
	err = WriteInt32(dst, int32(b.Len()))
	if err != nil {
		return err
	}
	_, err = b.WriteTo(dst)
	if err != nil {
		return err
	}
	return dst.WriteByte(0)
}
func WriteArray(dst *bufio.Writer, arr ...any) error {
	b := bytes.NewBuffer(nil)
	buf := bufio.NewWriter(b)
	for i, val := range arr {
		err := WriteElement(buf, strconv.Itoa(i), val)
		if err != nil {
			return err
		}
	}
	if b.Len() >= math.MaxInt32 {
		return errors.New("bson: too large document")
	}
	err := WriteInt32(dst, int32(b.Len()))
	if err != nil {
		return err
	}
	_, err = b.WriteTo(dst)
	if err != nil {
		return err
	}
	return dst.WriteByte(0)
}
func WriteObjectID(dst *bufio.Writer, b [12]byte) error {
	_, err := dst.Write(b[:])
	return err
}
func WriteElement(dst *bufio.Writer, name string, val any) error {
	var err error
	switch val := val.(type) {
	case string:
		err = dst.WriteByte(0x02)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteString(dst, val)
	case D: //TODO
		err = dst.WriteByte(0x03)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteDocument(dst, val)
	case []any:
		err = dst.WriteByte(0x04)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteArray(dst, val...)
	case []byte:
		err = dst.WriteByte(0x05)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteBinary(dst, val, SubtypeGenericBinary)
	case [12]byte:
		err = dst.WriteByte(0x07)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteObjectID(dst, val)
	case bool:
		err = dst.WriteByte(0x08)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteBool(dst, val)
	case time.Time:
		err = dst.WriteByte(0x09)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteInt64(dst, val.UnixMilli())
	case nil:
		err = dst.WriteByte(0x0A)
		if err != nil {
			return err
		}
		return writeCString(dst, name)
	case int32:
		err = dst.WriteByte(0x10)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteInt32(dst, val)
	case int64:
		err = dst.WriteByte(0x12)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteInt64(dst, val)
	case int:
		//nomalize
		err = dst.WriteByte(0x12)
		if err != nil {
			return err
		}
		err = writeCString(dst, name)
		if err != nil {
			return err
		}
		return WriteInt64(dst, int64(val))
	}
	return fmt.Errorf("bson: unsupported type %T", val)
}
func WriteString(dst *bufio.Writer, str string) error {
	if len(str) >= math.MaxInt32 {
		return errors.New("bson: too large string")
	}
	err := WriteInt32(dst, int32(len(str))+1)
	if err != nil {
		return err
	}
	return writeCString(dst, str)
}
func WriteBool(dst *bufio.Writer, b bool) error {
	if b {
		return dst.WriteByte(1)
	}
	return dst.WriteByte(0)
}
func WriteBinary(dst *bufio.Writer, bin []byte, subtype SubType) error {
	if len(bin) >= math.MaxInt32 {
		return errors.New("bson: too large binary")
	}
	err := WriteInt32(dst, int32(len(bin)))
	if err != nil {
		return err
	}
	err = dst.WriteByte(byte(subtype))
	if err != nil {
		return err
	}
	_, err = dst.Write(bin)
	return err
}
