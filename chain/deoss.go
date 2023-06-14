/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package chain

import (
	"log"

	"github.com/CESSProject/cess-go-sdk/core/pattern"
	"github.com/CESSProject/cess-go-sdk/core/utils"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
)

// QueryDeossPeerPublickey
func (c *ChainSDK) QueryDeossPeerPublickey(pubkey []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utils.RecoverError(err))
		}
	}()
	var data pattern.PeerId

	if !c.GetChainState() {
		return nil, pattern.ERR_RPC_CONNECTION
	}

	key, err := types.CreateStorageKey(c.metadata, pattern.OSS, pattern.OSS, pubkey)
	if err != nil {
		return nil, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return nil, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return nil, pattern.ERR_RPC_EMPTY_VALUE
	}
	return []byte(string(data[:])), nil
}

// QueryDeossPeerPublickey
func (c *ChainSDK) QueryDeossPeerIdList() ([]string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utils.RecoverError(err))
		}
	}()
	var result []string

	if !c.GetChainState() {
		return nil, pattern.ERR_RPC_CONNECTION
	}

	key := createPrefixedKey(pattern.OSS, pattern.OSS)
	keys, err := c.api.RPC.State.GetKeysLatest(key)
	if err != nil {
		return nil, errors.Wrap(err, "[GetKeysLatest]")
	}

	set, err := c.api.RPC.State.QueryStorageAtLatest(keys)
	if err != nil {
		return nil, errors.Wrap(err, "[QueryStorageAtLatest]")
	}

	for _, elem := range set {
		for _, change := range elem.Changes {
			var data pattern.PeerId
			if err := codec.Decode(change.StorageData, &data); err != nil {
				continue
			}
			result = append(result, base58.Encode([]byte(string(data[:]))))
		}
	}
	return result, nil
}

func (c *ChainSDK) QuaryAuthorizedAcc(puk []byte) (types.AccountID, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utils.RecoverError(err))
		}
	}()
	var data types.AccountID

	if !c.GetChainState() {
		return data, pattern.ERR_RPC_CONNECTION
	}

	acc, err := types.NewAccountID(puk)
	if err != nil {
		return data, errors.Wrap(err, "[NewAccountID]")
	}

	b, err := codec.Encode(*acc)
	if err != nil {
		return data, errors.Wrap(err, "[EncodeToBytes]")
	}

	key, err := types.CreateStorageKey(c.metadata, pattern.OSS, pattern.AUTHORITYLIST, b)
	if err != nil {
		return data, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return data, pattern.ERR_RPC_EMPTY_VALUE
	}
	return data, nil
}

func (c *ChainSDK) QuaryAuthorizedAccount(puk []byte) (string, error) {
	acc, err := c.QuaryAuthorizedAcc(puk)
	if err != nil {
		return "", err
	}
	return utils.EncodePublicKeyAsCessAccount(acc[:])
}

func (c *ChainSDK) CheckSpaceUsageAuthorization(puk []byte) (bool, error) {
	grantor, err := c.QuaryAuthorizedAcc(puk)
	if err != nil {
		if err.Error() == pattern.ERR_Empty {
			return false, nil
		}
		return false, err
	}
	account_chain, _ := utils.EncodePublicKeyAsCessAccount(grantor[:])
	account_local := c.GetSignatureAcc()

	return account_chain == account_local, nil
}
