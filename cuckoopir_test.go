package cuckoopir

import (
	"fmt"
	"testing"
	_ "reflect"
	_ "runtime"
	"sync"
	"math/rand"
	"time"
	"runtime"
	_ "runtime/debug"
)


const LOGQ = uint64(32)
const SEC_PARAM = uint64(1 << 10)

var rows = []uint64{1}
// var rows = []uint64{11,4,5,14,19,19,8,10}
var ell uint64 = uint64(len(rows))


func TestPIR(t *testing.T) {
	// fmt.Println("Number of CPUs:", runtime.NumCPU())

    // runtime.GOMAXPROCS(runtime.NumCPU())

	N := uint64(1 << 24)
	// Num        uint64 // number of DB entries.
	d := uint64(8)
	// Row_length uint64 // number of bits per DB entry.
	pir := CuckooPIR{}
	// p := pir.PickParams(N, d, SEC_PARAM, LOGQ)//return Params
	// p := Params{1<<2,6.4,1<<2,1<<3,32,1<<8}//toy example
	p := Params{993,12.8,1<<14,1<<14,32,1<<8}//return Params
	// p := Params{1<<10,6.4,1<<12,1<<12,32,1<<8}//return Params

	DB := MakeRandomDB(N, d, &p)//return *Database

	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 1; i++ {
		// go RunPIR(&pir, DB, p, []uint64{1,2,3},&wg)
		go RunPIR(&pir, DB, p, rows, &wg)
	}
	wg.Wait()
	fmt.Println("Done")
}

// func TestCuckoo(t *testing.T) {
// 	fmt.Println(n)
// 	fmt.Println("gmap:", gmap)
// 	// fmt.Println("gkeys:", gkeys)
// 	// fmt.Println("gvals:", gvals)
// 	fmt.Println("mapBytes:", mapBytes)

// 	c := NewCuckoo(DefaultLogSize)

// 	fmt.Printf("Number of items in a bucket is: %v\n", 1<<bshift)
// 	fmt.Printf("Number of hash functions is: %v\n", nhash)
// 	fmt.Printf("Length of buckets is: %v\n", len(c.buckets))

// 	fmt.Println("Inserting items...")
// 	for k, v := range gmap {
// 		c.Insert(k, v)
// 		ShowTable(c)
// 	}
// 	fmt.Printf("Length of buckets is: %v\n", len(c.buckets))

// 	for k, v := range gmap {
// 		cv, ok := c.Search(k)
// 		if !ok {
// 			t.Error("not ok:", k, v, cv)
// 			return
// 		}
// 		if reflect.DeepEqual(cv, v) == false {
// 			t.Error("got: ", cv, " expected: ", v)
// 			return
// 		}
// 	}

// 	if c.Len() != len(gmap) {
// 		t.Error("got: ", c.Len(), " expected: ", len(gmap))
// 		return
// 	}
// 	fmt.Println("LoadFactor:", c.LoadFactor())

// 	fmt.Println("Deleting items...")
// 	ndeleted := 0
// 	maxdelete := len(gmap) * 95 / 100
// 	for k := range gmap {
// 		if ndeleted >= maxdelete {
// 			break
// 		}

// 		c.Delete(k)
// 		if v, ok := c.Search(k); ok == true {
// 			t.Error("got: ", v)
// 			return
// 		}

// 		ndeleted++

// 		if c.Len() != len(gmap)-ndeleted {
// 			t.Error("got: ", c.Len(), " expected: ", len(gmap)-ndeleted)
// 			return
// 		}
// 	}
// }

// // cuckooPIR for both key and value of uint8
// func TestCuckooPIRUint8(t *testing.T){
// 	fmt.Printf("Totally %v items by %v hash functions, %v items in a bucket, %v buckets.\n", len(gmap), nhash, blen, tablen * nhash)
// 	// fmt.Println("Items to be inserted: ", gmap)
// 	c := NewCuckoo(DefaultLogSize)
// 	var m runtime.MemStats
// 	runtime.ReadMemStats(&m)

// 	start := time.Now()
// 	fmt.Println("Inserting items...")
// 	for k, v := range gmap {
// 		c.Insert(k, v)
// 	}
// 	printTime(start)
// 	ShowTable(c)
// 	fmt.Println("LoadFactor:", c.LoadFactor())

// 	fmt.Println("Creating a table...")
// 	TabMat := [nhash]*Matrix{}
// 	T := MatrixNew(uint64(tablen), uint64(blen))

