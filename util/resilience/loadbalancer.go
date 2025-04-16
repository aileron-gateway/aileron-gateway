// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright The AILERON Gateway Authors

package resilience

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/rand/v2"
	"slices"
	"sync"
)

type Entry interface {
	ID() string
	Active() bool
	Weight() int
	Hint() int
}

type LoadBalancer[T any] interface {
	Entries() []T
	Add(...T)
	Remove(ids string)
	Get(hint int) T
}

type RoundRobinLB[T Entry] struct {
	mu sync.Mutex

	// weights is the list of weights for the entries.
	// The length of the slice MUST be the same with lb.entries.
	// The length of the weights can not be 0.
	// weights can not be 0 or negative values.
	weights []int
	// entries is the list of all load balance targets.
	// The length of the slice MUST be the same with lb.weights.
	// The length of the entries can be 0.
	// entries MUST NOT contain nil.
	entries []T

	// curPos is the index of current target entry.
	// Current target = entries[curPos].
	// 0 <= curPos <= len(entries)
	curPos int
	// curWeights is the remaining weights of the entry at curPos.
	// That means len(curWeights) == len(entries).
	// curWeights will be decremented for each Get.
	curWeights []int
}

func (lb *RoundRobinLB[T]) Entries() []T {
	return lb.entries
}

func (lb *RoundRobinLB[T]) Add(es ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	size := len(es)
	lb.entries = append(lb.entries, make([]T, 0, size)...)
	lb.weights = append(lb.weights, make([]int, 0, size)...)
	lb.curWeights = append(lb.curWeights, make([]int, 0, size)...)

	for _, e := range es {
		if e.Weight() <= 0 {
			continue
		}
		lb.entries = append(lb.entries, e)
		lb.weights = append(lb.weights, e.Weight())
		lb.curWeights = append(lb.curWeights, e.Weight())
	}
}

func (lb *RoundRobinLB[T]) Remove(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, e := range lb.entries {
		if e.ID() != id {
			continue
		}
		lb.entries = append(lb.entries[:i], lb.entries[i+1:]...)
		lb.weights = append(lb.weights[:i], lb.weights[i+1:]...)
		lb.curWeights = append(lb.curWeights[:i], lb.curWeights[i+1:]...)
		if lb.curPos > i {
			lb.curPos -= 1
		}
		break // There might be an another entry with the same id.
	}
}

func (lb *RoundRobinLB[T]) Get(_ int) T {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	n := len(lb.entries)
	for i := 0; i < n; i++ {
		if lb.curPos >= n {
			lb.curPos = 0
		}

		e := lb.entries[lb.curPos]
		if !e.Active() {
			lb.curWeights[lb.curPos] = lb.weights[lb.curPos]
			lb.curPos += 1
			continue
		}

		lb.curWeights[lb.curPos] -= 1
		if lb.curWeights[lb.curPos] <= 0 {
			lb.curWeights[lb.curPos] = lb.weights[lb.curPos]
			lb.curPos += 1
		}
		return e
	}

	var tmp T
	return tmp
}

type RandomLB[T Entry] struct {
	mu sync.RWMutex

	// weights is the list of weights for the entries.
	// The length of the slice MUST be the same with lb.entries.
	// The length of the weights can not be 0.
	// weights can not be 0 or negative values.
	weights []int
	// entries is the list of all load balance targets.
	// The length of the slice MUST be the same with lb.weights.
	// The length of the entries can be 0.
	// entries MUST NOT contain nil.
	entries []T

	// shuffled is the shuffled index for entries.
	// Use like entries[shuffled[j]].
	shuffled []int
	// curPos is the index of current target entry.
	// Current target = entries[curPos].
	// 0 <= curPos <= len(entries)
	curPos int
	// curWeights is the remaining weights of the entry at curPos.
	// That means len(curWeights) == len(entries).
	// curWeights will be decremented for each Get.
	curWeights []int
}

func (lb *RandomLB[T]) Entries() []T {
	return lb.entries
}

func (lb *RandomLB[T]) Add(es ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	size := len(es)
	lb.entries = append(lb.entries, make([]T, 0, size)...)
	lb.weights = append(lb.weights, make([]int, 0, size)...)
	lb.shuffled = append(lb.shuffled, make([]int, 0, size)...)
	lb.curWeights = append(lb.curWeights, make([]int, 0, size)...)

	for _, e := range es {
		if e.Weight() <= 0 {
			continue
		}
		lb.entries = append(lb.entries, e)
		lb.weights = append(lb.weights, e.Weight())
		lb.shuffled = append(lb.shuffled, len(lb.entries)-1)
		lb.curWeights = append(lb.curWeights, e.Weight())
	}
}

