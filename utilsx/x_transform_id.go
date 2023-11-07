package utilsx

import (
	"log"
	"sync"

	"github.com/sqids/sqids-go"
)

const (
	TRANSFORMER_ID_ALPHABET  string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	TRANSFORMER_ID_MINLENGTH uint8  = 10
)

type idTransformer struct {
	alphabet   string
	minLength  uint8
	sequenceNo uint64

	transformer *sqids.Sqids
	once        sync.Once
}

type idTransformerInterface interface {
	SetAlphabet(string)
	SetMinLength(uint8)
	SetSequenceNo(uint64)
	Encode(id uint64) string
	Decode(id string) uint64
}

// NewIdTransformer creates a new instance of idTransformerInterface.
//
// It initializes a new idTransformer struct with the TRANSFORMER_ID_ALPHABET
// and TRANSFORMER_ID_MINLENGTH constants as its properties.
// Returns a pointer to the newly created idTransformer struct.
func NewIdTransformer() idTransformerInterface {
	return &idTransformer{
		alphabet:  TRANSFORMER_ID_ALPHABET,
		minLength: TRANSFORMER_ID_MINLENGTH,
	}
}

// sequenceMaxNumber calculates the maximum number that can be obtained by multiplying
// the given number with the sequence of numbers from i to 1.
//
// i: The number from which the sequence starts.
// Returns: The maximum number obtained by multiplying the sequence of numbers.
func (it *idTransformer) sequenceMaxNumber(i uint64) uint64 {
	var number uint64 = 1
	for i > 0 {
		number *= i
		i--
	}
	return number
}

// SetSequenceNo sets the sequence number for the idTransformer.
//
// It takes a uint64 value representing the sequence number to be set.
// It does not return anything.
func (it *idTransformer) SetSequenceNo(no uint64) {
	if it.transformer != nil {
		log.Fatal("idTransformer has already been initialized")
	}
	it.sequenceNo = no
}

// SetAlphabet sets the alphabet for the idTransformer.
//
// It takes a string parameter named alphabet.
// It does not return anything.
func (it *idTransformer) SetAlphabet(alphabet string) {
	if it.transformer != nil {
		log.Fatal("idTransformer has already been initialized")
	}
	it.alphabet = alphabet
}

// SetMinLength sets the minimum length of the idTransformer.
//
// minLength: The minimum length to be set.
func (it *idTransformer) SetMinLength(minLength uint8) {
	if it.transformer != nil {
		log.Fatal("idTransformer has already been initialized")
	}
	it.minLength = minLength
}

// Encode encodes the given ID using the idTransformer struct.
//
// It takes a uint64 ID as a parameter and returns a string.
func (it *idTransformer) Encode(id uint64) string {
	s, _ := it.getSqids()
	encodeStr, err := s.Encode([]uint64{id})
	if err != nil {
		log.Printf("invalid id: %d", id)
		return ""
	}
	return encodeStr
}

// Decode decodes the given ID and returns the corresponding uint64 value.
//
// It takes a string parameter 'id' which represents the ID to be decoded.
// The function returns a uint64 value which is the decoded number.
func (it *idTransformer) Decode(id string) uint64 {
	s, _ := it.getSqids()
	numbers := s.Decode(id)
	if len(numbers) == 0 {
		log.Printf("invalid id: %s", id)
		return 0
	}
	return numbers[0]
}

// getSqids returns the sqids.Sqids instance and an error.
//
// It initializes the transformer if it hasn't been done already,
// using the provided alphabet and minLength options.
// Then, it returns the transformer instance and nil error.
func (it *idTransformer) getSqids() (*sqids.Sqids, error) {
	it.once.Do(func() {

		originAlphabet := it.alphabet
		if it.sequenceNo > it.sequenceMaxNumber(uint64(len(originAlphabet))) {
			log.Fatal("字典序列号大于最大生成数量")
		}
		var customAlphabet []byte
		for len(originAlphabet) > 0 {
			alphabetLength := uint64(len(originAlphabet))
			index := it.sequenceNo % alphabetLength
			letter := []byte(originAlphabet)[index]
			originAlphabet = originAlphabet[:index] + originAlphabet[index+1:]
			customAlphabet = append(customAlphabet, letter)
			it.sequenceNo /= alphabetLength
		}
		it.alphabet = string(customAlphabet) + originAlphabet

		it.transformer, _ = sqids.New(sqids.Options{
			Alphabet:  it.alphabet,
			MinLength: it.minLength,
		})
	})
	return it.transformer, nil
}
