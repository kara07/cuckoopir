package cuckoopir

// configurable variables (for tuning the algorithm)
const (
	bshift                = 2   // Number of items in a bucket is 1<<bshift.
	// bshift                = 3   // Number of items in a bucket is 1<<bshift.
	nhashshift            = 1   // Number of hash functions is 1<<nhashshift. (With SSE2, we can do 4 at once).
	// nhashshift            = 2   // Number of hash functions is 1<<nhashshift. (With SSE2, we can do 4 at once).
	shrinkFactor          = 0   // A shrink will be triggered when the load factor goes below 2^(-shrinkFactor). Setting this to 0 will disable shrinking and avoid potential new allocations.
	rehashThreshold       = 0.9 // If the load factor is below rehashThreshold, Insert will try to rehash everything before actually growing.
	randomWalkCoefficient = 2   // A multiplicative coefficient best determined by benchmarks. The optimal value depends on bshift and nhashshift.
	stashSize             = 4   // Size of stash (see Kirsch, Adam, Michael Mitzenmacher, and Udi Wieder. "More robust hashing: Cuckoo hashing with a stash." SIAM Journal on Computing 39.4 (2009): 1543-1561.)
)

// other configurable variables
const (
	gc             = false      // trigger GC after every alloc (which happens during grow).
	// DefaultLogSize = 8 + bshift // A reasonable logsize value for NewCuckoo for use when the number of items to be inserted is not known ahead.
	DefaultLogSize = 2 + bshift	// minimum number of buckets
)

// Key must be an integer-type.
// type Key uint32
type Key uint8

// Value can be anything, replace this to match your needs (not using unsafe.Pointer to avoid the overhead to store additional pointer or interface{} which comes with a worse overhead).
// type Value uint32
type Value uint8
