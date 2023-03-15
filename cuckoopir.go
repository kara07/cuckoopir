package cuckoopir

// #cgo CFLAGS: -O3 -march=native
// #include "pir.h"
import "C"
import "fmt"
// import "sync"

type CuckooPIR struct{}

func (pi *CuckooPIR) Name() string {
	return "CuckooPIR"
}

func (pi *CuckooPIR) GetBW(info DBinfo, p Params) {
	offline_download := float64(p.M*p.N*p.Logq) / (8.0 * 1024.0)
	fmt.Printf("\t\tOffline download: %d KB\n", uint64(offline_download))

	online_upload := float64(p.M*p.Logq) / (8.0 * 1024.0)
	fmt.Printf("\t\tOnline upload: %d KB\n", uint64(online_upload))

	online_download := float64(p.M*p.Logq) / (8.0 * 1024.0)
	fmt.Printf("\t\tOnline download: %d KB\n", uint64(online_download))
}

func (pi *CuckooPIR) Init(info DBinfo, p Params) *Matrix {
	A := MatrixRand(p.M, p.N, p.Logq, 0)
	return A
}


func (pi *CuckooPIR) Setup(DB *Database, A *Matrix, p Params) *Matrix {
	AT := A.TransposeCopy()
	fmt.Println("A:", A.Rows, "x", A.Cols)
	fmt.Println("AT:", AT.Rows, "x", AT.Cols)
	fmt.Println("DB.Data:", DB.Data.Rows, "x", DB.Data.Cols)
	
	// M := MatrixMul(DB.Data, A)
	M := MatrixMul(AT, DB.Data)
	fmt.Println("M:", M.Rows, "x", M.Cols)
	
	return M
}

func (pi *CuckooPIR) Query(L []uint64, A *Matrix, p Params, info DBinfo) (*Matrix, *Matrix) {
	fmt.Println("A:", A.Rows, "x", A.Cols)

	S := MatrixRand(p.N, uint64(len(L)), p.Logq, 0)
	fmt.Println("S:", S.Rows, "x", S.Cols)

	Q := MatrixMul(A, S)//type *Matrix
	E := MatrixGaussian(p.M, uint64(len(L)))
	Q.MatrixAdd(E)
	// col := i % DB.Data.Cols
	// row := i / DB.Data.Cols
	for i, j := range L {
		index := uint64(j-1) * Q.Cols + uint64(i)
		Q.Data[index] += C.Elem(p.Delta())
	}
	// Qhat.Data[L[0]/p.T] += C.Elem(p.Delta())
	// fmt.Printf("query type is %T\n", query)

	// fmt.Println("query:", query)
	// fmt.Println("imodp.M:", i%p.M)
	// fmt.Println("query.Data[imodp.M]:", query.Data[i%p.M])
	fmt.Println("Q:", Q.Rows, "x", Q.Cols)

	return S, Q
}

func (pi *CuckooPIR) Response(DB *Database, Q *Matrix, shared *Matrix, p Params) *Matrix {
	fmt.Println("Q:", Q.Rows, "x", Q.Cols)

	QT := Q.TransposeCopy()
	fmt.Println("QT:", QT.Rows, "x", QT.Cols)

	R := MatrixMul(QT,DB.Data)
	fmt.Println("R:", R.Rows, "x", R.Cols)

	return R
}

func (pi *CuckooPIR) Extract(M *Matrix, R *Matrix, S *Matrix, p Params, info DBinfo) *Matrix {

	// col := i % p.M
	S.Transpose()
	Mhat := MatrixMul(S, M)
	R.MatrixSub(Mhat)
	// Recover each Z_p element that makes up the desired database entry
	
	R.Round(p)
	V := R
	fmt.Println("V:", V.Rows, "x", V.Cols)

	//fmt.Printf("Reconstructing row %d: %d\n", j, denoised)

	return V
}

func (pi *CuckooPIR) Reset(DB *Database, p Params) {
	// Uncompress the database, and map its entries to the range [-p/2, p/2].
	DB.Unsquish()
	DB.Data.Sub(p.P / 2)
}
