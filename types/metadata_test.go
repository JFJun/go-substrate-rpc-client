package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_ParseMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV13Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	d, _ := json.Marshal(meta.AsMetadataV13)
	fmt.Println(string(d))
}
