package encode

import(
	"testing"
	"runtime/debug"
)

type Struct struct{
	Number int
	Name string
	Real float32
}

func TestEncode(t *testing.T) {
	s := Struct{10, "name", 0.5}
	_, err := EncodeBytes(s)

	if err != nil {
		t.Error("Error while encoding", err)
		debug.PrintStack()
	}
}

func TestDecode(t *testing.T) {
	data := []byte{49, 255, 129, 3, 1, 1, 6, 83, 116, 114, 117, 99, 116, 1, 255, 130, 0, 1, 3, 1, 6, 78, 117, 109, 98, 101, 114, 1, 4, 0, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 4, 82, 101, 97, 108, 1, 8, 0, 0, 0, 15, 255, 130, 1, 20, 1, 4, 110, 97, 109, 101, 1, 254, 224, 63, 0}

	var s Struct
	expected  := Struct{10, "name", 0.5}

	err := DecodeBytes(data, &s)

	if err != nil {
		t.Error("Error while decoding", err)
		debug.PrintStack()
	} 

	if s != expected {
		t.Error(s, "is different from ", expected)
		debug.PrintStack()
	}
}