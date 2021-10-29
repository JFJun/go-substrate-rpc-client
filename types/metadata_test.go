package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_ParseMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	v14 := meta.AsMetadataV14
	d, _ := json.Marshal(v14)
	fmt.Println(string(d))
}

func TestMetadataV14FindCallIndex(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	callIdx, err := meta.FindCallIndex("Balances.transfer")
	if err != nil {
		panic(err)
	}
	fmt.Println(callIdx)
}
func TestMetadataV14FindEventNamesForEventID(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	id := EventID{}
	id[0] = 5
	id[1] = 2
	mod, event, err := meta.FindEventNamesForEventID(id)
	if err != nil {
		panic(err)
	}
	fmt.Println(mod, event)
}

func TestMetadataV14FindStorageEntryMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	st, err := meta.FindStorageEntryMetadata("System", "Account")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(st)
}

func TestMetadataV14ExistsModuleMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	s := meta.ExistsModuleMetadata("System")

	fmt.Println(s)
}
