package bson

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"io"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func randomFill(depth int, doc *D) {
	n := 11
	for range make([]struct{}, n) {
		if depth <= 0 {
			return
		}
		seed := time.Now().UnixNano()
		rand.Seed(seed)
		switch rand.Intn(n) {
		case 1:
			var _embed D
			randomFill(depth-1, &_embed)
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: _embed})
		case 2:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: int64(seed)})
		case 3:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: -int32(seed)})
		case 4:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: strconv.FormatInt(seed, 16)})
		case 5:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: -int(seed)})
		case 6:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: nil})
		case 7:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: seed%2 == 0})
		case 8:
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: time.UnixMilli(seed)})
		case 9:
			rand.Seed(time.Now().UnixNano())
			b := make([]byte, rand.Intn(12))
			_, err := io.ReadFull(crand.Reader, b)
			if err != nil {
				panic(err)
			}
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: b})
		case 10:
			var b [12]byte
			_, err := io.ReadFull(crand.Reader, b[:])
			if err != nil {
				panic(err)
			}
			*doc = append(*doc, E{Name: strconv.Itoa(int(seed)), Value: b})
		}
	}
}
func TestCorrectness(t *testing.T) {
	for range make([]struct{}, 10) {
		var doc D
		randomFill(11, &doc)
		b, err := doc.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		doc2, err := ReadDocument(bufio.NewReader(bytes.NewReader(b)))
		if err != nil {
			t.Fatal(err)
		}
		b2, err := doc2.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b2, b) {
			t.FailNow()
		}
		if doc.String() != doc.String() {
			t.FailNow()
		}
	}
}
