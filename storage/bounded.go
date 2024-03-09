package storage

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/joaovictorsl/gollections/maps"
	"golang.org/x/exp/slices"
)

// A storage that has a fixed capacity and can store values of different sizes.
type BoundedStorage struct {
	// Maps a key to the bucket it is stored in
	keyToBucketMap map[string]int
	// Maps a bucket to the index where the value is storage in the data slice
	bucketToIndexMap map[int]*maps.OccupationMap[string, int]
	// The data slice
	data []byte
	// The different sizes of buckets
	sizes []int
	// The total capacity of objects the storage can store
	cap int
}

// Creates a new BoundedStorage.
//
// # Params
//   - sizeToCapMap: a map of size in bytes and capacity (how many objects with this size or less can be stored).
//
// # Returns
//
// A new BoundedStorage.
func NewBoundedStorage(sizeToCapMap map[int]int) *BoundedStorage {
	if len(sizeToCapMap) == 0 {
		log.Fatal("Invalid args for new BoundedStorage. sizeToCapMap must not be empty.")
	}

	// Calc total size and amount of possible keys
	totalSize := 0
	totalCap := 0
	biMap := make(map[int]*maps.OccupationMap[string, int])
	sizes := make([]int, 0, len(sizeToCapMap))
	lastIdx := 0
	for size, cap := range sizeToCapMap {
		if size == 0 || cap == 0 {
			continue
		}
		// Adds length bytes
		size += 4

		totalSize += size * cap
		totalCap += cap

		places := make([]int, cap)
		for i := 0; i < cap; i++ {
			places[i] = lastIdx
			lastIdx += size
		}

		biMap[size] = maps.NewOccupationMap[string, int](places...)
		sizes = append(sizes, size)
	}

	slices.Sort(sizes)

	bs := &BoundedStorage{
		bucketToIndexMap: biMap,
		sizes:            sizes,
		data:             make([]byte, totalSize),
		keyToBucketMap:   make(map[string]int, totalCap),
		cap:              totalCap,
	}

	return bs
}

func (bs *BoundedStorage) Put(k string, v []byte) (err error) {
	prevBucket, hadPrevBucket := bs.keyToBucketMap[k]
	bucket := -1
	// Find smallest bucket that fits the value and has empty place
	for _, s := range bs.sizes {
		fits := len(v)+4 <= s
		// Fits in bucket and bucket has empty place
		if (fits && !bs.bucketToIndexMap[s].IsFull()) ||
			// Fits in bucket, which may be full, but prevBucket is the current bucket
			(fits && hadPrevBucket && prevBucket == s) {
			bucket = s
			break
		}
	}

	if bucket == -1 {
		// Nowhere to store value
		return fmt.Errorf("value doesn't fit in any bucket")
	}

	if hadPrevBucket {
		// Free same key from previous bucket
		// A key should be unique across all buckets
		bs.bucketToIndexMap[prevBucket].Free(k)
	}

	p, _ := bs.bucketToIndexMap[bucket].Occupy(k)

	binary.LittleEndian.PutUint32(bs.data[p:], uint32(len(v)))
	copy(bs.data[p+4:], v)

	bs.keyToBucketMap[k] = bucket

	return nil
}

func (bs *BoundedStorage) Get(k string) (v []byte, ok bool) {
	bucket, ok := bs.keyToBucketMap[k]
	if !ok {
		return nil, false
	}

	if idx, ok := bs.bucketToIndexMap[bucket].Get(k); ok {
		vLen := int(binary.LittleEndian.Uint32(bs.data[idx : idx+4]))
		return bs.data[idx+4 : idx+4+vLen], true
	}

	return nil, false
}

func (bs *BoundedStorage) Remove(k string) {
	bucket, ok := bs.keyToBucketMap[k]
	if ok {
		bs.bucketToIndexMap[bucket].Free(k)
		delete(bs.keyToBucketMap, k)
	}
}

func (bs *BoundedStorage) Size() int {
	currSize := 0
	for _, om := range bs.bucketToIndexMap {
		currSize += om.Size()
	}

	return currSize
}

func (bs *BoundedStorage) Capacity() int {
	return bs.cap
}
