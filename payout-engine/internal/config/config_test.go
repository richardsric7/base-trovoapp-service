package config

import (
	"strings"
	"testing"
)

func TestParseSigners(t *testing.T) {
	k := []string{
		"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		"0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		"0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
	}
	s, err := ParseSigners(strings.Join(k, ";"), "X")
	if err != nil || len(s) != 3 || s[0].Address.Hex() != "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" {
		t.Fatalf("%v %v", s, err)
	}
	if s, err := ParseSigners(strings.Join(k[:3], ","), "X"); err != nil || len(s) != 3 {
		t.Fatalf("legacy separator: %v", err)
	}
	if _, err := ParseSigners(strings.Join(k[:2], ";"), "X"); err == nil {
		t.Fatal("fewer than 3 signers must fail")
	}
	if _, err := ParseSigners("nope;"+strings.Join(k[:2], ";"), "X"); err == nil || !strings.Contains(err.Error(), "entry 1") {
		t.Fatalf("bad key: %v", err)
	}
}
