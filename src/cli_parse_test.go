package main

import (
	"testing"

	"github.com/paisano-nix/paisano/flake"
)

func TestParseSpec(t *testing.T) {
	tests := []struct {
		input string
		want  Spec
	}{
		{
			input: "github:nixos/nixpkgs/nixos-unstable#__std//path/to/something:action",
			want:  Spec{FlakeRef: "github:nixos/nixpkgs/nixos-unstable", Registry: "__std", Cell: "path", Block: "to", Target: "something", Action: "action"},
		},
		{
			input: "//old/way/to:write",
			want:  Spec{FlakeRef: ".", Registry: flake.BrandedRegistry, Cell: "old", Block: "way", Target: "to", Action: "write"},
		},
		{
			input: "github:nixos/nixpkgs#//python3/packages/hypercorn:build",
			want:  Spec{FlakeRef: "github:nixos/nixpkgs", Registry: flake.BrandedRegistry, Cell: "python3", Block: "packages", Target: "hypercorn", Action: "build"},
		},
		{
			input: "#__rebranded//devops/containers/service-foo:deploy",
			want:  Spec{FlakeRef: ".", Registry: "__rebranded", Cell: "devops", Block: "containers", Target: "service-foo", Action: "deploy"},
		},
		{
			input: "fh:ryantm/agenix/*#__rebranded//devops/containers/service-foo:deploy",
			// TODO: this is better of to be put under golden/snap test
			want: Spec{FlakeRef: "https://api.flakehub.com/f/ryantm/agenix/*.tar.gz", Registry: "__rebranded", Cell: "devops", Block: "containers", Target: "service-foo", Action: "deploy"},
		},
	}
	for _, test := range tests {
		got, err := parseSpec(test.input)
		if err != nil {
			t.Errorf("Regex failed: %v", err)
		}
		if *got != test.want {
			t.Errorf("Got %v; Want %v", *got, test.want)
		}
	}
}
