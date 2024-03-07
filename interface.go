package fooche

import "time"

type ICache interface {
	Set(k string, v []byte, ttl time.Duration) error
	Has(k string) bool
	Get(k string) (v []byte, err error)
	Delete(k string)
	String() string
}
