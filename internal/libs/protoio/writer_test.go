package protoio_test

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
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/state/internal/libs/protoio"
	v1 "github.com/bhojpur/state/pkg/api/v1/types"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/types"
)

func aVote(t testing.TB) *types.Vote {
	t.Helper()
	var stamp, err = time.Parse(types.TimeFormat, "2017-12-25T03:00:01.234Z")
	require.NoError(t, err)

	return &types.Vote{
		Type:      v1.SignedMsgType(byte(v1.PrevoteType)),
		Height:    12345,
		Round:     2,
		Timestamp: stamp,
		BlockID: types.BlockID{
			Hash: crypto.Checksum([]byte("blockID_hash")),
			PartSetHeader: types.PartSetHeader{
				Total: 1000000,
				Hash:  crypto.Checksum([]byte("blockID_part_set_header_hash")),
			},
		},
		ValidatorAddress: crypto.AddressHash([]byte("validator_address")),
		ValidatorIndex:   56789,
	}
}

type excludedMarshalTo struct {
	msg proto.Message
}

func (emt *excludedMarshalTo) ProtoMessage() {}
func (emt *excludedMarshalTo) String() string {
	return emt.msg.String()
}
func (emt *excludedMarshalTo) Reset() {
	emt.msg.Reset()
}
func (emt *excludedMarshalTo) Marshal() ([]byte, error) {
	return proto.Marshal(emt.msg)
}

var _ proto.Message = (*excludedMarshalTo)(nil)

var sink interface{}

func BenchmarkMarshalDelimitedWithMarshalTo(b *testing.B) {
	msgs := []proto.Message{
		aVote(b).ToProto(),
	}
	benchmarkMarshalDelimited(b, msgs)
}

func BenchmarkMarshalDelimitedNoMarshalTo(b *testing.B) {
	msgs := []proto.Message{
		&excludedMarshalTo{aVote(b).ToProto()},
	}
	benchmarkMarshalDelimited(b, msgs)
}

func benchmarkMarshalDelimited(b *testing.B, msgs []proto.Message) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, msg := range msgs {
			blob, err := protoio.MarshalDelimited(msg)
			require.Nil(b, err)
			sink = blob
		}
	}

	if sink == nil {
		b.Fatal("Benchmark did not run")
	}

	// Reset the sink.
	sink = (interface{})(nil)
}
