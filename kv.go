// Copyright (c) 2014-2015 Utkan Güngördü <utkan@freeconsole.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cuckoopir

import (
	"math/rand"
	"fmt"
)

func generateRandomBytes(size int) []byte {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

func randKV() ([]Key, []Value) {
	keys := make([]Key, n)
	values := make([]Value, n)

	for i := 0; i < n; i++ {
		keys[i] = generateRandomBytes(keySize)
		values[i] = generateRandomBytes(valSize)
	}

	// print them out
	for i := 0; i < n; i++ {
		fmt.Printf("Key %d: %x ", i, keys[i])
		fmt.Printf("Value %d: %x\n", i, values[i])
	}

	return keys, values
}



// var (
// 	mapBytes    uint64
// 	cuckooBytes uint64
// )

// var (
// 	mbench map[Key]Value
// 	cbench *Cuckoo
// )

// func readAlloc() uint64 {
// 	var ms runtime.MemStats
// 	runtime.ReadMemStats(&ms)
// 	return ms.Alloc
// }

// func mkmap(n int) (map[Key]Value, []Key, []Value, uint64) {

// 	keys := make([]Key, n, n)
// 	vals := make([]Value, n, n)

// 	runtime.GC()
// 	before := readAlloc()

// 	m := make(map[Key]Value)
// 	for i := 0; i < n; i++ {
// 		k := Key(rand.Uint32())
// 		v := Value(k)
// 		m[k] = v
// 		keys[i] = k
// 		vals[i] = v
// 	}

// 	after := readAlloc()

// 	// fmt.Println("m, keys, vals, after - before", m, keys, vals, after - before)
// 	return m, keys, vals, after - before
// }

// func init() {
// 	gmap, gkeys, gvals, mapBytes = mkmap(n)
// }

// func TestZero(t *testing.T) {
// 	c := NewCuckoo(logsize)
// 	var v Value

// 	for i := 0; i < 10; i++ {
// 		c.Insert(0, v)
// 		_, ok := c.Search(0)
// 		if !ok {
// 			t.Error("search failed")
// 		}
// 	}
// }

// func TestSimple(t *testing.T) {
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

// func TestMem(t *testing.T) {
// 	runtime.GC()
// 	before := readAlloc()

// 	c := NewCuckoo(logsize)
// 	for k, v := range gmap {
// 		c.Insert(k, v)
// 	}

// 	after := readAlloc()

// 	cuckooBytes = after - before

// 	t.Log("LoadFactor:", c.LoadFactor())
// 	t.Log("Built-in map memory usage (MiB):", float64(mapBytes)/float64(1<<20))
// 	t.Log("Cuckoo hash  memory usage (MiB):", float64(cuckooBytes)/float64(1<<20))
// }

// func BenchmarkCuckooInsert(b *testing.B) {
// 	cbench = NewCuckoo(logsize)
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		cbench.Insert(gkeys[i%n], gvals[i%n])
// 	}
// }

// func BenchmarkCuckooSearch(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		cbench.Search(gkeys[i%n])
// 	}
// }

// func BenchmarkCuckooDelete(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		cbench.Delete(gkeys[i%n])
// 	}
// }

// func BenchmarkMapInsert(b *testing.B) {
// 	mbench = make(map[Key]Value)
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		mbench[gkeys[i%n]] = gvals[i%n]
// 	}
// }

// func BenchmarkMapSearch(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		_, _ = mbench[gkeys[i%n]]
// 	}
// }

// func BenchmarkMapDelete(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		delete(mbench, gkeys[i%n])
// 	}
// }
