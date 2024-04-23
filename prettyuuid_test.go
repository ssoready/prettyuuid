package prettyuuid_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/ssoready/prettyuuid"
)

func Example() {
	format := prettyuuid.MustNewFormat("invoice_", "0123456789abcdefghijklmnopqrstuvwxyz")

	fmt.Println(format.Format(uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")))

	id, _ := format.Parse("invoice_eoswzolg3bsx0zn8otq1p8oom")
	fmt.Println(uuid.UUID(id).String())
	// Output:
	// invoice_eoswzolg3bsx0zn8otq1p8oom
	// f81d4fae-7dec-11d0-a765-00a0c91e6bf6
}

func TestNewFormat(t *testing.T) {
	testCases := []struct {
		Alphabet string
		WantErr  string
	}{
		{
			Alphabet: "",
			WantErr:  `alphabet must have len >= 2: ""`,
		},
		{
			Alphabet: "a",
			WantErr:  `alphabet must have len >= 2: "a"`,
		},
		{
			Alphabet: "aa",
			WantErr:  `alphabet must not contain duplicate chars: "aa"`,
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%q", tt.Alphabet), func(t *testing.T) {
			_, err := prettyuuid.NewFormat("", tt.Alphabet)
			if d := cmp.Diff(tt.WantErr, err.Error()); d != "" {
				t.Fatalf("NewFormat(%v) did not return expected err (-want +got):\n%s", tt.Alphabet, d)
			}
		})
	}
}

func TestFormat_Parse(t *testing.T) {
	testCases := []struct {
		Name    string
		Format  prettyuuid.Format
		ID      string
		WantErr string
	}{
		{
			Name:    "bad length",
			Format:  prettyuuid.MustNewFormat("", "0123456789abcdef"),
			ID:      "00000000000000000000000000000000x",
			WantErr: `"00000000000000000000000000000000x" does not have expected length 32`,
		},
		{
			Name:    "bad prefix",
			Format:  prettyuuid.MustNewFormat("prefix_", "0123456789abcdefghijklmnopqrstuvwxyz"),
			ID:      "xxxxxx_0000000000000000000000000",
			WantErr: `"xxxxxx_0000000000000000000000000" does not have expected prefix "prefix_"`,
		},
		{
			Name:    "bad char",
			Format:  prettyuuid.MustNewFormat("", "0123456789abcdef"),
			ID:      "0000000000000000000000000000000x",
			WantErr: `"0000000000000000000000000000000x" contains illegal char at position 31`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := tt.Format.Parse(tt.ID)
			if d := cmp.Diff(tt.WantErr, err.Error()); d != "" {
				t.Fatalf("%v.Parse(%v) did not return expected err (-want +got):\n%s", tt.Format, tt.ID, d)
			}
		})
	}
}

func TestFormat_roundtrip(t *testing.T) {
	testCases := []struct {
		Name   string
		Format prettyuuid.Format
		UUID   string
		Want   string
	}{
		{
			Name:   "round-trip zero uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdef"),
			UUID:   "00000000-0000-0000-0000-000000000000",
			Want:   "00000000000000000000000000000000",
		},
		{
			Name:   "round-trip one uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdef"),
			UUID:   "00000000-0000-0000-0000-000000000001",
			Want:   "00000000000000000000000000000001",
		},
		{
			Name:   "round-trip max uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdef"),
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Want:   "ffffffffffffffffffffffffffffffff",
		},
		{
			Name:   "round-trip sample uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdef"),
			UUID:   "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
			Want:   "f81d4fae7dec11d0a76500a0c91e6bf6",
		},
		{
			Name:   "alphanumeric zero uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "00000000-0000-0000-0000-000000000000",
			Want:   "0000000000000000000000000",
		},
		{
			Name:   "alphanumeric one uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "00000000-0000-0000-0000-000000000001",
			Want:   "0000000000000000000000001",
		},
		{
			Name:   "alphanumeric max uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Want:   "f5lxx1zz5pnorynqglhzmsp33",
		},
		{
			Name:   "alphanumeric sample uuid",
			Format: prettyuuid.MustNewFormat("", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
			Want:   "eoswzolg3bsx0zn8otq1p8oom",
		},
		{
			Name:   "alphanumeric zero uuid with prefix",
			Format: prettyuuid.MustNewFormat("prefix_", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "00000000-0000-0000-0000-000000000000",
			Want:   "prefix_0000000000000000000000000",
		},
		{
			Name:   "alphanumeric one uuid with prefix",
			Format: prettyuuid.MustNewFormat("prefix_", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "00000000-0000-0000-0000-000000000001",
			Want:   "prefix_0000000000000000000000001",
		},
		{
			Name:   "alphanumeric max uuid with prefix",
			Format: prettyuuid.MustNewFormat("prefix_", "0123456789abcdefghijklmnopqrstuvwxyz"),
			UUID:   "ffffffff-ffff-ffff-ffff-ffffffffffff",
			Want:   "prefix_f5lxx1zz5pnorynqglhzmsp33",
		},
		{
			Name:   "alphanumeric sample uuid with prefix",
			Format: prettyuuid.MustNewFormat("prefix_", "0123456789abcdefghijklmnopqrstuvwxyz"),
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

func FuzzFormat_roundtrip(f *testing.F) {
	f.Fuzz(func(t *testing.T, prefix, alphabet string, fuzzID []byte) {
		// we have to take a byte slice because go fuzz does not take arrays
		if len(fuzzID) != 16 {
			return
		}
		id := (*[16]byte)(fuzzID)

		format, err := prettyuuid.NewFormat(prefix, alphabet)
		if err != nil {
			return // we don't care about invalid formats
		}

		formatted := format.Format(*id)
		parsed, err := format.Parse(formatted)
		if err != nil {
			t.Fatalf("%v.Parse(%v) err: %v", format, formatted, err)
		}
		if d := cmp.Diff(*id, parsed); d != "" {
			t.Fatalf("%v.Parse(Format(%v)) does not round-trip: (-want +got)\n%s", format, id, d)
		}
	})
}
