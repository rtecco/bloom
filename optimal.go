package bloom

/*
	This is a port of:
       https://raw.githubusercontent.com/facebookarchive/cassandra/master/src/org/apache/cassandra/utils/BloomCalculations.java

	These functions are used to choose correct values of "buckets per element" and "number of hash functions, k".

	The following calculations are taken from:
       http://www.cs.wisc.edu/~cao/papers/summary-cache/node8.html
       "Bloom Filters - the math"
*/

const (
	maxBuckets = 15
	minBuckets = 2

	maxK = 8
	minK = 1
)

var optimalKPerBuckets [16]int = [...]int{
	1, // dummy K for 0 buckets per element
	1, // dummy K for 1 buckets per elements
	1, 2, 3, 3, 4, 5, 5, 6, 7, 8, 8, 8, 8, 8,
}

/*
	In the following table, the row 'i' shows false positive rates if i buckets
	per element are used.  Column 'j' shows false positive rates if j hash
	functions are used.  The first row is 'i=0', the first column is 'j=0'.
	Each cell (i,j) the false positive rate determined by using i buckets per
	element and j hash functions.
*/

// the first column is a dummy column representing K=0
var probs [16][9]float32 = [16][9]float32{
	{1.0},      // dummy row representing 0 buckets per element
	{1.0, 1.0}, // dummy row representing 1 buckets per element
	{1.0, 0.393, 0.400},
	{1.0, 0.283, 0.237, 0.253},
	{1.0, 0.221, 0.155, 0.147, 0.160},
	{1.0, 0.181, 0.109, 0.092, 0.092, 0.101},
	{1.0, 0.154, 0.0804, 0.0609, 0.0561, 0.0578, 0.0638},
	{1.0, 0.133, 0.0618, 0.0423, 0.0359, 0.0347, 0.0364},
	{1.0, 0.118, 0.0489, 0.0306, 0.024, 0.0217, 0.0216, 0.0229},
	{1.0, 0.105, 0.0397, 0.0228, 0.0166, 0.0141, 0.0133, 0.0135, 0.0145},
	{1.0, 0.0952, 0.0329, 0.0174, 0.0118, 0.00943, 0.00844, 0.00819, 0.00846},
	{1.0, 0.0869, 0.0276, 0.0136, 0.00864, 0.0065, 0.00552, 0.00513, 0.00509},
	{1.0, 0.08, 0.0236, 0.0108, 0.00646, 0.00459, 0.00371, 0.00329, 0.00314},
	{1.0, 0.074, 0.0203, 0.00875, 0.00492, 0.00332, 0.00255, 0.00217, 0.00199},
	{1.0, 0.0689, 0.0177, 0.00718, 0.00381, 0.00244, 0.00179, 0.00146, 0.00129},
	{1.0, 0.0645, 0.0156, 0.00596, 0.003, 0.00183, 0.00128, 0.001, 0.000852},
}

/*
	Given a maximum tolerable false positive probability, compute a Bloom specification which will
	give less than the specified false positive rate, but minimize the number of buckets per element
	and the number of hash functions used.  Because bandwidth (and therefore total bitvector size) is
	considered more expensive than computing power, preference is given to minimizing buckets per element
	rather than number of hash functions.
*/
func ComputeOptimal(maxFalsePositiveProb float32) (bucketsPerElement int, k int) {

	bucketsPerElement = 2
	k = optimalKPerBuckets[bucketsPerElement]

	//
	// handle trivial cases

	if maxFalsePositiveProb >= probs[minBuckets][minK] {
		return
	}

	if maxFalsePositiveProb < probs[maxBuckets][maxK] {

		bucketsPerElement = maxBuckets
		k = maxK

		return
	}

	//
	// find the minimum required number of buckets

	for probs[bucketsPerElement][k] > maxFalsePositiveProb {
		bucketsPerElement++
		k = optimalKPerBuckets[bucketsPerElement]
	}

	//
	// the number of buckets is sufficient, see if we can relax k w/o losing precision

	for probs[bucketsPerElement][k-1] <= maxFalsePositiveProb {
		k--
	}

	return
}
