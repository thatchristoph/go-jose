/*-
 * Copyright 2014 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package josecipher

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/binary"
)

// DeriveECDHES derives a shared encryption key using ECDH/ConcatKDF as described in JWE/JWA.
func DeriveECDHES(alg string, apuData, apvData []byte, priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, size int) []byte {
	// algId, partyUInfo, partyVInfo inputs must be prefixed with the length
	algData := []byte(alg)
	algID := make([]byte, 4)
	binary.BigEndian.PutUint32(algID, uint32(len(algData)))
	algID = append(algID, algData...)

	ptyUInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(ptyUInfo, uint32(len(apuData)))
	ptyUInfo = append(ptyUInfo, apuData...)

	ptyVInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(ptyVInfo, uint32(len(apvData)))
	ptyVInfo = append(ptyVInfo, apvData...)

	// suppPubInfo is the encoded length of the output size in bits
	supPubInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(supPubInfo, uint32(size)*8)

	z, _ := priv.PublicKey.Curve.ScalarMult(pub.X, pub.Y, priv.D.Bytes())
	reader := NewConcatKDF(crypto.SHA256, z.Bytes(), algID, ptyUInfo, ptyVInfo, supPubInfo, []byte{})

	key := make([]byte, size)

	// Read on the KDF will never fail
	_, _ = reader.Read(key)
	return key
}
