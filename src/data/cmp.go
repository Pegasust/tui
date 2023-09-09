package data

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/paisano-nix/paisano/flake"
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
)

type Completion struct {
	local CompeIndex
}

type CompeIndex struct {
	cells   []string
	blocks  map[string][]string
	targets map[string]map[string][]string
	actions map[string]map[string]map[string][]string
}

func (root Root) indexRegistry() CompeIndex {
	var cells = []string{}
	var blocks = map[string][]string{}
	var targets = map[string]map[string][]string{}
	var actions = map[string]map[string]map[string][]string{}
	for _, c := range root.Cells {
		blocks[c.Name] = []string{}
		targets[c.Name] = map[string][]string{}
		actions[c.Name] = map[string]map[string][]string{}
		cells = append(cells, c.Name, "cell")
		for _, b := range c.Blocks {
			targets[c.Name][b.Name] = []string{}
			actions[c.Name][b.Name] = map[string][]string{}
			blocks[c.Name] = append(blocks[c.Name], b.Name, "block")
			for _, t := range b.Targets {
				actions[c.Name][b.Name][t.Name] = []string{}
				targets[c.Name][b.Name] = append(targets[c.Name][b.Name], t.Name, t.Description())
				for _, a := range t.Actions {
					actions[c.Name][b.Name][t.Name] = append(
						actions[c.Name][b.Name][t.Name],
						a.Name,
						a.Description(),
					)
				}
			}
		}
	}
	return CompeIndex{
		cells:   cells,
		blocks:  blocks,
		targets: targets,
		actions: actions,
	}
}

func StdPathCompe(compe *CompeIndex, c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	// start with <tab>; no typing
	case 0:
		return carapace.ActionValuesDescribed(
			compe.cells...,
		).Invoke(c).Suffix("/").ToA().Style(
			style.Of(style.Bold, style.Carapace.Highlight(1)))
	// only a single / typed
	case 1:
		return carapace.ActionValuesDescribed(
			compe.cells...,
		).Invoke(c).Prefix("/").Suffix("/").ToA()
	// start typing cell
	case 2:
		return carapace.ActionValuesDescribed(
			compe.cells...,
		).Invoke(c).Suffix("/").ToA().Style(
			style.Carapace.Highlight(1))
	// start typing block
	case 3:
		return carapace.ActionValuesDescribed(
			compe.blocks[c.Parts[2]]...,
		).Invoke(c).Suffix("/").ToA().Style(
			style.Carapace.Highlight(2))
	// start typing target
	case 4:
		return carapace.ActionMultiParts(":", func(d carapace.Context) carapace.Action {
			switch len(d.Parts) {
			// start typing target
			case 0:
				return carapace.ActionValuesDescribed(
					compe.targets[c.Parts[2]][c.Parts[3]]...,
				).Invoke(c).Suffix(":").ToA().Style(
					style.Carapace.Highlight(3))
				// start typing action
			case 1:
				return carapace.ActionValuesDescribed(
					compe.actions[c.Parts[2]][c.Parts[3]][d.Parts[0]]...,
				).Invoke(c).ToA()
			default:
				return carapace.ActionValues()
			}
		})
	default:
		return carapace.ActionValues()
	}
}

var (
	FlakeHubProto    = "fh"
	GitHubProto      = "github"
	GitLabProto      = "gitlab"
	SourceHutProto   = "sourcehut"
	NixRegistryProto = "flake"
	FileProto        = "file"
	PathProto        = "path"
)

func flakehubFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMessage("TODO: implement")
}

func githubFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMessage("TODO: implement")
}

func gitlabFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {

	return carapace.ActionMessage("TODO: implement")
}

func sourcehutFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {

	return carapace.ActionMessage("TODO: implement")
}

func nixRegistryFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {

	return carapace.ActionMessage("TODO: implement")
}

func gitQueryCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	// Requires ctx.Value to point at `?|query=value`
	q, err := url.ParseQuery(ctx.Value)
	if err != nil {
		return carapace.ActionMessage("Programmer err, invalid query: %v", err)
	}
	prefix := ""
	if len(q) != 0 {
		prefix = "&"
	}
	desc := []string{
		"dir", "subdirectory",
		"rev", "branch name, tag name, or tagged commit-ish",
		"ref", "commit-ish reference",
	}

	descs := make([]string, 0, len(desc))
	for i := 0; i < len(desc); i += 2 {
		if !q.Has(desc[i]) {
			descs = append(descs, desc[i], desc[i+1])
		}
	}
	return carapace.Batch(
		carapace.ActionValuesDescribed(descs...),
	).ToA().Prefix(prefix)
}

func gitFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMultiParts("?", func(c carapace.Context) carapace.Action {
		switch len(ctx.Parts) {
		case 0:
			return carapace.Batch(
				carapace.ActionValuesDescribed(
					"git://", "git over path",
					"git+path://", "git over path",
					"git+https://", "git over https (Gitea or such)",
					"git+ssh://", "git over ssh",
				),
				carapace.ActionMessage("TODO: complete git compe"),
			).ToA()
		case 1:
			return gitQueryCompe(compe, ctx)
		default:
			return carapace.ActionMessage(
				"Err: Should only have <= 1 '?'",
			).Style(style.Red)
		}
	}).Invoke(ctx).ToA()
}

func fileFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMultiParts("?", func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		case 0:
			// TODO: Filter out only flake directory
			return carapace.ActionDirectories()
		case 1:
			return gitQueryCompe(compe, ctx)
		default:
			return carapace.ActionMessage(
				"Not conforming to URL+query. Should have <= 1 '?', got %v",
				len(c.Parts),
			).Style(style.Red)
		}
	})
}

type compeSource func(*CompeIndex, carapace.Context) carapace.Action

func CarapaceRootCmdCompletion(cmd *carapace.Carapace, argv0 string) {
	// completes: '//cell/block/target:action'
	cmd.PositionalCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {

			// TODO: support remote caching, maybe in XDG cache
			local := flake.LocalFlakeRegistry()
			cache, key, _, _, err := local.LoadFlakeCmd()
			if err != nil {
				return carapace.ActionMessage(fmt.Sprintf("%v\n", err))
			}
			cached, _, err := cache.GetBytes(*key)
			var root *Root
			if err == nil {
				root, err = LoadJson(bytes.NewReader(cached))
				if err != nil {
					return carapace.ActionMessage(fmt.Sprintf("%v\n", err))
				}
			} else {
				return carapace.ActionMessage(fmt.Sprintf("No completion cache: please initialize by running '%[1]s re-cache'.", argv0))
			}
			index := root.indexRegistry()
			return carapace.ActionMultiParts("#", func(c carapace.Context) carapace.Action {
				switch len(c.Parts) {
				case 0:
					if strings.HasPrefix(c.Value, "//") {
						std := StdPathCompe(&index, c)
						return carapace.Batch(
							std,
							carapace.ActionValues("DEBUG: entering std compe"),
						).ToA()
					}
					return carapace.Batch(
						carapace.ActionMessage("DEBUG: Does this block other values?").Style(style.Red),
						fileFlakeCompe(&index, c),
						carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {

							switch len(c.Parts) {
							case 0:
								return carapace.Batch(
									carapace.ActionValuesDescribed(
										FlakeHubProto, "Flakehub",
										GitHubProto, "GitHub",
										GitLabProto, "GitLab",
										SourceHutProto, "SourceHut",
										NixRegistryProto, "Nix registry",
									).Suffix(":"),
									carapace.ActionValuesDescribed(
										FileProto, "Local file",
										PathProto, "Local path",
									).Suffix("://"),
								).ToA()
							case 1:
								dispatch := map[string](compeSource){
									// NOTE: API specs
									FlakeHubProto:    flakehubFlakeCompe,
									GitHubProto:      githubFlakeCompe,
									GitLabProto:      gitlabFlakeCompe,
									SourceHutProto:   sourcehutFlakeCompe,
									NixRegistryProto: nixRegistryFlakeCompe,

									FileProto: fileFlakeCompe,
									PathProto: fileFlakeCompe,
								}
								call := dispatch[c.Parts[0]]
								if call == nil {
									return carapace.ActionMessage("Unknown %s", c.Parts[0])
								}
								return call(&index, c)
							default:
								return carapace.ActionMessage("Ensure Nix support, docs: https://nixos.org/manual/nix/stable/command-ref/new-cli/nix3-flake.html#types")
							}
						}),
						gitFlakeCompe(&index, c),
					).ToA().Invoke(c).ToA()
					// return carapace.ActionValuesDescribed(
					// 	"./", "Local Flake",
					// 	"", "Local Flake",
					// 	"github:", "Github Flake",
					// 	"fh:", "Flakehub",
					// 	"git+ssh://", "Git Flake over ssh",
					// 	"git+https://", "Git Flake over https",
					// 	"git+http://", "Git Flake over http",
					// ).Invoke(c).ToA()
				case 1:
					return carapace.ActionMultiParts("//", func(c carapace.Context) carapace.Action {
						switch len(c.Parts) {
						case 0:
							// paisano registry
							well_known := map[string]string{
								flake.BrandedRegistry: "Custom branded registry",
								"__std":               "github:divnix/std",
							}
							dedup_registries := make([]string, 0, len(well_known)*2)
							for name, desc := range well_known {
								dedup_registries = append(dedup_registries, name, desc)
							}
							return carapace.ActionValuesDescribed(
								dedup_registries...,
							).NoSpace().Invoke(c).ToA()
						case 1:
							// std
							return StdPathCompe(&index, c)
						default:
							return carapace.ActionMessage("Should not have more than one '//'").Style(style.Red)
						}
					})
				default:
					return carapace.ActionValues()
				}

			})

		}),
	)
}
