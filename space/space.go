package space

import (
	"fmt"
	"crypto/md5"
	"encoding/base64"
)

type Entry [][]byte

func MakeEntry(entries [][]byte) Entry {
	return entries
}

type TupleSpace struct{
	entryList []Entry
	// TODO: The tuple space has more stuff
}

func Hash(data []byte) string {
	hash := md5.Sum(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (ts *TupleSpace) Write(tuple Entry, ok *bool) error {
	fmt.Println("TupleSpace.Write Called!")
	fmt.Println("Tuple provided:", tuple)
	
	fmt.Println("Hash of all the entry fields:")
	for i, v := range tuple {
		fmt.Println(i, ":", Hash(v))
	}
	fmt.Println()
	
	// Check the type (caller can send anything)
	// TODO: Is this really necessary?
	
	// Done!
	*ok = true
	return nil
}