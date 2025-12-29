package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/unix-world/smartgoext/cloud/vcard"

	"github.com/unix-world/smartgoplus/cloud/dav/carddav"
)

func TestFilter(t *testing.T) {
	newAO := func(str string) carddav.AddressObject {
		card, err := vcard.NewDecoder(strings.NewReader(str)).Decode()
		if err != nil {
			t.Fatal(err)
		}
		return carddav.AddressObject{
			Card: card,
		}
	}

	alice := newAO(`BEGIN:VCARD
VERSION:4.0
UID:urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b1
FN;PID=1.1:Alice Gopher
N:Gopher;Alice;;;
EMAIL;PID=1.1:alice@example.com
CLIENTPIDMAP:1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0551
END:VCARD`)

	bob := newAO(`BEGIN:VCARD
VERSION:4.0
UID:urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b2
FN;PID=1.1:Bob Gopher
N:Gopher;Bob;;;
EMAIL;PID=1.1:bob@example.com
CLIENTPIDMAP:1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0552
END:VCARD`)

	carla := newAO(`BEGIN:VCARD
VERSION:4.0
UID:urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b3
FN;PID=1.1:Carla Gopher
N:Gopher;Carla;;;
EMAIL;PID=1.1:carla@example.com
CLIENTPIDMAP:1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0553
END:VCARD`)
	carlaFiltered := newAO(`BEGIN:VCARD
VERSION:4.0
UID:urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b3
EMAIL;PID=1.1:carla@example.com
END:VCARD`)

	for _, tc := range []struct {
		name  string
		query *carddav.AddressBookQuery
		addrs []carddav.AddressObject
		want  []carddav.AddressObject
		err   error
	}{
		{
			name:  "nil-query",
			query: nil,
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{alice, bob, carla},
		},
		{
			name: "no-limit-query",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{alice, bob, carla},
		},
		{
			name: "limit-1-query",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				Limit: 1,
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{alice},
		},
		{
			name: "limit-4-query",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				Limit: 4,
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{alice, bob, carla},
		},
		{
			name: "email-match",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "carla"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{carla},
		},
		{
			name: "email-match-any",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{
							{Text: "carla@example"},
							{Text: "alice@example"},
						},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{alice, carla},
		},
		{
			name: "email-match-all",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				PropFilters: []carddav.PropFilter{{
					Name: vcard.FieldEmail,
					TextMatches: []carddav.TextMatch{
						{Text: ""},
					},
				}},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{alice, bob, carla},
		},
		{
			name: "email-no-match",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					AllProp: true,
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.org"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{},
		},
		{
			name: "email-match-filter-properties",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldVersion,
						vcard.FieldUID,
						vcard.FieldEmail,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "carla"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{carlaFiltered},
		},
		{
			name: "email-match-filter-properties-always-returns-version",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldUID,
						vcard.FieldEmail,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "carla"}},
					},
				},
			},
			addrs: []carddav.AddressObject{alice, bob, carla},
			want:  []carddav.AddressObject{carlaFiltered},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := carddav.Filter(tc.query, tc.addrs)
			switch {
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %q\nwant=%q", got, want)
				}
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error:\ngot= %+v\nwant=%+v", err, tc.err)
			case err == nil && tc.err == nil:
				if got, want := got, tc.want; !reflect.DeepEqual(got, want) {
					t.Fatalf("invalid filter values:\ngot= %+v\nwant=%+v", got, want)
				}
			}
		})
	}
}

