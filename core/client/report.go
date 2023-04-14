/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package client

import (
	"github.com/CESSProject/sdk-go/core/chain"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func (c *Cli) ReportFile(roothash []string) (string, []chain.FileHash, error) {
	var hashs = make([]chain.FileHash, len(roothash))
	for i := 0; i < len(roothash); i++ {
		for j := 0; j < len(roothash[i]); j++ {
			hashs[i][j] = types.U8(roothash[i][j])
		}
	}
	return c.Chain.SubmitFileReport(hashs)
}
