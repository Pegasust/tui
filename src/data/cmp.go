package data

import (
	"bytes"
	"fmt"
	"os"
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
	return carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		// start with <tab>; no typing
		case 0:
			return carapace.ActionValuesDescribed(
				compe.cells...,
			).Invoke(c).Prefix("//").Suffix("/").ToA().Style(
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
	})
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
	return carapace.ActionMessage("TODO: implement flakehub")
}

func githubFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMessage("TODO: implement github")
}

func gitlabFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMessage("TODO: implement gitlab")
}

func sourcehutFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMessage("TODO: implement sourcehut")
}

func flakeRegistryCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMessage("TODO: implement flake reg")
}

func gitQueryCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	// Requires ctx.Value to point at `?|query=value`
	return carapace.ActionMultiParts("&", func(c carapace.Context) carapace.Action {
		return carapace.ActionMultiParts("=", func(c carapace.Context) carapace.Action {
			git_queries := []string{
				"dir", "subdirectory",
				"rev", "branch name, tag name, or tagged commit-ish",
				"ref", "commit-ish reference",
			}
			query_vars := carapace.ActionValuesDescribed(
				git_queries...,
			).Suffix("=").NoSpace()
			switch len(c.Parts) {
			case 0:
				return query_vars
			case 1:
				return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					next_query_var := query_vars.Invoke(c).Prefix(c.Value + "&").NoSpace()
					value_cmp := func() carapace.Action {
						switch c.Parts[0] {
						case "dir":
							return carapace.ActionMessage("git subdirectory")
						case "rev":
							return carapace.ActionMessage("branch name, tag name, or tagged commit-ish")
						case "ref":
							return carapace.ActionMessage("commit-ish ref")
						default:
							return carapace.ActionValues()
						}
					}
					if len(c.Value) == 0 {
						return value_cmp()
					}
					return carapace.Batch(
						next_query_var,
						value_cmp(),
					).Invoke(c).Merge().Action
				})
			default:
				return carapace.ActionValues()
			}
		}).Invoke(c).ToA()
	}).Invoke(ctx).ToA()
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func localGitFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMultiParts("?", func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		case 0:
			prefix := "?"
			if isDir(c.Value) && len(c.Value) > 0 {
				prefix = c.Value + "?"
			}
			return carapace.Batch(
				carapace.ActionDirectories(),
				gitQueryCompe(compe, c).Prefix(prefix),
			).ToA()
		case 1:
			return gitQueryCompe(compe, c)
		default:
			return carapace.ActionMessage(
				"Not conforming to URL+query. Should have <= 1 '?', got %v",
				len(c.Parts),
			).Style(style.Red)
		}
	}).Invoke(ctx).ToA()
}

var localTarballAction = carapace.ActionFiles(
	".zip", ".tar", ".tgz", ".tar.gz", ".tar.xz", ".tar.bz2", ".tar.zst",
)

func localTarballCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	// As long as you end with .tar.gz, I'm happy
	return localTarballAction
}

func remoteGitFlakeCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMultiParts("?", func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		case 0:
			// TODO: Filter out only git + flake url
			return carapace.ActionMessage("[(:username)@|(:username):(:password)@|](:domain.name)((\"/\":path))*")
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

func remoteTarballCompe(compe *CompeIndex, ctx carapace.Context) carapace.Action {
	return carapace.ActionMultiParts("?", func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		case 0:
			// TODO: Filter out only git + flake url
			return carapace.ActionMultiParts("/", func(c carapace.Context) carapace.Action {
				switch len(c.Parts) {
				case 0:
					return carapace.ActionMessage(
						"[(:username)@|(:username):(:password)@|](:domain.name)",
					)
				default:
					c.Value = strings.Join(c.Parts, "/") + c.Value
					c.Parts = make([]string, 0, 0)
					// NOTE: test this, I'm not sure how internal mutation of c is like
					return localTarballAction.Invoke(c).Chdir(os.DevNull)
				}

			})
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

func compeSourceFromAction(a carapace.Action) compeSource {
	return func(_ *CompeIndex, c carapace.Context) carapace.Action {
		return a.Invoke(c).ToA()
	}
}

type compeProduct struct {
	compeSource
	desc  string
	style string
}

func (prod *compeProduct) overrideCompeSource(cmp compeSource) compeProduct {
	return compeProduct{
		compeSource: cmp,
		desc:        prod.desc,
		style:       prod.style,
	}
}
func xferDispatch(compe *CompeIndex) carapace.Action {
	return carapace.ActionMultiParts("://", func(c carapace.Context) carapace.Action {
		localTarball := compeProduct{
			compeSource: localTarballCompe,
			desc:        "Local tarball",
			style:       style.Yellow,
		}

		remoteTarball := compeProduct{
			compeSource: remoteTarballCompe,
			desc:        "Remote tarball",
			style:       style.Green,
		}

		localGitFlake := compeProduct{
			compeSource: localGitFlakeCompe,
			desc:        "Local git tree",
			style:       style.Default,
		}

		remoteGitFlake := compeProduct{
			compeSource: remoteGitFlakeCompe,
			desc:        "Remote git tree",
			style:       style.Green,
		}

		remoteMercurialFlake := compeProduct{
			compeSource: remoteGitFlakeCompe,
			desc:        "Remote mg tree",
			style:       style.Cyan,
		}

		localMercurialFlake := compeProduct{
			compeSource: localGitFlakeCompe,
			desc:        "Local mg tree",
			style:       style.Blue,
		}

		dispatch := map[string]compeProduct{
			"file":      localTarball,
			"file+file": localTarball,
			// NOTE: this actually hints to `nix` that we acknowledge
			// tarballs with unconventional extension
			"tarball+file": localTarball.overrideCompeSource(
				compeSourceFromAction(carapace.ActionFiles()),
			),

			"file+http":     remoteTarball,
			"file+https":    remoteTarball,
			"tarball+http":  remoteTarball,
			"tarball+https": remoteTarball,
			"https":         remoteTarball,

			"path":     localGitFlake,
			"git+file": localGitFlake,
			"git+git":  localGitFlake,

			"git+http":  remoteGitFlake,
			"git+https": remoteGitFlake,
			"git+ssh":   remoteGitFlake,

			"hg+http":  remoteMercurialFlake,
			"hg+https": remoteMercurialFlake,
			"hg+ssh":   remoteMercurialFlake,

			"hg+file": localMercurialFlake,
		}
		switch len(c.Parts) {
		case 0:
			parts := make([]string, 0, len(dispatch)*3)
			for k, v := range dispatch {
				parts = append(parts, k, v.desc, v.style)
			}
			return carapace.ActionStyledValuesDescribed(parts...).Suffix("://").NoSpace()
		case 1:
			call, contains := dispatch[c.Parts[0]]
			if !contains {
				return carapace.ActionValues()
			}
			return call.compeSource(compe, c)
		default:
			return carapace.ActionMessage("Should not have more than one '://'").Style(style.Red)
		}
	})
}

func FlakeRefCompletion(index *CompeIndex, c carapace.Context) carapace.Action {
	return carapace.Batch(
		// //<std...>
		StdPathCompe(index, c),
		// .#__std//
		localGitFlakeCompe(index, c),
		// nixpkgs#__std//
		flakeRegistryCompe(index, c),
		// <xfer-protocol>://<...>#__std//
		xferDispatch(index),
		// <nix-special>:<...>#__std//
		carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValuesDescribed(
					FlakeHubProto, "Flakehub",
					GitHubProto, "GitHub",
					GitLabProto, "GitLab",
					SourceHutProto, "SourceHut",
					NixRegistryProto, "Nix registry",
				).Invoke(c).Suffix(":").NoSpace()
			case 1:
				dispatch := map[string](compeSource){
					// NOTE: API specs
					FlakeHubProto:    flakehubFlakeCompe,
					GitHubProto:      githubFlakeCompe,
					GitLabProto:      gitlabFlakeCompe,
					SourceHutProto:   sourcehutFlakeCompe,
					NixRegistryProto: flakeRegistryCompe,
				}
				call, contains := dispatch[c.Parts[0]]
				if !contains {
					return carapace.ActionValues()
				}
				return call(index, c)
			default:
				return carapace.ActionValues()
				// return carapace.ActionMessage("Ensure Nix support, docs: https://nixos.org/manual/nix/stable/command-ref/new-cli/nix3-flake.html#types")
			}
		}),
	).Invoke(c).Merge().Action
}

func PaisanoActionCompletion(cmd *carapace.Carapace, argv0 string) {
	// completes: '//cell/block/target:action'
	cmd.PositionalCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {

			// TODO: support remote caching, maybe in XDG cache
			local := flake.LocalPaisanoRegistry()
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
					return FlakeRefCompletion(&index, c)
				case 1:
					// paisano registry
					well_known := map[string]string{
						flake.BrandedRegistry: "Custom branded registry",
						"__std":               "github:divnix/std",
					}
					dedup_registries := make([]string, 0, len(well_known)*2)
					for name, desc := range well_known {
						dedup_registries = append(dedup_registries, name, desc)
					}
					registries := carapace.ActionValuesDescribed(
						dedup_registries...,
					).NoSpace().Suffix("//").Invoke(c).ToA()
					return carapace.Batch(
						registries,
						StdPathCompe(&index, c),
					).ToA()
				default:
					return carapace.ActionValues()
				}

			}).Invoke(c).ToA()

		}),
	)
}