func (lb *RandomLB[T]) Remove(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, e := range lb.entries {
		if e.ID() != id {
			continue
		}
		lb.entries = append(lb.entries[:i], lb.entries[i+1:]...)
		lb.weights = append(lb.weights[:i], lb.weights[i+1:]...)

		j := slices.Index(lb.shuffled, i)
		lb.shuffled = append(lb.shuffled[:j], lb.shuffled[j+1:]...)
		lb.curWeights = append(lb.curWeights[:j], lb.curWeights[j+1:]...)
		if lb.curPos > j {
			lb.curPos -= 1
		}
		for k := range lb.shuffled {
			if lb.shuffled[k] > i {
				lb.shuffled[k] -= 1
			}
		}
		break // There might be an another entry with the same id.
	}
}

func (lb *RandomLB[T]) Get(_ int) T {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	n := len(lb.entries)

	for i := 0; i < n; i++ {
		if lb.curPos >= n {
			lb.curPos = 0
			rand.Shuffle(n, func(x, y int) {
				if x < i || y < i {
					return
				}
				lb.shuffled[x], lb.shuffled[y] = lb.shuffled[y], lb.shuffled[x]
				lb.curWeights[x], lb.curWeights[y] = lb.curWeights[y], lb.curWeights[x]
			})
		}

		e := lb.entries[lb.shuffled[lb.curPos]]
		if !e.Active() {
			lb.curWeights[lb.curPos] = e.Weight()
			lb.curPos += 1
			continue
		}

		lb.curWeights[lb.curPos] -= 1
		if lb.curWeights[lb.curPos] <= 0 {
			lb.curWeights[lb.curPos] = e.Weight()
			lb.curPos += 1
		}

		return e
	}

	var tmp T
	return tmp
}

type DirectHashLB[T Entry] struct {
	mu sync.RWMutex

	// weights is the list of weights for the entries.
	// The length of the slice MUST be the same with lb.entries.
	// The length of the weights can not be 0.
	// weights can not be 0 or negative values.
	weights []int
	// entries is the list of all load balance targets.
	// The length of the slice MUST be the same with lb.weights.
	// The length of the entries can be 0.
	// entries MUST NOT contain nil.
	entries []T

	// sumWeights is the sum of the all entry wights.
	sumWeights int
	// cumulates is the cumulative weights of the entries.
	// cumulates[i] = weights[0] + weights[1] + ... + weights[i-1]
	cumulates []int
}

func (lb *DirectHashLB[T]) Entries() []T {
	return lb.entries
}

func (lb *DirectHashLB[T]) updateCumulates() {
	cum := 0
	for i, e := range lb.entries {
		cum += e.Weight()
		lb.cumulates[i] = cum
	}
	lb.sumWeights = cum
}

func (lb *DirectHashLB[T]) Add(es ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	size := len(es)
	lb.entries = append(lb.entries, make([]T, 0, size)...)
	lb.weights = append(lb.weights, make([]int, 0, size)...)
	lb.cumulates = append(lb.cumulates, make([]int, 0, size)...)

	for _, e := range es {
		if e.Weight() <= 0 {
			continue
		}
		lb.entries = append(lb.entries, e)
		lb.weights = append(lb.weights, e.Weight())
		lb.cumulates = append(lb.cumulates, 0) // Update later.
	}
	slices.SortStableFunc(lb.entries, func(e1, e2 T) int { return e2.Weight() - e1.Weight() })
	slices.SortStableFunc(lb.weights, func(w1, w2 int) int { return w2 - w1 })
	lb.updateCumulates()
}

func (lb *DirectHashLB[T]) Remove(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, e := range lb.entries {
		if e.ID() != id {
			continue
		}
		lb.entries = append(lb.entries[:i], lb.entries[i+1:]...)
		lb.weights = append(lb.weights[:i], lb.weights[i+1:]...)
		lb.cumulates = append(lb.cumulates[:i], lb.cumulates[i+1:]...)
		break // There might be an another entry with the same id.
	}
	lb.updateCumulates()
}

func (lb *DirectHashLB[T]) Get(hint int) T {
	if hint < 0 {
		hint = -hint
	}

	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if lb.sumWeights <= 0 {
		var tmp T
		return tmp // entry not available
	}

	target := hint%lb.sumWeights + 1 // 1 <= target <= lb.sumWeights
	for i, c := range lb.cumulates {
		if c < target {
			continue
		}
		return lb.entries[i]
	}

	var tmp T
	return tmp // Entry not available
}

