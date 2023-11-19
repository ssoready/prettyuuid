package prettyuuid_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/ssoready/prettyuuid"
)

func TestFormat(t *testing.T) {
	testCases := []struct {
		Name   string
		Format prettyuuid.Format
		UUID   string
		Want   string
	}{
		{
			Name:   "round-trip zero uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdef"},
			UUID:   "00000000-0000-0000-0000-000000000000",
			Want:   "00000000000000000000000000000000",
		},
		{
			Name:   "round-trip one uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdef"},
			UUID:   "00000000-0000-0000-0000-000000000001",
			Want:   "00000000000000000000000000000001",
		},
		{
			Name:   "round-trip max uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdef"},
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Want:   "ffffffffffffffffffffffffffffffff",
		},
		{
			Name:   "round-trip sample uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdef"},
			UUID:   "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
			Want:   "f81d4fae7dec11d0a76500a0c91e6bf6",
		},
		{
			Name:   "alphanumeric zero uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "00000000-0000-0000-0000-000000000000",
			Want:   "0000000000000000000000000",
		},
		{
			Name:   "alphanumeric one uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "00000000-0000-0000-0000-000000000001",
			Want:   "0000000000000000000000001",
		},
		{
			Name:   "alphanumeric max uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Want:   "f5lxx1zz5pnorynqglhzmsp33",
		},
		{
			Name:   "alphanumeric sample uuid",
			Format: prettyuuid.Format{Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
			Want:   "eoswzolg3bsx0zn8otq1p8oom",
		},
		{
			Name:   "alphanumeric zero uuid with prefix",
			Format: prettyuuid.Format{Prefix: "prefix_", Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "00000000-0000-0000-0000-000000000000",
			Want:   "prefix_0000000000000000000000000",
		},
		{
			Name:   "alphanumeric one uuid with prefix",
			Format: prettyuuid.Format{Prefix: "prefix_", Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "00000000-0000-0000-0000-000000000001",
			Want:   "prefix_0000000000000000000000001",
		},
		{
			Name:   "alphanumeric max uuid with prefix",
			Format: prettyuuid.Format{Prefix: "prefix_", Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Want:   "prefix_f5lxx1zz5pnorynqglhzmsp33",
		},
		{
			Name:   "alphanumeric sample uuid with prefix",
			Format: prettyuuid.Format{Prefix: "prefix_", Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"},
			UUID:   "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
			Want:   "prefix_eoswzolg3bsx0zn8otq1p8oom",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			id := uuid.MustParse(tt.UUID)
			got := tt.Format.Format(id)
			if d := cmp.Diff(tt.Want, got); d != "" {
				t.Fatalf("%v.Format(%v) != %v (-want +got)\n%s", tt.Format, tt.UUID, tt.Want, d)
			}

			parsed, err := tt.Format.Parse(tt.Want)
			if err != nil {
				t.Fatalf("%v.Parse(%v) returned unexpected err: %v", tt.Format, tt.Want, err)
			}

			if d := cmp.Diff(id, uuid.UUID(parsed)); d != "" {
				t.Fatalf("%v.Parse(%v) != %v (-want +got)\n%s", tt.Format, tt.Want, id, d)
			}
		})
	}
}
