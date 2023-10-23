package localvectordb

import (
	"math"
)

type Float interface {
	~float32 | ~float64
}

// Cosine 计算cosine [-1,1],越大越相似
func Cosine[T Float](a []T, b []T) T {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(float64(b[k]), 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(float64(a[k]), 2)
			continue
		}
		sumA += float64(a[k]) * float64(b[k])
		s1 += math.Pow(float64(a[k]), 2)
		s2 += math.Pow(float64(b[k]), 2)
	}
	return T(sumA / (math.Sqrt(s1) * math.Sqrt(s2)))
}

type DocWithScore struct {
	*Doc
	Score float32
}
type bigHeap []*DocWithScore

func (b bigHeap) Len() int {
	return len(b)
}

func (b bigHeap) Less(i, j int) bool {
	return b[i].Score < b[j].Score
}

func (b bigHeap) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b *bigHeap) Push(x interface{}) {
	*b = append(*b, x.(*DocWithScore))
}

func (b *bigHeap) Pop() interface{} {
	l := len(*b) - 1
	a := (*b)[l]
	(*b)[l] = nil
	*b = (*b)[:l]
	return a
}
