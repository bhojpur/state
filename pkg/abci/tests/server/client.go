package testsuite

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	mrand "math/rand"

	abciclient "github.com/bhojpur/state/pkg/abci/client"
	"github.com/bhojpur/state/pkg/abci/types"
	librand "github.com/bhojpur/state/pkg/libs/rand"
)

func InitChain(ctx context.Context, client abciclient.Client) error {
	total := 10
	vals := make([]types.ValidatorUpdate, total)
	for i := 0; i < total; i++ {
		pubkey := librand.Bytes(33)
		// nolint:gosec // G404: Use of weak random number generator
		power := mrand.Int()
		vals[i] = types.UpdateValidator(pubkey, int64(power), "")
	}
	_, err := client.InitChain(ctx, &types.RequestInitChain{
		Validators: vals,
	})
	if err != nil {
		fmt.Printf("Failed test: InitChain - %v\n", err)
		return err
	}
	fmt.Println("Passed test: InitChain")
	return nil
}

func Commit(ctx context.Context, client abciclient.Client, hashExp []byte) error {
	res, err := client.Commit(ctx)
	data := res.Data
	if err != nil {
		fmt.Println("Failed test: Commit")
		fmt.Printf("error while committing: %v\n", err)
		return err
	}
	if !bytes.Equal(data, hashExp) {
		fmt.Println("Failed test: Commit")
		fmt.Printf("Commit hash was unexpected. Got %X expected %X\n", data, hashExp)
		return errors.New("commitTx failed")
	}
	fmt.Println("Passed test: Commit")
	return nil
}

func FinalizeBlock(ctx context.Context, client abciclient.Client, txBytes [][]byte, codeExp []uint32, dataExp []byte) error {
	res, _ := client.FinalizeBlock(ctx, &types.RequestFinalizeBlock{Txs: txBytes})
	for i, tx := range res.TxResults {
		code, data, log := tx.Code, tx.Data, tx.Log
		if code != codeExp[i] {
			fmt.Println("Failed test: FinalizeBlock")
			fmt.Printf("FinalizeBlock response code was unexpected. Got %v expected %v. Log: %v\n",
				code, codeExp, log)
			return errors.New("FinalizeBlock error")
		}
		if !bytes.Equal(data, dataExp) {
			fmt.Println("Failed test:  FinalizeBlock")
			fmt.Printf("FinalizeBlock response data was unexpected. Got %X expected %X\n",
				data, dataExp)
			return errors.New("FinalizeBlock  error")
		}
	}
	fmt.Println("Passed test: FinalizeBlock")
	return nil
}

func CheckTx(ctx context.Context, client abciclient.Client, txBytes []byte, codeExp uint32, dataExp []byte) error {
	res, _ := client.CheckTx(ctx, &types.RequestCheckTx{Tx: txBytes})
	code, data, log := res.Code, res.Data, res.Log
	if code != codeExp {
		fmt.Println("Failed test: CheckTx")
		fmt.Printf("CheckTx response code was unexpected. Got %v expected %v. Log: %v\n",
			code, codeExp, log)
		return errors.New("checkTx")
	}
	if !bytes.Equal(data, dataExp) {
		fmt.Println("Failed test: CheckTx")
		fmt.Printf("CheckTx response data was unexpected. Got %X expected %X\n",
			data, dataExp)
		return errors.New("checkTx")
	}
	fmt.Println("Passed test: CheckTx")
	return nil
}
