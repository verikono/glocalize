package cmd

import "testing"

func TestRuntimeFlagsRoundTrip(t *testing.T) {
	in := runtimeFlags{Port: 6660, Hot: true, DevDir: "/srv/riot"}
	data, err := marshalFlags(in)
	if err != nil {
		t.Fatalf("marshalFlags() error = %v", err)
	}
	out, err := unmarshalFlags(data)
	if err != nil {
		t.Fatalf("unmarshalFlags() error = %v", err)
	}
	if out != in {
		t.Fatalf("round trip = %+v, want %+v", out, in)
	}
}
