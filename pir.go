package cuckoopir

import (
	"fmt"
	_ "os"
	"runtime"
	"runtime/debug"
	_ "runtime/pprof"
	"time"
	"sync"
)

// Defines the interface for PIR with preprocessing schemes
type PIR interface {
	Name() string

	GetBW(info DBinfo, p Params)

	Init(info DBinfo, p Params) *Matrix

	Setup(DB *Database, A *Matrix, p Params) *Matrix

	Query(L []uint64, A *Matrix, p Params, info DBinfo) (*Matrix, *Matrix)

	Response(DB *Database, Q *Matrix, shared *Matrix, p Params) *Matrix

	Extract(M *Matrix, R *Matrix, S *Matrix, p Params, info DBinfo) *Matrix

	Reset(DB *Database, p Params) // reset DB to its correct state, if modified during execution
}

// Run full PIR scheme (offline + online phases).
// func RunPIR(pi PIR, DB *Database, p Params, rows []uint64) (float64, float64) {
func RunPIR(pi PIR, DB *Database, p Params, rows []uint64, wg *sync.WaitGroup) (float64, float64) {
	fmt.Printf("Executing %s\n", pi.Name())
	//fmt.Printf("Memory limit: %d\n", debug.SetMemoryLimit(math.MaxInt64))
	// debug.SetGCPercent(-1)

	num_queries := uint64(len(rows))
	fmt.Printf("Number of %d queries at once: %v\n", num_queries, rows)
	bw := float64(0)

	// Print database
	fmt.Printf("Database: ")
	// DB.Data.Print()

	// derive *Matrix A as State.Data[0]
	shared_state := pi.Init(DB.Info, p)

	fmt.Println("Setup...")
	start := time.Now()
	offline_download := pi.Setup(DB, shared_state, p)
	printTime(start)
	comm := float64(offline_download.Size() * uint64(p.Logq) / (8.0 * 1024.0))
	fmt.Printf("\t\tOffline download: %f KB\n", comm)
	bw += comm
	runtime.GC()

	fmt.Println("Building query...")
	start = time.Now()
	indexes_to_query := rows
	client_state, query := pi.Query(indexes_to_query, shared_state, p, DB.Info)
		
	printTime(start)
	runtime.GC()

	comm = float64(query.Size() * uint64(p.Logq) / (8.0 * 1024.0))
	fmt.Printf("\t\tOnline upload: %f KB\n", comm)
	bw += comm
	runtime.GC()

	fmt.Println("Responsing query...")
	start = time.Now()
	response := pi.Response(DB, query, shared_state, p)
	elapsed := printTime(start)
	rate := printRate(p, elapsed, len(rows))
	comm = float64(response.Size() * uint64(p.Logq) / (8.0 * 1024.0))
	fmt.Printf("\t\tOnline download: %f KB\n", comm)
	bw += comm
	runtime.GC()

	fmt.Println("Extracting...")
	start = time.Now()
	V := pi.Extract(offline_download, response, client_state, p, DB.Info)
	printTime(start)

	// Recover to [-p/2, p/2] and verify
	expectedRows := DB.Data.SelectSparseRows(rows)
	for i, _:= range V.Data {
		if uint64(expectedRows.Data[i]) <= p.P/2 {
			if V.Data[i] != expectedRows.Data[i]{
				fmt.Printf("Expected result: %d, %d\n",i, expectedRows.Data[i])
				fmt.Printf("Actual result: %d, %d\n",i, V.Data[i])
				panic("Result Failure!")
			}
		}else{
			// V.Data[i] -= C.Elem(p.P)
			V.AddByIndex(-p.P, uint64(i))
			if V.Data[i] != expectedRows.Data[i]{
				fmt.Printf("Expected result: %d, %d\n",i, expectedRows.Data[i])
				fmt.Printf("Actual result: %d, %d\n",i, V.Data[i])
				panic("Result Failure!")
			}
			// if uint32(v) != uint32(expectedRows.Data[i]) + uint32(p.P){
			// 	fmt.Printf("Expected result: %d, %d\n",i, uint32(expectedRows.Data[i]) + uint32(p.P))
			// 	fmt.Printf("Actual result: %d, %d\n",i, v)
			// }
		}
		V.AddByIndex(p.P/2, uint64(i))
	}
	fmt.Println("Extracted: ")
	// V.Print()
	fmt.Println("Success!")
	

	runtime.GC()
	debug.SetGCPercent(100)
	// 
	wg.Done()
	return rate, bw
}