// 	k := 0
// 	for t := 0; t < len(c.buckets); t += tablen {
// 		for i := t; i < t + tablen; i++ {
// 			for j := 0; j < blen; j++ {
// 				T.Set(uint64(c.buckets[i].vals[j]), uint64(i - t), uint64(j))
// 			}
// 		}
// 		TabMat[k] = T
// 		k += 1
// 		T = MatrixNew(uint64(tablen), uint64(blen))
// 	}

// 	// print nhash tables
// 	// fmt.Println(len(TabMat))
// 	// for j := 0; j < len(TabMat); j++ {
// 	// 	TabMat[j].Print()
// 	// }

// 	// run CuckooPIR for Tables[]
// 	N := uint64(tablen * blen)
// 	d := uint64(8)
// 	pir := CuckooPIR{}
// 	p := Params{1<<10,6.4,uint64(tablen),blen,32,1<<8}

// 	var Tables [nhash]*Database
// 	var wg sync.WaitGroup
// 	wg.Add(nhash)
// 	for i := 0; i < nhash; i++ {
// 		Tables[i] = MakeDBFromMat(N, d, &p, TabMat[i])
// 		// Tables[i].Data.Print()
// 		go RunPIR(&pir, Tables[i], p, rows, &wg)
// 	}
// 	wg.Wait()
// 	fmt.Println("Done")
// }


// setMatrixValue sets the value from the byte slice into the matrix at the given row and column.
func setMatrixValue(matrix *Matrix, value []byte, row, col uint64) {
    for i, b := range value {
        matrix.Set(uint64(b), row, col+uint64(i))
    }
}

// cuckooPIR for arbitray item, i.e., byte[] key with byte[] value
func TestCuckooPIRByte(t *testing.T){
	fmt.Printf("Totally %v items by %v hash functions, %v items in a bucket, %v buckets.\n", n, nhash, blen, tablen * nhash)
	// fmt.Println("Items to be inserted: ", gmap)
	c := NewCuckoo(DefaultLogSize)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	keys, values := randKV()

	start := time.Now()
	fmt.Println("Inserting items...")
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		v := values[i]
		c.Insert(k, v)
	}

	printTime(start)
	ShowTable(c)
	fmt.Println("LoadFactor:", c.LoadFactor())

	fmt.Println("Creating a table...")
	TabMat := [nhash]*Matrix{}
	k := 0
	for t := 0; t < len(c.buckets); t += tablen {
		T := MatrixNew(uint64(tablen), uint64(blen*valSize))
	
		for i := t; i < t+tablen; i++ {
			for j := 0; j < blen; j++ {
				row := uint64(i - t)
				col := uint64(j * valSize)
				setMatrixValue(T, c.buckets[i].vals[j], row, col)
			}
		}
	
		TabMat[k] = T
		k += 1
	}

	for i := range TabMat {
		TabMat[i].Print()
	}

	// run CuckooPIR for Tables[]
	N 	:= uint64(tablen * blen)
	d 	:= uint64(8)
	pir := CuckooPIR{}
	p 	:= Params{1<<10,6.4,uint64(tablen),blen,32,1<<8}

	var Tables [nhash]*Database
	var wg sync.WaitGroup
	wg.Add(nhash)
	for i := 0; i < nhash; i++ {
		Tables[i] = MakeDBFromMat(N, d, &p, TabMat[i])
		Tables[i].Data.Print()
		go RunPIR(&pir, Tables[i], p, rows, &wg)
	}
	wg.Wait()
	fmt.Println("Done")
}

func TestHash(t *testing.T){
    // seed 		:= []byte("my-random-seed")
    input 		:= []byte("Hello, World!")
    outputSize 	:= 4

	seedSize := 16 // 你需要的 seed 长度
	rand.Seed(time.Now().UnixNano())

	seed := make([]byte, seedSize)
	for i := range seed {
		seed[i] = byte(rand.Intn(256))
	}

    hash, err := FNV1a(input, seed, outputSize)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("FNV: %x\n", hash)
	fmt.Printf("SHA256: %x\n", sha256mac(input, seed))
}

func TestReseed(t *testing.T){
	c := NewCuckoo(DefaultLogSize)
	// fmt.Println(rand.Intn(256))
	fmt.Printf("%x\n", c.seed)
	c.reseed()
	fmt.Printf("%x\n", c.seed)
}
