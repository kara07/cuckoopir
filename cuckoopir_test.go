package cuckoopir

import (
	_ "encoding/csv"
	"fmt"
	_ "math"
	_ "os"
	_ "strconv"
	"testing"
	_ "strings"
	"reflect"
	"runtime"
	"sync"
	"math/rand"
	"time"
)


const LOGQ = uint64(32)
const SEC_PARAM = uint64(1 << 10)


func TestPIR(t *testing.T) {
	fmt.Println("Number of CPUs:", runtime.NumCPU())

    runtime.GOMAXPROCS(runtime.NumCPU())

	N := uint64(1 << 20)
	// Num        uint64 // number of DB entries.
	d := uint64(8)
	// Row_length uint64 // number of bits per DB entry.
	pir := CuckooPIR{}
	// p := pir.PickParams(N, d, SEC_PARAM, LOGQ)//return Params
	p := Params{1024,6.4,1<<10,1<<10,32,512}//return Params
	// p := Params{1024,6.4,1<<16,1<<14,32,512}//return Params
	// p := Params{1024,6.4,5120,1024,32,991}//return Params
	// type Params struct {
	// 	N     uint64  // LWE secret dimension
	// 	Sigma float64 // LWE error distribution stddev
	
	// 	M uint64 // DB height
	// 	T uint64 // DB width
	
	// 	Logq uint64 // (logarithm of) ciphertext modulus
	// 	P    uint64 // plaintext modulus
	// }

	DB := MakeRandomDB(N, d, &p)//return *Database
	// type Database struct {
	// 	Info DBinfo
	// 	Data *Matrix
	// }
	// type Matrix struct {
	// 	Rows uint64
	// 	Cols uint64
	// 	Data []C.Elem		//typedef uint32_t Elem;
	// }
	// fmt.Println(*DB.Data)
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 1; i++ {
		go RunPIR(&pir, DB, p, []uint64{11,45,14,19,19,810},&wg)
	}
	wg.Wait()
	fmt.Println("Done")
}

func TestCuckoo(t *testing.T) {
	fmt.Println(n)
	fmt.Println("gmap:", gmap)
	// fmt.Println("gkeys:", gkeys)
	// fmt.Println("gvals:", gvals)
	fmt.Println("mapBytes:", mapBytes)

	c := NewCuckoo(DefaultLogSize)

	fmt.Printf("Number of items in a bucket is: %v\n", 1<<bshift)
	fmt.Printf("Number of hash functions is: %v\n", nhash)
	fmt.Printf("Length of buckets is: %v\n", len(c.buckets))

	fmt.Println("Inserting items...")
	for k, v := range gmap {
		c.Insert(k, v)
		ShowTable(c)
	}
	fmt.Printf("Length of buckets is: %v\n", len(c.buckets))

	for k, v := range gmap {
		cv, ok := c.Search(k)
		if !ok {
			t.Error("not ok:", k, v, cv)
			return
		}
		if reflect.DeepEqual(cv, v) == false {
			t.Error("got: ", cv, " expected: ", v)
			return
		}
	}

	if c.Len() != len(gmap) {
		t.Error("got: ", c.Len(), " expected: ", len(gmap))
		return
	}
	fmt.Println("LoadFactor:", c.LoadFactor())

	fmt.Println("Deleting items...")
	ndeleted := 0
	maxdelete := len(gmap) * 95 / 100
	for k := range gmap {
		if ndeleted >= maxdelete {
			break
		}

		c.Delete(k)
		if v, ok := c.Search(k); ok == true {
			t.Error("got: ", v)
			return
		}

		ndeleted++

		if c.Len() != len(gmap)-ndeleted {
			t.Error("got: ", c.Len(), " expected: ", len(gmap)-ndeleted)
			return
		}
	}
}

func TestCuckooPIR(t *testing.T){
	fmt.Println("gmap:", gmap)
	c := NewCuckoo(DefaultLogSize)
	
	fmt.Printf("Number of items in a bucket is: %v\n", 1<<bshift)
	fmt.Printf("Number of hash functions is: %v\n", nhash)
	fmt.Printf("Initial length of buckets is: %v\n", len(c.buckets))

	fmt.Println("Inserting items...")
	for k, v := range gmap {
		c.Insert(k, v)
		// ShowTable(c)
	}
	fmt.Printf("Length of buckets is: %v\n", len(c.buckets))
	ShowTable(c)
	fmt.Println("LoadFactor:", c.LoadFactor())

	fmt.Println("Creating a table...")
	Tables := []*Matrix{}
	T := MatrixNew(uint64(tabLen), uint64(blen))

	for t := 0; t < len(c.buckets); t += tabLen {
		// fmt.Println("t:", t)
		for i := t; i < t + tabLen; i++ {
			// fmt.Println("i:", i)
			for j := 0; j < blen; j++ {
				T.Set(uint64(c.buckets[i].vals[j]), uint64(i - t), uint64(j))
			}
		}
		// T.Print()
		Tables = append(Tables, T)
		T = MatrixNew(uint64(tabLen), uint64(blen))
	}
	// print Tables
	for j := 0; j < len(Tables); j++ {
		Tables[j].Print()
	}
}

func TestMatrix(t *testing.T){

	// a := 3
	// ptr := &a
	// fmt.Println(ptr)
	fmt.Println(MatrixMul(Atest,Btest))
	// outpyt: missing type in composite literal
}

func TestInt(t *testing.T){
	// N := uint64(1 << 20)
	// Num        uint64 // number of DB entries.
	// d := uint64(8)
	// Row_length uint64 // number of bits per DB entry.
	p := Params{1024,6.4,4,2,32,991}//return Params
	a := uint32(1)
	fmt.Println(a-2)
	RandomMatrix := MatrixRand(4,2,0,991)
	RandomMatrix.Print()

	RandomMatrix.Sub(p.P / 2)
	RandomMatrix.Print()

	// RandomDB := MakeRandomDB(N, d, &p)
	// RandomDB.Data.Print()

	Rows := RandomMatrix.SelectRows(0,2)
	Rows.Print()

	row := RandomMatrix.SelectRow(2)
	row.Print()

	slice := []uint64{3,2}
	rows := RandomMatrix.SelectSparseRows(slice)
	rows.Print()
}

func TestMulAdd(t *testing.T){

	rand.Seed(time.Now().UnixNano())
	const numMultiplications = 9999999
	// results := make([]int, numMultiplications)

	c := 0
	for i := 0; i < numMultiplications; i++ {
		a := rand.Intn(1<<32)
		b := rand.Intn(1<<32)
		c = a + b
		fmt.Println(c)
	}

	// 输出结果
	// for i, result := range results {
	// 	fmt.Printf("第 %d 次乘法的结果是: %d\n", i+1, result)
	// }

}