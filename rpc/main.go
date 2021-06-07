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

package rpc

import (
	"github.com/JFJun/go-substrate-rpc-client/v3/client"
	"github.com/JFJun/go-substrate-rpc-client/v3/rpc/author"
	"github.com/JFJun/go-substrate-rpc-client/v3/rpc/chain"
	"github.com/JFJun/go-substrate-rpc-client/v3/rpc/offchain"
	"github.com/JFJun/go-substrate-rpc-client/v3/rpc/state"
	"github.com/JFJun/go-substrate-rpc-client/v3/rpc/system"
	"github.com/JFJun/go-substrate-rpc-client/v3/types"
)

type RPC struct {
	Author   *author.Author
	Chain    *chain.Chain
	Offchain *offchain.Offchain
	State    *state.State
	System   *system.System
	client   client.Client
}

func NewRPC(cl client.Client) (*RPC, error) {
	st := state.NewState(cl)
	meta, err := st.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	opts := types.SerDeOptionsFromMetadata(meta)
	types.SetSerDeOptions(opts)

	return &RPC{
		Author:   author.NewAuthor(cl),
		Chain:    chain.NewChain(cl),
		Offchain: offchain.NewOffchain(cl),
		State:    st,
		System:   system.NewSystem(cl),
		client:   cl,
	}, nil
}
