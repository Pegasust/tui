package main

import (
	"fmt"
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
			input: "./#//old/way/to:write",
			want:  Spec{FlakeRef: "./", Registry: flake.BrandedRegistry, Cell: "old", Block: "way", Target: "to", Action: "write"},
		},
		{
			input: "git+ssh://some-gl-profile?ref=bleed&dir=subdir#//old/way/to:write",
			want:  Spec{FlakeRef: "git+ssh://some-gl-profile?ref=bleed&dir=subdir", Registry: flake.BrandedRegistry, Cell: "old", Block: "way", Target: "to", Action: "write"},
		},
		{
			input: "github:nixos/nixpkgs#//registry/has/default:build",
			want:  Spec{FlakeRef: "github:nixos/nixpkgs", Registry: flake.BrandedRegistry, Cell: "registry", Block: "has", Target: "default", Action: "build"},
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
		{
			input: "fh:ryantm/agenix/*#//devops/containers/service-foo:deploy",
			// TODO: this is better of to be put under golden/snap test
			want: Spec{FlakeRef: "https://api.flakehub.com/f/ryantm/agenix/*.tar.gz", Registry: flake.BrandedRegistry, Cell: "devops", Block: "containers", Target: "service-foo", Action: "deploy"},
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

func TestParseUnspec(t *testing.T) {
	tests := []string{
		"github:nixos/nixpkgs//registry/has/default:build",
	}
	for _, test := range tests {
		got, err := parseSpec(test)
		if err == nil {
			fmt.Printf("Failed to assert error for '%s', partial: %v\n", test, got)
		}
	}
}
