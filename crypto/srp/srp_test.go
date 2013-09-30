// Copyright 2013 Tad Glines
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package srp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
)

var groups []string = []string{
	"rfc5054.1024",
	"rfc5054.1536",
	"rfc5054.2048",
	"rfc5054.3072",
	"rfc5054.4096",
	"rfc5054.6144",
	"rfc5054.8192",
	"stanford.1024",
	"stanford.1536",
	"stanford.2048",
	"stanford.3072",
	"stanford.4096",
	"stanford.6144",
	"stanford.8192",
}

var passwords []string = []string{
	"0",
	"a",
	"password",
	"This Is A Long Password",
	"This is a really long password a;lsdfkjauiwjenfasueifxl3847tq8374y(*&^JHG&*^$.kjbh()&*^KJG",
}

type hashFunc func() hash.Hash

var hashes []hashFunc = []hashFunc{
	sha1.New,
	sha256.New,
	sha512.New,
}

func testSRP(t *testing.T, group string, h func() hash.Hash, username, password []byte) {
	srp, _ := NewSRP(group, h)
	cs := srp.NewClientSession(username, password)
	salt, v, err := srp.ComputeVerifier(password)
	if err != nil {
		t.Fatal(err)
	}
	ss := srp.NewServerSession(username, salt, v)

	_, err = cs.ComputeKey(salt, ss.GetB())
	if err != nil {
		t.Fatal(err)
	}

	_, err = ss.ComputeKey(cs.GetA())
	if err != nil {
		t.Fatal(err)
	}

	cauth := cs.ComputeAuthenticator()
	if !ss.VerifyClientAuthenticator(cauth) {
		t.Fatal("Client Authenticator is not valid")
	}

	sauth := ss.ComputeAuthenticator(cauth)
	if !ss.VerifyClientAuthenticator(sauth) {
		t.Fatal("Server Authenticator is not valid")
	}
}

func TestSRPSimple(t *testing.T) {
	for _, g := range groups {
		for _, h := range hashes {
			for _, p := range passwords {
				testSRP(t, g, h, []byte("test"), []byte(p))
			}
		}
	}
}