func TestMatch(t *testing.T) {
	newAO := func(str string) carddav.AddressObject {
		card, err := vcard.NewDecoder(strings.NewReader(str)).Decode()
		if err != nil {
			t.Fatal(err)
		}
		return carddav.AddressObject{
			Card: card,
		}
	}

	alice := newAO(`BEGIN:VCARD
VERSION:4.0
UID:urn:uuid:4fbe8971-0bc3-424c-9c26-36c3e1eff6b1
FN;PID=1.1:Alice Gopher
N:Gopher;Alice;;;
EMAIL;PID=1.1:alice@example.com
CLIENTPIDMAP:1;urn:uuid:53e374d9-337e-4727-8803-a1e9c14e0556
END:VCARD`)

	for _, tc := range []struct {
		name  string
		query *carddav.AddressBookQuery
		addr  carddav.AddressObject
		want  bool
		err   error
	}{
		{
			name:  "nil-query",
			query: nil,
			addr:  alice,
			want:  true,
		},
		{
			name: "match-email-contains",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-email-equals-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{
							Text:      "alice@example.com",
							MatchType: carddav.MatchEquals,
						}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-email-equals-not",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{
							Text:      "example.com",
							MatchType: carddav.MatchEquals,
						}},
					},
				},
			},
			addr: alice,
			want: false,
		},
		{
			name: "match-email-equals-ok-negate",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{
							Text:            "bob@example.com",
							NegateCondition: true,
							MatchType:       carddav.MatchEquals,
						}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-email-starts-with-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{
							Text:      "alice@",
							MatchType: carddav.MatchStartsWith,
						}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-email-ends-with-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{
							Text:      "com",
							MatchType: carddav.MatchEndsWith,
						}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-email-ends-with-not",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{
							Text:      ".org",
							MatchType: carddav.MatchEndsWith,
						}},
					},
				},
			},
			addr: alice,
			want: false,
		},
		{
			name: "match-name-contains-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldName,
						TextMatches: []carddav.TextMatch{{
							Text:      "Alice",
							MatchType: carddav.MatchContains,
						}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-name-contains-all-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldName,
						Test: carddav.FilterAllOf,
						TextMatches: []carddav.TextMatch{
							{
								Text:      "Alice",
								MatchType: carddav.MatchContains,
							},
							{
								Text:      "Gopher",
								MatchType: carddav.MatchContains,
							},
						},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-name-contains-all-prop-not",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				FilterTest: carddav.FilterAllOf,
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldName,
						TextMatches: []carddav.TextMatch{{
							Text:      "Alice",
							MatchType: carddav.MatchContains,
						}},
					},
					{
						Name: vcard.FieldName,
						TextMatches: []carddav.TextMatch{{
							Text:      "GopherXXX",
							MatchType: carddav.MatchContains,
						}},
					},
				},
			},
			addr: alice,
			want: false,
		},
		{
			name: "match-name-contains-all-text-match-not",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name: vcard.FieldName,
						Test: carddav.FilterAllOf,
						TextMatches: []carddav.TextMatch{
							{
								Text:      "Alice",
								MatchType: carddav.MatchContains,
							},
							{
								Text:      "GopherXXX",
								MatchType: carddav.MatchContains,
							},
						},
					},
				},
			},
			addr: alice,
			want: false,
		},
		{
			name: "missing-prop-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
						"XXX-not-THERE", // but AllProp is false.
					},
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "match-all-prop-ok",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
					AllProp: true,
				},
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addr: alice,
			want: true,
		},
		{
			name: "invalid-query-filter",
			query: &carddav.AddressBookQuery{
				DataRequest: carddav.AddressDataRequest{
					Props: []string{
						vcard.FieldFormattedName,
						vcard.FieldEmail,
						vcard.FieldUID,
					},
				},
				FilterTest: "XXX-invalid-filter",
				PropFilters: []carddav.PropFilter{
					{
						Name:        vcard.FieldEmail,
						TextMatches: []carddav.TextMatch{{Text: "example.com"}},
					},
				},
			},
			addr: alice,
			err:  fmt.Errorf("unknown query filter test \"XXX-invalid-filter\""),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := carddav.Match(tc.query, &tc.addr)
			switch {
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %q\nwant=%q", got, want)
				}
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error:\ngot= %+v\nwant=%+v", err, tc.err)
			case err == nil && tc.err == nil:
				if got, want := got, tc.want; got != want {
					t.Fatalf("invalid match value: got=%v, want=%v", got, want)
				}
			}
		})
	}
}
