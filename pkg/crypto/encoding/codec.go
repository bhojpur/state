package encoding

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
	"fmt"

	"github.com/bhojpur/state/internal/jsontypes"
	cryptopb "github.com/bhojpur/state/pkg/api/v1/crypto"
	"github.com/bhojpur/state/pkg/crypto"
	"github.com/bhojpur/state/pkg/crypto/ed25519"
	"github.com/bhojpur/state/pkg/crypto/secp256k1"
	"github.com/bhojpur/state/pkg/crypto/sr25519"
)

func init() {
	jsontypes.MustRegister((*cryptopb.PublicKey)(nil))
	jsontypes.MustRegister((*cryptopb.PublicKey_Ed25519)(nil))
	jsontypes.MustRegister((*cryptopb.PublicKey_Secp256K1)(nil))
}

// PubKeyToProto takes crypto.PubKey and transforms it to a protobuf Pubkey
func PubKeyToProto(k crypto.PubKey) (cryptopb.PublicKey, error) {
	var kp cryptopb.PublicKey
	switch k := k.(type) {
	case ed25519.PubKey:
		kp = cryptopb.PublicKey{
			Sum: &cryptopb.PublicKey_Ed25519{
				Ed25519: k,
			},
		}
	case secp256k1.PubKey:
		kp = cryptopb.PublicKey{
			Sum: &cryptopb.PublicKey_Secp256K1{
				Secp256K1: k,
			},
		}
	case sr25519.PubKey:
		kp = cryptopb.PublicKey{
			Sum: &cryptopb.PublicKey_Sr25519{
				Sr25519: k,
			},
		}
	default:
		return kp, fmt.Errorf("toproto: key type %v is not supported", k)
	}
	return kp, nil
}

// PubKeyFromProto takes a protobuf Pubkey and transforms it to a crypto.Pubkey
func PubKeyFromProto(k cryptopb.PublicKey) (crypto.PubKey, error) {
	switch k := k.Sum.(type) {
	case *cryptopb.PublicKey_Ed25519:
		if len(k.Ed25519) != ed25519.PubKeySize {
			return nil, fmt.Errorf("invalid size for PubKeyEd25519. Got %d, expected %d",
				len(k.Ed25519), ed25519.PubKeySize)
		}
		pk := make(ed25519.PubKey, ed25519.PubKeySize)
		copy(pk, k.Ed25519)
		return pk, nil
	case *cryptopb.PublicKey_Secp256K1:
		if len(k.Secp256K1) != secp256k1.PubKeySize {
			return nil, fmt.Errorf("invalid size for PubKeySecp256k1. Got %d, expected %d",
				len(k.Secp256K1), secp256k1.PubKeySize)
		}
		pk := make(secp256k1.PubKey, secp256k1.PubKeySize)
		copy(pk, k.Secp256K1)
		return pk, nil
	case *cryptopb.PublicKey_Sr25519:
		if len(k.Sr25519) != sr25519.PubKeySize {
			return nil, fmt.Errorf("invalid size for PubKeySr25519. Got %d, expected %d",
				len(k.Sr25519), sr25519.PubKeySize)
		}
		pk := make(sr25519.PubKey, sr25519.PubKeySize)
		copy(pk, k.Sr25519)
		return pk, nil
	default:
		return nil, fmt.Errorf("fromproto: key type %v is not supported", k)
	}
}
