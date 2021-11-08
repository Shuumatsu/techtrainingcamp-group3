package bloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
	"techtrainingcamp-group3/config"
)

var User *bloom.BloomFilter
var Envelope *bloom.BloomFilter

func init() {
	User = bloom.NewWithEstimates(6e8, 0.001)
	Envelope = bloom.NewWithEstimates(config.MaxSnatchAmount*6e8, 0.001)
}
