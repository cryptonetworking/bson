package bson

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type E struct {
	Name  string
	Value any
}

type D []E

func (doc D) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	buf := bufio.NewWriter(b)
	err := WriteDocument(buf, doc)
	if err != nil {
		return nil, err
	}
	err = buf.Flush()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
func (doc D) Normalize() {
	for i, elem := range doc {
		switch val := elem.Value.(type) {
		case int:
			doc[i].Value = int64(val)
		}
	}
}

func FromMap(_map map[string]any) D {
	doc := make(D, 0, len(_map))
	for name, val := range _map {
		doc = append(doc, E{Name: name, Value: val})
	}
	return doc
}
func FromArray(values ...any) D {
	doc := make(D, len(values))
	for i, val := range values {
		doc[i].Name = strconv.Itoa(i)
		doc[i].Value = val
	}
	return doc
}

func (doc D) Get(name string) (any, int) {
	if len(doc) == 0 {
		return nil, -1
	}
	for i, elem := range doc {
		if elem.Name == name {
			return elem.Value, i
		}
	}
	return nil, -1
}

func (doc D) ToArray() ([]any, error) {
	if len(doc) == 0 {
		return nil, nil
	}
	arr := make([]any, len(doc))
	for i := range arr {
		val, index := doc.Get(strconv.Itoa(i))
		if index < 0 {
			return nil, errors.New("bson: invalid array")
		}
		arr[i] = val
	}
	return arr, nil
}
func (doc D) ToMap() map[string]any {
	_map := make(map[string]any)
	for _, elem := range doc {
		if doc, ok := elem.Value.(D); ok {
			_map[elem.Name] = doc.ToMap()
		} else {
			_map[elem.Name] = elem.Value
		}
	}
	return _map
}

func (doc D) String() string {
	return fmt.Sprintf("%+v", doc.ToMap())
}
