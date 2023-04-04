package cuckoopir

import "math"
import "fmt"

type DBinfo struct {
	Num        uint64 // number of DB entries.
	Row_length uint64 // number of bits per DB entry.

	Packing uint64 // number of DB entries per Z_p elem, if log(p) > DB entry size.
	Ne      uint64 // number of Z_p elems per DB entry, if DB entry size > log(p).

	X uint64 // tunable param that governs communication,
	         // must be in range [1, ne] and must be a divisor of ne;
	         // represents the number of times the scheme is repeated.
	P    uint64 // plaintext modulus.
	Logq uint64 // (logarithm of) ciphertext modulus.

	// For in-memory DB compression
	Basis     uint64 
	Squishing uint64
	Cols      uint64
}

type Database struct {
	Info DBinfo
	Data *Matrix
}

func (DB *Database) Squish() {
	//fmt.Printf("Original DB dims: ")
	//DB.Data.Dim()

	DB.Info.Basis = 10
	DB.Info.Squishing = 3 
	DB.Info.Cols = DB.Data.Cols
	DB.Data.Squish(DB.Info.Basis, DB.Info.Squishing)

	//fmt.Printf("After squishing, with compression factor %d: ", DB.Info.Squishing)
	//DB.Data.Dim()

	// Check that params allow for this compression
	if (DB.Info.P > (1 << DB.Info.Basis)) || (DB.Info.Logq < DB.Info.Basis * DB.Info.Squishing) {
		panic("Bad params")
	}
}

func (DB *Database) Unsquish() {
	DB.Data.Unsquish(DB.Info.Basis, DB.Info.Squishing, DB.Info.Cols)
}

// Store the database with entries decomposed into Z_p elements, and mapped to [-p/2, p/2]
// Z_p elements that encode the same database entry are stacked vertically below each other.
func ReconstructElem(vals []uint64, index uint64, info DBinfo) uint64 {
	q := uint64(1 << info.Logq)

	for i, _ := range vals {
		vals[i] = (vals[i] + info.P/2) % q
		vals[i] = vals[i] % info.P
	}

	val := Reconstruct_from_base_p(info.P, vals)

	if info.Packing > 0 {
		val = Base_p((1 << info.Row_length), val, index%info.Packing)
	}

	return val
}

func (DB *Database) GetElem(i uint64) uint64 {
	if i >= DB.Info.Num {
		panic("Index out of range")
	}

	// fmt.Println("DB rows:", DB.Data.Rows, "DB cols:", DB.Data.Cols)
	// fmt.Println("Getting elem", i, "of DB row", i/DB.Data.Cols, "col", i%DB.Data.Cols)
	col := i % DB.Data.Cols
	row := i / DB.Data.Cols

	if DB.Info.Packing > 0 {
		new_i := i / DB.Info.Packing
		col = new_i % DB.Data.Cols
		row = new_i / DB.Data.Cols
	}

	var vals []uint64
	for j := row * DB.Info.Ne; j < (row+1)*DB.Info.Ne; j++ {
		vals = append(vals, DB.Data.Get(j, col))
	}

	return ReconstructElem(vals, i, DB.Info)
}

func SetupDB(Num, row_length uint64, p *Params) *Database {
	if (Num == 0) || (row_length == 0) {
		panic("Empty database!")
	}

	D := new(Database)

	D.Info.Num = Num
	D.Info.Row_length = row_length
	D.Info.P = p.P
	D.Info.Logq = p.Logq

	db_elems, elems_per_entry, entries_per_elem := Num_DB_entries(Num, row_length, p.P)
	fmt.Println("DB elems:", db_elems, "elems per entry:", elems_per_entry, "entries per elem:", entries_per_elem)
	D.Info.Ne = elems_per_entry
	D.Info.X = D.Info.Ne
	D.Info.Packing = entries_per_elem

	for D.Info.Ne%D.Info.X != 0 {
		D.Info.X += 1
	}

	D.Info.Basis = 0
	D.Info.Squishing = 0

	fmt.Printf("Total packed DB size is ~%f MB\n",
		float64(p.M*p.T)*math.Log2(float64(p.P))/(1024.0*1024.0*8.0))

	// if db_elems > p.M*p.T {
	// 	panic("Params and database size don't match")
	// }

	if p.M%D.Info.Ne != 0 {
		panic("Number of DB elems per entry must divide DB height")
	}

	return D
}

func MakeRandomDB(Num, row_length uint64, p *Params) *Database {
	D := SetupDB(Num, row_length, p)
	D.Data = MatrixRand(p.M, p.T, 0, p.P)

	// Map DB elems to [-p/2; p/2]
	// D.Data.Sub(p.P / 2)
	D.Data.Sub(p.P / 2)

	return D
}

// func MakeDB(Num, row_length uint64, p *Params, vals []uint64) *Database {
// 	D := SetupDB(Num, row_length, p)
// 	D.Data = MatrixZeros(p.M, p.T)

// 	if uint64(len(vals)) != Num {
// 		panic("Bad input DB")
// 	}

// 	if D.Info.Packing > 0 {
// 		// Pack multiple DB elems into each Z_p elem
// 		at := uint64(0)
// 		cur := uint64(0)
// 		coeff := uint64(1)
// 		for i, elem := range vals {
// 			cur += (elem * coeff)
// 			coeff *= (1 << row_length)
// 			if ((i+1)%int(D.Info.Packing) == 0) || (i == len(vals)-1) {
// 				D.Data.Set(cur, at/p.M, at%p.M)
// 				at += 1
// 				cur = 0
// 				coeff = 1
// 			}
// 		}
// 	} else {
// 		// Use multiple Z_p elems to represent each DB elem
// 		for i, elem := range vals {
// 			for j := uint64(0); j < D.Info.Ne; j++ {
// 				D.Data.Set(Base_p(D.Info.P, elem, j), (uint64(i)/p.M)*D.Info.Ne+j, uint64(i)%p.M)
// 			}
// 		}
// 	}

// 	// Map DB elems to [-p/2; p/2]
// 	D.Data.Sub(p.P / 2)

// 	return D
// }

func MakeDBFromMat(Num, row_length uint64, p *Params, mat *Matrix) *Database {
	D := SetupDB(Num, row_length, p)
	D.Data = mat

	// Map DB elems to [-p/2; p/2]
	D.Data.Sub(p.P / 2)

	return D
}