type RingHashLB[T Entry] struct {
	mu sync.RWMutex

	// weights is the list of weights for the entries.
	// The length of the slice MUST be the same with lb.entries.
	// The length of the weights can not be 0.
	// weights can not be 0 or negative values.
	weights []int
	// entries is the list of all load balance targets.
	// The length of the slice MUST be the same with lb.weights.
	// The length of the entries can be 0.
	// entries MUST NOT contain nil.
	entries []T

	// Size is the length of ring hash table.
	// This MUST be 100 times larger than the length of entries.
	// When Size<100*len(entries), it is fixed to 100*len(entries)
	// This size can not be change dynamically.
	Size int

	// size is the determined size of the ring hash table.
	// size and Size may be different if an invalid value was set to the Size.
	size int
	// ring is the ring hash table.
	// len(ring) == size
	ring []int
}

func (lb *RingHashLB[T]) Entries() []T {
	return lb.entries
}

func (lb *RingHashLB[T]) Add(es ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	size := len(es)
	lb.entries = append(lb.entries, make([]T, 0, size)...)
	lb.weights = append(lb.weights, make([]int, 0, size)...)

	for _, e := range es {
		if e.Weight() <= 0 {
			continue
		}
		lb.entries = append(lb.entries, e)
		lb.weights = append(lb.weights, e.Weight())
	}

	lb.size = lb.Size
	if lb.Size < 100*len(lb.entries) {
		lb.size = 100 * len(lb.entries)
	}
	if len(lb.ring) < lb.size {
		lb.ring = append(lb.ring, make([]int, lb.size-len(lb.ring))...)
	}
	for i := 0; i < lb.size; i++ {
		lb.ring[i] = -1 // Reset tables.
	}

	sumWeight := 0
	for i := 0; i < len(lb.entries); i++ {
		sumWeight += lb.weights[i]
	}

	for i := 0; i < len(lb.entries); i++ {
		targetNum := (lb.size * lb.weights[i] / (2 * sumWeight))
		count := 0
		h := fnv.New64()
		h.Write([]byte(lb.entries[i].ID()))
		for {
			h.Write(h.Sum(nil))
			n := binary.BigEndian.Uint64(h.Sum(nil))
			pos := int(n % uint64(lb.size)) //nolint:gosec // G115: integer overflow conversion int -> uint64
			if lb.ring[pos] < 0 {
				lb.ring[pos] = i
				count += 1
				if count >= targetNum {
					break
				}
			}
		}
	}
}

func (lb *RingHashLB[T]) Remove(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, e := range lb.entries {
		if e.ID() != id {
			continue
		}
		lb.entries = append(lb.entries[:i], lb.entries[i+1:]...)
		lb.weights = append(lb.weights[:i], lb.weights[i+1:]...)
		for j := 0; j < lb.size; j++ {
			if lb.ring[j] == i {
				lb.ring[j] = -1
			} else if lb.ring[j] > i {
				lb.ring[j] -= 1
			}
		}
		break // There might be an another entry with the same id.
	}
}

func (lb *RingHashLB[T]) Get(hint int) T {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if lb.size <= 0 {
		var tmp T
		return tmp
	}

	// pos is the initial ring position
	// determined by the request's hash.
	pos := hint % lb.size

	// inActives holds the index of inactive entries.
	// When the length of inActives is the same as the length of lb.entries,
	// it means all entries are inactive.
	var inActives []int

	// Lookup entries from the initial position in clockwise order.
	for i := 0; i < lb.size; i++ {
		if pos >= lb.size {
			pos = 0 // Reset position.
		}
		index := lb.ring[pos]
		if index < 0 {
			pos += 1
			continue // Entry not assigned.
		}
		if slices.Contains(inActives, index) {
			pos += 1
			continue // Entry is inactive.
		}

		e := lb.entries[index]
		if !e.Active() {
			inActives = append(inActives, index)
			pos += 1
			if len(inActives) == len(lb.entries) {
				var tmp T
				return tmp // All entries is inactive.
			}
			continue // Entry is inactive.
		}

		return e
	}

	var tmp T
	return tmp
}

type MaglevLB[T Entry] struct {
	mu sync.RWMutex

	// weights is the list of weights for the entries.
	// The length of the slice MUST be the same with lb.entries.
	// The length of the weights can not be 0.
	// weights can not be 0 or negative values.
	weights []int
	// entries is the list of all load balance targets.
	// The length of the slice MUST be the same with lb.weights.
	// The length of the entries can be 0.
	// entries MUST NOT contain nil.
	entries []T

	// Size is the size of lookup table.
	// Size MUST be a prime number and should be grater than 10*len(lb.entries).
	// See https://research.google/pubs/maglev-a-fast-and-reliable-software-network-load-balancer/
	Size int

	// sumWeight is the sum of the all entry weights.
	sumWeight int
	// size is the determined size of the lookup table.
	// size and Size may be different if an invalid value was set to the Size.
	size int
	// table is the lookup table.
	// The length is the same as size.
	table []int
	// updated notifies if the table update finished.
	updated chan struct{}
}

func (lb *MaglevLB[T]) Entries() []T {
	return lb.entries
}

