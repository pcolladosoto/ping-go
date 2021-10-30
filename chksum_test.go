package main

import "testing"

func TestChecksum(t *testing.T) {
	var tests = []struct {
		in   []byte
		want uint16
	}{
		{
			ICMPMsg{
				Type:     8,
				Code:     0,
				Checkusm: 0,
				Id:       5,
				SeqNum:   10,
				Payload:  []byte("Hello there!"),
			}.Marshal(),
			0x9502,
		},
		{
			ICMPMsg{
				Type:     8,
				Code:     0,
				Checkusm: 0,
				Id:       5,
				SeqNum:   10,
				Payload:  []byte("0"),
			}.Marshal(),
			0xc7f0,
		},
		{
			ICMPMsg{
				Type:     8,
				Code:     0,
				Checkusm: 0,
				Id:       5,
				SeqNum:   10,
				Payload:  []byte("Am I in an ICMP message O_o?"),
			}.Marshal(),
			0x8c8c,
		},
	}

	for _, test := range tests {
		if got := CheckSum(test.in); got != test.want {
			t.Errorf("convert(%v) = %q; wanted %v", test.in, got, test.want)
		}
	}
}
