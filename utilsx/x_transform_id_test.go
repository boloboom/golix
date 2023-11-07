package utilsx

import "testing"

func TestTransformId(t *testing.T) {
	var testId uint64 = 1

	it1 := NewIdTransformer()
	it1.SetAlphabet("0123456789")
	it1.SetSequenceNo(1)
	encodeId1 := it1.Encode(testId)
	decodeId1 := it1.Decode(encodeId1)
	t.Log(decodeId1, encodeId1)

	it2 := NewIdTransformer()
	it2.SetSequenceNo(2)
	encodeId2 := it2.Encode(testId)
	decodeId2 := it2.Decode(encodeId2)
	t.Log(decodeId2, encodeId2)

	it3 := NewIdTransformer()
	it3.SetSequenceNo(3)
	encodeId3 := it3.Encode(testId)
	decodeId3 := it3.Decode(encodeId3)
	t.Log(decodeId3, encodeId3)
}
