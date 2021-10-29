package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/JFJun/go-substrate-rpc-client/v3/scale"
	"hash"
	"strings"
	"sync"
)

type MetadataV14 struct {
	Lookup     PortableRegistry
	ldLk       sync.Mutex
	LookUpData map[int64]*Si1Type
	Pallets    []PalletMetadataV14
	Extrinsic  ExtrinsicMetadataV14
}

func (d *MetadataV14) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	for _, mod := range d.Pallets {
		if !mod.HasCalls {
			continue
		}
		if string(mod.Name) != s[0] {
			continue
		}
		callType := mod.Calls.Type

		for _, lookUp := range d.Lookup {
			if lookUp.Id.Int64() == callType.Int64() {
				if len(lookUp.Type.Def.Variant.Variants) > 0 {
					for _, vars := range lookUp.Type.Def.Variant.Variants {
						if string(vars.Name) == s[1] {
							return CallIndex{uint8(mod.Index), uint8(vars.Index)}, nil
						}
					}
				}
			}
		}
	}
	return CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)
}

func (d *MetadataV14) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	for _, mod := range d.Pallets {
		if !mod.HasEvents {
			continue
		}
		if uint8(mod.Index) != eventID[0] {
			continue
		}
		eventType := mod.Events.Type.Int64()

		for _, lookUp := range d.Lookup {
			if lookUp.Id.Int64() == eventType {
				if len(lookUp.Type.Def.Variant.Variants) > 0 {
					for _, vars := range lookUp.Type.Def.Variant.Variants {
						if uint8(vars.Index) == eventID[1] {
							return mod.Name, vars.Name, nil
						}
					}
				}
			}
		}
	}
	return "", "", fmt.Errorf("module index %v out of range", eventID[0])
}

func (d *MetadataV14) FindStorageEntryMetadata(module string, fn string) (StorageEntryMetadata, error) {
	for _, mod := range d.Pallets {
		if !mod.HasStorage {
			continue
		}
		if string(mod.Storage.Prefix) != module {
			continue
		}
		for _, s := range mod.Storage.Items {
			if string(s.Name) != fn {
				continue
			}
			return s, nil
		}
		return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
	}
	return nil, fmt.Errorf("module %v not found in metadata", module)
}

func (d *MetadataV14) ExistsModuleMetadata(module string) bool {
	for _, mod := range d.Pallets {
		if string(mod.Name) == module {
			return true
		}
	}
	return false
}

func (d *MetadataV14) Decode(decoder scale.Decoder) error {
	var err error
	err = decoder.Decode(&d.Lookup)
	if err != nil {
		return err
	}
	// 处理lookUp
	d.LookUpData = make(map[int64]*Si1Type)
	d.ldLk.Lock()
	for _, lookUp := range d.Lookup {
		d.LookUpData[lookUp.Id.Int64()] = &lookUp.Type
	}
	d.ldLk.Unlock()

	err = decoder.Decode(&d.Pallets)
	if err != nil {
		return err
	}

	err = decoder.Decode(&d.Extrinsic)
	if err != nil {
		return err
	}
	return nil
}

type PortableRegistry GenericPortableRegistry
type GenericPortableRegistry []PortableTypeV14

type PalletMetadataV14 struct {
	Name       Text
	HasStorage bool
	Storage    PalletStorageMetadataV14
	HasCalls   bool
	Calls      PalletCallMetadataV14
	HasEvents  bool
	Events     PalletEventMetadataV14
	Constants  []PalletConstantMetadataV14
	HasErrors  bool
	Errors     PalletErrorMetadataV14
	Index      U8
}

func (m *PalletMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}
	err = decoder.Decode(&m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = decoder.Decode(&m.Storage)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = decoder.Decode(&m.Calls)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = decoder.Decode(&m.Events)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}
	err = decoder.Decode(&m.HasErrors)
	if err != nil {
		return err
	}
	if m.HasErrors {
		err = decoder.Decode(&m.Errors)
		if err != nil {
			return err
		}
	}
	return decoder.Decode(&m.Index)
}

func (m PalletMetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = encoder.Encode(m.Storage)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = encoder.Encode(m.Calls)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = encoder.Encode(m.Events)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.Constants)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Errors)
	if err != nil {
		return err
	}

	return encoder.Encode(m.Index)
}

/*
Storage
*/

type PalletStorageMetadataV14 struct {
	Prefix Text
	Items  []StorageEntryMetadataV14
}

func (d *PalletStorageMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&d.Prefix)
	if err != nil {
		return err
	}
	return decoder.Decode(&d.Items)
}

type StorageEntryMetadataV14 struct {
	Name          Text
	Modifier      StorageEntryModifierV14
	Type          StorageEntryTypeV14
	Fallback      Bytes
	Documentation []Text
}

func (s StorageEntryMetadataV14) IsPlain() bool {
	return s.Type.IsPlainType
}

func (s StorageEntryMetadataV14) IsMap() bool {
	return false
}