func (lb *MaglevLB[T]) Add(es ...T) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	size := len(es)
	lb.entries = append(lb.entries, make([]T, 0, size)...)
	lb.weights = append(lb.weights, make([]int, 0, size)...)

	for _, e := range es {
		if e.Weight() <= 0 {
			continue
		}
		lb.entries = append(lb.entries, e)
		lb.weights = append(lb.weights, e.Weight())
		lb.sumWeight += e.Weight()
	}
	lb.update()
}

func (lb *MaglevLB[T]) Remove(id string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i, e := range lb.entries {
		if e.ID() != id {
			continue
		}
		lb.entries = append(lb.entries[:i], lb.entries[i+1:]...)
		lb.weights = append(lb.weights[:i], lb.weights[i+1:]...)
		lb.sumWeight -= e.Weight()
		break // There might be an another entry with the same id.
	}
	lb.update()
}

func (lb *MaglevLB[T]) update() {
	if lb.updated == nil {
		lb.updated = make(chan struct{})
		defer func() {
			close(lb.updated)
			lb.updated = nil
		}()
	} else {
		<-lb.updated
		return
	}

	lb.size = lb.Size
	if lb.Size < 10*len(lb.entries) {
		lb.size = 10 * len(lb.entries)
	}
	if lb.sumWeight > lb.size {
		rate := float32(lb.size) / float32(lb.sumWeight)
		lb.size = 0
		for i := 0; i < len(lb.entries); i++ {
			lb.weights[i] = int(1 + rate*float32(lb.entries[i].Weight()))
			lb.size += lb.weights[i]
		}
	}
	lb.size = genPrimeEuler(lb.size)

	n := len(lb.entries)
	offset := make([]int, n)
	skip := make([]int, n)

	var inActives []int
	for i, e := range lb.entries {
		if !e.Active() || lb.weights[i] <= 0 {
			inActives = append(inActives, i)
			continue
		}
		h1 := fnv.New32()
		h1.Write([]byte(lb.entries[i].ID()))
		h2 := fnv.New32a()
		h2.Write([]byte(lb.entries[i].ID()))
		offset[i] = int(binary.BigEndian.Uint32(h1.Sum(nil)) % uint32(lb.size))   //nolint:gosec // G115: integer overflow conversion int -> uint32
		skip[i] = 1 + int(binary.BigEndian.Uint32(h2.Sum(nil))%uint32(lb.size-1)) //nolint:gosec // G115: integer overflow conversion int -> uint32
	}

	permTable := make([][]int, lb.size)
	lookupTable := make([]int, lb.size)

	for i := 0; i < lb.size; i++ {
		permTable[i] = make([]int, n)
		lookupTable[i] = -1
		for j := 0; j < n; j++ {
			permTable[i][j] = (offset[j] + skip[j]*i) % lb.size
		}
	}
	lb.table = lookupTable

	if len(inActives) == n {
		return // All entries are inactive, return early.
	}

	indexes := make([]int, n)
	total := 0

loop:
	for i := 0; i < lb.size; i++ {
		for j := 0; j < n; j++ {
			if slices.Contains(inActives, j) {
				continue
			}

			index := indexes[j]
			count := 0
			for index < lb.size {
				pos := permTable[index][j]
				if lookupTable[pos] < 0 {
					lookupTable[pos] = j
					count += 1
					total += 1
					if total >= lb.size {
						break loop // End generating the table.
					}
					if count >= lb.weights[j] {
						break
					}
				}
				index += 1
			}
			indexes[j] = index
		}
	}
}

func (lb *MaglevLB[T]) Get(hint int) T {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if lb.size <= 0 {
		var tmp T
		return tmp
	}

	pos := hint % lb.size
	i := lb.table[pos]
	if i < 0 {
		var tmp T
		return tmp
	}
	e := lb.entries[i]

	if !e.Active() {
		lb.update()
		i := lb.table[pos]
		if i < 0 {
			var tmp T
			return tmp
		}
		e = lb.entries[i]
	}

	return e
}

// genPrimeEuler returns a prime number grater than
// the given min value using euler's formula.
// See https://en.wikipedia.org/wiki/Lucky_numbers_of_Euler
//   - x*x + x + 41
//   - x*x - x + 41
func genPrimeEuler(min int) int {
	var val int
	i := -1
	for {
		i++
		val = i*i + i + 41 // (i*i - i + 41)
		if val >= min {
			break
		}
	}
	for {
		if isPrime(val) {
			return val
		}
		println(i, val)
		i++
		val = i*i + i + 41 // Equal to i*i - i + 41
	}
}

// isPrime returns if the given number is
// a prime number of not.
func isPrime(n int) bool {
	switch {
	case n <= 1:
		return false
	case n == 2:
		return true
	case n%2 == 0:
		return false
	}
	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
