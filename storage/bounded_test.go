package storage

import (
	"strconv"
	"testing"
)

const (
	TEST_MIN_SIZE = 5
	TEST_MED_SIZE = 10
	TEST_MAX_SIZE = 15
	TEST_MIN_CAP  = 3
	TEST_MED_CAP  = 2
	TEST_MAX_CAP  = 1
)

var (
	initMap = map[int]int{
		TEST_MIN_SIZE: TEST_MIN_CAP,
		TEST_MED_SIZE: TEST_MED_CAP,
		TEST_MAX_SIZE: TEST_MAX_CAP,
	}
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestOverall(t *testing.T) {
	storage := NewBoundedStorage(initMap)
	if storage.Size() != 0 {
		t.Errorf("expected storage to be empty, but got %v", storage.Size())
	}

	key := "foo"
	value := "bar"
	err := storage.Put(key, []byte(value))
	if err != nil {
		t.Error("expected put to succeed, but failed")
	}

	if storage.Size() != 1 {
		t.Errorf("expected storage to have size 1, but got %v", storage.Size())
	}

	if v, ok := storage.Get(key); !ok {
		t.Error("expected get to succeed, but failed")
	} else if string(v) != value {
		t.Errorf("expected get on key %s to be: %s, but was: %s", key, value, string(v))
	}

	storage.Remove(key)
	if storage.Size() != 0 {
		t.Errorf("expected storage to be empty, but got %v", storage.Size())
	}
}

func TestFull(t *testing.T) {
	storage := NewBoundedStorage(initMap)
	key := "foo"
	value := "bar"
	for i := 0; i < TEST_MIN_CAP; i++ {
		storage.Put(key, []byte(value))
	}

	if err := storage.Put(key, []byte(value)); err != nil {
		t.Error("expected put to succeed, since  all keys were equal it should not be full")
	} else if storage.Size() != 1 {
		t.Errorf("expected size to be 1, was: %v", storage.Size())
	}

	for i := 0; i < TEST_MIN_CAP-1; i++ {
		v := []byte(value)
		v = append(v, byte(i+1))
		storage.Put(key+strconv.FormatInt(int64(i+1), 10), v)
	}

	if err := storage.Put(key+"last", []byte(value)); err != nil {
		t.Errorf("expected put to succeed since it would store in a different size bucket, but failed")
	} else if storage.Size() != TEST_MIN_CAP+1 {
		t.Errorf("expected size to be %d, was: %d", TEST_MIN_CAP+1, storage.Size())
	} else if storage.keyToBucketMap[key+"last"] != TEST_MED_SIZE+4 /*This +4 is accounting for the four bytes for the value length*/ {
		t.Errorf("expected key to be in bucket %d, was in: %d", TEST_MED_SIZE+4, storage.keyToBucketMap[key+"last"])
	}

	// Rewrite same key in different size bucket
	err := storage.Put(key, []byte("barbar"))
	if err != nil {
		t.Error("expected put to succeed since we are inserting in another bucket")
	} else if storage.Size() != TEST_MIN_CAP+1 {
		t.Errorf("same key in different size bucket should've been removed. Expected size to be %d, was: %d", TEST_MIN_CAP+1, storage.Size())
	}

	storage.Put("focus", []byte("abcdefghijklmn"))

	if storage.bucketToIndexMap[TEST_MIN_SIZE+4].Size() != TEST_MIN_CAP-1 {
		t.Errorf("expected min bucket size to be %d, was: %d", TEST_MIN_CAP-1, storage.bucketToIndexMap[TEST_MIN_SIZE+4].Size())
	} else if storage.bucketToIndexMap[TEST_MED_SIZE+4].Size() != TEST_MED_CAP {
		t.Errorf("expected med bucket size to be full (%d), was: %d", TEST_MED_CAP, storage.bucketToIndexMap[TEST_MED_SIZE+4].Size())
	} else if storage.bucketToIndexMap[TEST_MAX_SIZE+4].Size() != TEST_MAX_CAP {
		t.Errorf("expected max bucket size to be full (%d), was: %d", TEST_MAX_CAP, storage.bucketToIndexMap[TEST_MAX_SIZE+4].Size())
	}

	err = storage.Put("foo", []byte("abcdefghijklmn"))
	if err == nil {
		t.Errorf("expected put to fail since MAX_CAP bucket is full")
	}

	storage.Remove("focus")

	err = storage.Put("foo", []byte("abcdefghijklmn"))
	if err != nil {
		t.Errorf("expected put of key foo in MAX_CAP bucket to succeed, but it failed")
	}

	err = storage.Put("foo1", []byte("abcdefghi"))
	if err != nil {
		t.Errorf("expected put of key foo1 in MED_CAP bucket to succeed, but it failed")
	}

	if storage.bucketToIndexMap[TEST_MIN_SIZE+4].Size() != TEST_MIN_CAP-2 {
		t.Errorf("expected min bucket size to be %d, was: %d", TEST_MIN_CAP-2, storage.bucketToIndexMap[TEST_MIN_SIZE+4].Size())
	} else if storage.bucketToIndexMap[TEST_MED_SIZE+4].Size() != TEST_MED_CAP {
		t.Errorf("expected min bucket size to be full (%d), was: %d", TEST_MED_CAP, storage.bucketToIndexMap[TEST_MED_SIZE+4].Size())
	} else if storage.bucketToIndexMap[TEST_MAX_SIZE+4].Size() != TEST_MAX_CAP {
		t.Errorf("expected min bucket size to be full (%d), was: %d", TEST_MAX_CAP, storage.bucketToIndexMap[TEST_MAX_SIZE+4].Size())
	}

	if v, ok := storage.Get(key); !ok {
		t.Errorf("expected to find foo key in storage, but didn't")
	} else if string(v) != "abcdefghijklmn" {
		t.Errorf("expected value of key foo to be %s, but was %s", "abcdefghijklmn", string(v))
	}
}