func (s StorageEntryMetadataV14) IsDoubleMap() bool {
	return false
}
func (s StorageEntryMetadataV14) IsNMap() bool {
	return s.Type.IsMap
}
func (s StorageEntryMetadataV14) GetHashers() ([]hash.Hash, error) {
	var (
		hashes []hash.Hash
	)
	if s.Type.IsMap {
		for _, hasher := range s.Type.AsMap.Hasher {
			h, err := hasher.HashFunc()
			if err != nil {
				return nil, err
			}
			hashes = append(hashes, h)
		}
	}
	return hashes, nil
}

type StorageEntryModifierV14 struct {
	IsOptional bool // 0
	IsDefault  bool // 1
	IsRequired bool // 2
}

func (s *StorageEntryModifierV14) Decode(decoder scale.Decoder) error {
	var t uint8
	err := decoder.Decode(&t)
	if err != nil {
		return err
	}

	switch t {
	case 0:
		s.IsOptional = true
	case 1:
		s.IsDefault = true
	case 2:
		s.IsRequired = true
	default:
		return fmt.Errorf("received unexpected storage function modifier type %v", t)
	}
	return nil
}

type StorageEntryTypeV14 struct {
	IsPlainType bool
	AsPlainType Si1LookupTypeId
	IsMap       bool
	AsMap       MapTypeV14
}

func (d *StorageEntryTypeV14) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		d.IsPlainType = true
		err = decoder.Decode(&d.AsPlainType)
		if err != nil {
			return err
		}
	case 1:
		d.IsMap = true
		err = decoder.Decode(&d.AsMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("StorageFunctionTypeV14 is not support this type: %d", b)
	}
	return nil
}

type MapTypeV14 struct {
	Hasher  []StorageHasherV10
	KeysId  Si1LookupTypeId
	ValueId Si1LookupTypeId
}

/*
Call
*/

type PalletCallMetadataV14 struct {
	Type Si1LookupTypeId
}

/*
Event
*/

type PalletEventMetadataV14 struct {
	Type Si1LookupTypeId
}

/*
Constant
*/

type PalletConstantMetadataV14 struct {
	Name  Text
	Type  Si1LookupTypeId
	Value Bytes
	Docs  []Text
}

/*
Error
*/

type PalletErrorMetadataV14 struct {
	Type Si1LookupTypeId
}

type ExtrinsicMetadataV14 struct {
	Type             Si1LookupTypeId
	Version          U8
	SignedExtensions []SignedExtensionMetadataV14
}

type SignedExtensionMetadataV14 struct {
	Identifier       Text
	Type             Si1LookupTypeId
	AdditionalSigned Si1LookupTypeId
}

/*
func: 为了适应bifrost-go这个包而添加以下功能
author: flynn
date: 2021/10/29
*/

type IMetadataFunc interface {
	getCallIndex(moduleName, fn string) (callIdx string, err error)
	findNameByCallIndex(callIdx string) (moduleName, fn string, err error)
	getConstants(modName, constantsName string) (constantsType string, constantsValue []byte, err error)
}

func (d *MetadataV14) getCallIndex(moduleName, fn string) (string, error) {
	idx, err := d.FindCallIndex(fmt.Sprintf("%s.%s", moduleName, fn))
	if err != nil {
		return "", err
	}
	return idx.String(), nil

}

func (d *MetadataV14) findNameByCallIndex(callIdx string) (string, string, error) {
	if len(callIdx) != 4 {
		return "", "", fmt.Errorf("call index length is not equal 4: length: %d", len(callIdx))
	}
	data, err := hex.DecodeString(callIdx)
	if err != nil {
		return "", "", fmt.Errorf("call index is not hex string")
	}
	for _, mod := range d.Pallets {
		if !mod.HasCalls {
			continue
		}
		if uint8(mod.Index) == data[0] {
			callType := mod.Calls.Type.Int64()
			d.ldLk.Lock()
			call := d.LookUpData[callType]
			if call == nil {
				return "", "", fmt.Errorf("%s do not have this call id: %d", mod.Name, data[1])
			}
			if len(call.Def.Variant.Variants) == 0 {
				return "", "", fmt.Errorf("%s  call.Def.Variant.Variants len is 0", mod.Name)
			}
			for _, vars := range call.Def.Variant.Variants {
				if uint8(vars.Index) == data[1] {
					return string(mod.Name), string(vars.Name), nil
				}
			}
			d.ldLk.Unlock()
		}
	}
	return "", "", errors.New("do not find")
}

func (d *MetadataV14) getConstants(modName, constantsName string) (constantsType string, constantsValue []byte, err error) {
	for _, mod := range d.Pallets {
		if string(mod.Name) == modName {
			for _, constants := range mod.Constants {
				if string(constants.Name) == constantsName {
					constantsTypeId := constants.Type.Int64()
					d.ldLk.Lock()
					siType := d.LookUpData[constantsTypeId]
					if siType == nil {
						return "", nil, fmt.Errorf("%s.%s constants type is nil ptr", mod.Name, constants.Name)
					}
					constantsType = siType.Def.Primitive.Value
					d.ldLk.Unlock()
					constantsValue = constants.Value
					return constantsType, constantsValue, nil
				}
			}
		}
	}
	return "", nil, fmt.Errorf("do not find this constants,moduleName=%s,"+
		"constantsName=%s", modName, constantsName)
}
