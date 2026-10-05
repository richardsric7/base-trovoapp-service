package vaultclient

import (
	"reflect"
	"testing"
)

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "current format",
			raw:  "0xaaa;0xbbb;0xccc",
			want: []string{"0xaaa", "0xbbb", "0xccc"},
		},
		{
			name: "commas are part of an entry",
			raw:  "0xaaa,0xbbb",
			want: []string{"0xaaa,0xbbb"},
		},
		{
			name: "single entry, no separators either way",
			raw:  "0xaaa",
			want: []string{"0xaaa"},
		},
		{
			name: "a pasted mnemonic's internal commas stay inside their own entry",
			raw:  "word1,word2,word3;0xbbb",
			want: []string{"word1,word2,word3", "0xbbb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitCSV(tt.raw)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitCSV(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
		})
	}
}
