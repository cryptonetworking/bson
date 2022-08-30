package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

func readCString(buf *bufio.Reader) (string, error) {
	str, err := buf.ReadString(0)
	if err != nil {
		return "", err
	}
	return str[:len(str)-1], nil
}

func ReadInt32(buf *bufio.Reader) (int32, error) {
	b := make([]byte, 4)
	_, err := io.ReadFull(buf, b)
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(b)), nil
}

func ReadInt64(buf *bufio.Reader) (int64, error) {
	b := make([]byte, 8)
	_, err := io.ReadFull(buf, b)
	if err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(b)), nil
}

func SkipZero(buf *bufio.Reader) error {
	b, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if b != 0 {
		return errors.New("bson: not zero")
	}
	return nil
}

func ReadDocument(buf *bufio.Reader) (doc D, err error) {
	_size, err := ReadInt32(buf)
	if err != nil {
		return doc, err
	}
	b := make([]byte, _size)
	_, err = io.ReadFull(buf, b)

	if err != nil {
		return doc, err
	}
	if err = SkipZero(buf); err != nil {
		return nil, err
	}
	buf = bufio.NewReader(bytes.NewReader(b))
	for {
		name, val, err := ReadElement(buf)
		if err != nil {
			if err == io.EOF {

				return doc, nil
			}
			return doc, err
		}
		doc = append(doc, E{name, val})
	}
}

func ReadElement(buf *bufio.Reader) (string, any, error) {
	_type, err := buf.ReadByte()
	if err != nil {
		return "", nil, err
	}
	name, err := readCString(buf)
	if err != nil {
		return "", nil, err
	}
	switch _type {
	// case 0x01: TODO
	//	n, err := ReadInt64(buf)
	//	return name, float64(n), err
	case 0x02:
		str, err := ReadString(buf)
		return name, str, err
	case 0x03:
		doc, err := ReadDocument(buf)
		return name, doc, err
	case 0x04:
		doc, err := ReadDocument(buf)
		if err != nil {
			return name, nil, err
		}
		arr, err := doc.ToArray()
		return name, arr, err
	case 0x05:
		bin, _, err := ReadBinary(buf)
		return name, bin, err
	case 0x06:
		return name, nil, nil
	case 0x07:
		oid, err := ReadObjectID(buf)
		return name, oid, err
	case 0x08:
		b, err := ReadBool(buf)
		return name, b, err
	case 0x09:
		n, err := ReadInt64(buf)
		return name, time.UnixMilli((n)), err
	case 0x0A:
		return name, nil, nil
	// case 0x0B: TODO regexp
	// case 0x0C: deprecated
	// case 0x0D: deprecated
	// case 0x0E: deprecated
	// case 0x0F: deprecated
	case 0x10:
		n, err := ReadInt32(buf)
		return name, n, err
	// case 0x11: mongodb internal
	case 0x12:
		n, err := ReadInt64(buf)
		return name, (n), err
		//
		// case 0x13: TODO decimal128
		// case 0xff: mongodb internal
		// case 0x7f: mongodb internal
	}
	return "", nil, fmt.Errorf("unknown element type %x", _type)
}

func ReadObjectID(buf *bufio.Reader) ([12]byte, error) {
	var b [12]byte
	_, err := io.ReadFull(buf, b[:])
	return b, err
}

func ReadBool(buf *bufio.Reader) (bool, error) {
	b, err := buf.ReadByte()
	return b != 0, err
}

func ReadBinary(buf *bufio.Reader) ([]byte, SubType, error) {
	n, err := ReadInt32(buf)
	if err != nil {
		return nil, 0, err
	}
	st, err := buf.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	b := make([]byte, n)
	_, err = io.ReadFull(buf, b)
	if err != nil {
		return nil, 0, err
	}
	return b, SubType(st), nil
}

func ReadString(buf *bufio.Reader) (string, error) {
	n, err := ReadInt32(buf)
	if err != nil {
		return "", err
	}
	b := make([]byte, n-1)
	_, err = io.ReadFull(buf, b)
	if err != nil {
		return "", err
	}
	if err = SkipZero(buf); err != nil {
		return "", err
	}
	return string(b), nil
}
