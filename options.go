/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package sdkgo

import (
	"time"

	"github.com/CESSProject/cess-go-sdk/chain"
)

// ConnectRpcAddrs configuration rpc address
func ConnectRpcAddrs(s []string) Option {
	return func(cfg *Config) error {
		cfg.Rpc = s
		return nil
	}
}

// Mnemonic configures the mnemonic of the signature account
func Mnemonic(mnemonic string) Option {
	return func(cfg *Config) error {
		cfg.Mnemonic = mnemonic
		return nil
	}
}

// TransactionTimeout configures the waiting timeout for a transaction
func TransactionTimeout(timeout time.Duration) Option {
	return func(cfg *Config) error {
		if timeout < chain.BlockInterval {
			cfg.Timeout = chain.BlockInterval
		} else {
			cfg.Timeout = timeout
		}
		return nil
	}
}

// Name configuration sdk name
func Name(name string) Option {
	return func(cfg *Config) error {
		cfg.Name = name
		return nil
	}
}
