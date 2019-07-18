package tree

import (
	"fmt"
	"testing"
	"time"
	"unsafe"

	"math/rand"
)

const Record = 10000

func Test_Reader(t *testing.T) {
	r := newTree()
	testTreeMemGet(r)
	testTreeMemSeek(r)

	byteArray, _ := NewWriter(r).Encode()
	fmt.Println()
	fmt.Println("file-size:", len(byteArray))

	reader := NewReader(byteArray)

	testTreeReaderGet(reader)
	testTreeReaderSeek(reader)
}

func Test_SeekToFirst(t *testing.T) {
	r := newTree()

	byteArray, _ := NewWriter(r).Encode()
	reader := NewReader(byteArray)

	it := reader.SeekToFirst()
	success := 0
	for it.Next() {
		success++
	}
	fmt.Println("success: ", success)
}

func Test_Range(t *testing.T) {
	r := newTree()

	byteArray, _ := NewWriter(r).Encode()
	reader := NewReader(byteArray)

	it := reader.Range([]byte("key-1000000"), []byte("key-200000"))
	success := 0
	for it.Next() {
		success++
	}
	fmt.Println("success: ", success)
}

func Test_Random(t *testing.T) {
	r := NewBTree()
	var records = make(map[string]int)
	for i := 0; i < Record; i++ {
		key := RandStringBytesMaskImprSrc(8)
		records[key] = i
		r.Put([]byte(key), i)
	}

	by, _ := NewWriter(r).Encode()
	reader := NewReader(by)

	startTime := time.Now().UnixNano()
	success := 0
	for k, v := range records {
		value, _ := reader.Get([]byte(k))
		if value == v {
			success++
		} else {
			fmt.Println("xx", k)
		}
	}
	fmt.Println("file-get:", (time.Now().UnixNano()-startTime)/1000000, "ms")
	fmt.Println("file-success:", success)
}

func testTreeReaderSeek(reader *Reader) {
	success := 0
	startTime := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		success = 0
		key := fmt.Sprintf("%s%d", "key-", i)
		it := reader.Seek([]byte(key))
		if nil != it {
			for it.Next() {
				success++
				//fmt.Println("k:", string(it.GetKey()), " v:", it.GetValue())
			}
		}
		//fmt.Println(i, " seek-success:", success)
	}
	fmt.Println("file-seek:", (time.Now().UnixNano()-startTime)/1000000, "ms")
}

func newTree() *BTree {
	t := NewBTree()

	for i := 0; i < Record; i++ {
		t.Put([]byte(fmt.Sprintf("%s%d", "key-", i)), i)
	}
	return t
}

func testTreeMemGet(t *BTree) {
	success := 0
	startTime := time.Now().UnixNano()
	for i := 0; i < Record; i++ {
		key := fmt.Sprintf("%s%d", "key-", i)
		v, _ := t.Get([]byte(key))
		if i == v {
			success++
		} else {
			fmt.Println("error...", i)
		}
	}
	fmt.Println("mem-get:", (time.Now().UnixNano()-startTime)/1000000, "ms")
	fmt.Println("mem-success:", success)
}

func testTreeMemSeek(t *BTree) {
	success := 0
	startTime := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		success = 0
		key := fmt.Sprintf("%s%d", "key-", i)
		e, ok := t.tree.Seek([]byte(key))
		if ok {
			for {
				_, _, err := e.Next()
				if nil == err {
					success++
				} else {
					break
				}
			}
		}
		//fmt.Println(i, " seek-success:", success)
	}
	fmt.Println("mem-seek:", (time.Now().UnixNano()-startTime)/1000000, "ms")
}

func testTreeReaderGet(r *Reader) {
	success := 0
	startTime := time.Now().UnixNano()
	for i := 0; i < Record; i++ {
		key := fmt.Sprintf("%s%d", "key-", i)
		v, _ := r.Get([]byte(key))
		if i == v {
			success++
		} else {
			fmt.Println("error...", i)
		}
	}
	fmt.Println("file-get:", (time.Now().UnixNano()-startTime)/1000000, "ms")
	fmt.Println("file-success:", success)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxMask = 1<<6 - 1 // All 1-bits, as many as 6
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for 10 characters!
	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		b[i] = letterBytes[int(cache&letterIdxMask)%len(letterBytes)]
		i--
		cache >>= 6
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
