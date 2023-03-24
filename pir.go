package cuckoopir

import (
	"fmt"
	_ "os"
	"runtime"
	"runtime/debug"
	_ "runtime/pprof"
	"time"
	"sync"
//	"math"
)

// Defines the interface for PIR with preprocessing schemes
// type PIR interface {
// 	Name() string

// 	GetBW(info DBinfo, p Params)

// 	Init(info DBinfo, p Params) State

// 	Setup(DB *Database, shared State, p Params) (State, Msg)

// 	Query(L []uint64, shared State, p Params, info DBinfo) (State, Msg)

// 	Response(DB *Database, query Msg, server State, shared State, p Params) Msg

// 	Extract(offline Msg, query Msg, answer Msg, shared State, client State, p Params, info DBinfo) Msg

// 	Reset(DB *Database, p Params) // reset DB to its correct state, if modified during execution
// }
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
	DB.Data.Print()

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
		
	runtime.GC()
	printTime(start)
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

	expectedRows := DB.Data.SelectSparseRows(rows)
	for i, v := range V.Data {
		if v != expectedRows.Data[i]{
			fmt.Printf("Expected result: %d, %d\n",i, expectedRows.Data[i])
			fmt.Printf("Actual result: %d, %d\n",i, v)
			// panic("Result Failure!")
		}
	}
	fmt.Println("Success!")
	

	runtime.GC()
	debug.SetGCPercent(100)
	// 
	wg.Done()
	return rate, bw
}

// Run full PIR scheme (offline + online phases), where the transmission of the A matrix is compressed.
// func RunPIRCompressed(pi PIR, DB *Database, p Params, i []uint64) (float64, float64) {
//         fmt.Printf("Executing %s\n", pi.Name())
//         //fmt.Printf("Memory limit: %d\n", debug.SetMemoryLimit(math.MaxInt64))
//         debug.SetGCPercent(-1)

//         num_queries := uint64(len(i))
//         if DB.Data.Rows/num_queries < DB.Info.Ne {
//                 panic("Too many queries to handle!")
//         }
//         batch_sz := DB.Data.Rows / (DB.Info.Ne * num_queries) * DB.Data.Cols
//         bw := float64(0)

//         server_shared_state, comp_state := pi.InitCompressed(DB.Info, p)
//         client_shared_state := pi.DecompressState(DB.Info, p, comp_state)

//         fmt.Println("Setup...")
//         start := time.Now()
//         server_state, offline_download := pi.Setup(DB, server_shared_state, p)
//         printTime(start)
//         comm := float64(offline_download.Size() * uint64(p.Logq) / (8.0 * 1024.0))
//         fmt.Printf("\t\tOffline download: %f KB\n", comm)
//         bw += comm
//         runtime.GC()

//         fmt.Println("Building query...")
//         start = time.Now()
//         var client_state []State
//         var query MsgSlice
//         for index, _ := range i {
//                 index_to_query := i[index] + uint64(index)*batch_sz
//                 cs, q := pi.Query(index_to_query, client_shared_state, p, DB.Info)
//                 client_state = append(client_state, cs)
//                 query.Data = append(query.Data, q)
//         }
//         runtime.GC()
//         printTime(start)
//         comm = float64(query.Size() * uint64(p.Logq) / (8.0 * 1024.0))
//         fmt.Printf("\t\tOnline upload: %f KB\n", comm)
//         bw += comm
//         runtime.GC()

//         fmt.Println("Answering query...")
//         start = time.Now()
//         answer := pi.Answer(DB, query, server_state, server_shared_state, p)
//         elapsed := printTime(start)
//         rate := printRate(p, elapsed, len(i))
//         comm = float64(answer.Size() * uint64(p.Logq) / (8.0 * 1024.0))
//         fmt.Printf("\t\tOnline download: %f KB\n", comm)
//         bw += comm
//         runtime.GC()

//         pi.Reset(DB, p)
//         fmt.Println("Reconstructing...")
//         start = time.Now()

//         for index, _ := range i {
//                 index_to_query := i[index] + uint64(index)*batch_sz
//                 val := pi.Recover(index_to_query, uint64(index), offline_download,
//                                   query.Data[index], answer, client_shared_state,
//                                   client_state[index], p, DB.Info)

//                 if DB.GetElem(index_to_query) != val {
//                         fmt.Printf("Batch %d (querying index %d -- row should be >= %d): Got %d instead of %d\n",
//                                 index, index_to_query, DB.Data.Rows/4, val, DB.GetElem(index_to_query))
//                         panic("Reconstruct failed!")
//                 }
//         }
//         fmt.Println("Success!")
//         printTime(start)

//         runtime.GC()
//         debug.SetGCPercent(100)
//         return rate, bw
// }
