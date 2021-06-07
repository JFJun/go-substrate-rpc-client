// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"errors"
	"fmt"
	"github.com/JFJun/go-substrate-rpc-client/v3/client"
	"github.com/JFJun/go-substrate-rpc-client/v3/types"
)

// GetStorage retreives the stored data and decodes them into the provided interface. Ok is true if the value is not
// empty.
func (s *State) GetStorage(key types.StorageKey, target interface{}, blockHash types.Hash) (ok bool, err error) {
	raw, err := s.getStorageRaw(key, &blockHash)
	if err != nil {
		return false, err
	}
	if len(*raw) == 0 {
		return false, nil
	}
	return true, types.DecodeFromBytes(*raw, target)
}

// GetStorageLatest retreives the stored data for the latest block height and decodes them into the provided interface.
// Ok is true if the value is not empty.
func (s *State) GetStorageLatest(key types.StorageKey, target interface{}) (ok bool, err error) {
	raw, err := s.getStorageRaw(key, nil)
	if err != nil {
		return false, err
	}
	if len(*raw) == 0 {
		return false, nil
	}
	return true, types.DecodeFromBytes(*raw, target)
}

// GetStorageRaw retreives the stored data as raw bytes, without decoding them
func (s *State) GetStorageRaw(key types.StorageKey, blockHash types.Hash) (*types.StorageDataRaw, error) {
	return s.getStorageRaw(key, &blockHash)
}

// GetStorageRawLatest retreives the stored data for the latest block height as raw bytes, without decoding them
func (s *State) GetStorageRawLatest(key types.StorageKey) (*types.StorageDataRaw, error) {
	return s.getStorageRaw(key, nil)
}

func (s *State) getStorageRaw(key types.StorageKey, blockHash *types.Hash) (*types.StorageDataRaw, error) {
	var res string
	err := client.CallWithBlockHash(s.client, &res, "state_getStorage", blockHash, key.Hex())
	if err != nil {
		return nil, err
	}

	bz, err := types.HexDecodeString(res)
	if err != nil {
		return nil, err
	}

	data := types.NewStorageDataRaw(bz)
	return &data, nil
}

/*
func: 因为AccountInfo这个结构老是变，所以我在解析这个结构的时候，在这里处理
author: flynn
date: 2021-06-07
*/

func (s *State) GetStorageAccountInfo(key types.StorageKey, blockHash types.Hash) (*types.AccountInfo, error) {
	raw, err := s.getStorageRaw(key, &blockHash)
	if err != nil {
		return nil, err
	}
	if len(*raw) == 0 {
		return nil, errors.New("get storage account info data raw error: raw len is 0")
	}
	var accountInfo types.AccountInfo
	switch len(*raw) {
	case 80:
		var aiwtr types.AccountInfoWithTripleRefCount
		err = types.DecodeFromBytes(*raw, &aiwtr)
		if err != nil {
			return nil, err
		}
		accountInfo.Nonce = aiwtr.Nonce
		accountInfo.Consumers = aiwtr.Consumers
		accountInfo.Providers = aiwtr.Providers
		accountInfo.Data.Free = aiwtr.Data.Free
		accountInfo.Data.FreeFrozen = aiwtr.Data.FreeFrozen
		accountInfo.Data.MiscFrozen = aiwtr.Data.MiscFrozen
		accountInfo.Data.Reserved = aiwtr.Data.Reserved
	case 76:
		var aiwpr types.AccountInfoWithProviders
		err = types.DecodeFromBytes(*raw, &aiwpr)
		if err != nil {
			return nil, err
		}
		accountInfo.Nonce = aiwpr.Nonce
		accountInfo.Consumers = aiwpr.Consumers
		accountInfo.Providers = aiwpr.Providers
		accountInfo.Data.Free = aiwpr.Data.Free
		accountInfo.Data.FreeFrozen = aiwpr.Data.FreeFrozen
		accountInfo.Data.MiscFrozen = aiwpr.Data.MiscFrozen
		accountInfo.Data.Reserved = aiwpr.Data.Reserved
	case 72:
		var aio types.AccountInfoOld
		err = types.DecodeFromBytes(*raw, &aio)
		if err != nil {
			return nil, err
		}
		accountInfo.Nonce = aio.Nonce
		accountInfo.Consumers = aio.Refcount
		accountInfo.Providers = 0
		accountInfo.Data.Free = aio.Data.Free
		accountInfo.Data.FreeFrozen = aio.Data.FreeFrozen
		accountInfo.Data.MiscFrozen = aio.Data.MiscFrozen
		accountInfo.Data.Reserved = aio.Data.Reserved
	default:
		return nil, fmt.Errorf("can not parse account info ,raw length is not standard len(80,76,72),len=%d", len(*raw))
	}
	return &accountInfo, err
}
