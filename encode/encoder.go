package encode

import (
	"bytes"
	"encoding/gob"
	"crypto/md5"
	"encoding/base64"
)

func EncodeBytes(stuff interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(stuff)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeBytes(stuff []byte, ptr interface{}) error {
	buf := bytes.NewBuffer(stuff)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(ptr)

	if err != nil {
		return err
	}

	return nil
}


func Hash(data []byte) string {
	hash := md5.Sum(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}
