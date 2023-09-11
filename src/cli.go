package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/oriser/regroup"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"

	"github.com/paisano-nix/paisano/data"
	"github.com/paisano-nix/paisano/flake"
)

type Spec struct {
	FlakeRef string `regroup:"flake_ref,optional"`
	Registry string `regroup:"registry,optional"`
	Cell     string `regroup:"cell,required"`
	Block    string `regroup:"block,required"`
	Target   string `regroup:"target,required"`
	Action   string `regroup:"action,required"`
}

var Use = fmt.Sprintf("%[1]s [flakeref]//[cell]/[block]/[target]:[action] [args...]", argv0)

var paisanoRegistryRegex = `(?:(?P<flake_ref>[^#]*))?#?(?:(?P<registry>[^/]+))?`
var paisanoRegistryRe = regroup.MustCompile("^" + paisanoRegistryRegex)

var specRe = regroup.MustCompile("^" + paisanoRegistryRegex + `//(?P<cell>[^/]+)/(?P<block>[^/]+)/(?P<target>.+):(?P<action>[^:]+)`)

func parseSpec(specRef string) (*Spec, error) {
	s := &Spec{}
	if err := specRe.MatchToTarget(specRef, s); err != nil {
		return nil, fmt.Errorf("invalid argument format: %s, should follow %s: %w", specRef, Use, err)
	}
	if s.FlakeRef == "" {
		s.FlakeRef = "."
	}

	if strings.HasPrefix(s.FlakeRef, data.FlakeHubProto+":") {
		// NOTE: handle the case of fh:org/repo/semver since not all Nix versions support this
		// https://flakehub.com/docs
		s.FlakeRef = s.FlakeRef[len(data.FlakeHubProto)+1:]
		tup := strings.SplitN(s.FlakeRef, "/", 3)
		org, repo, semver := tup[0], tup[1], (func() string {
			if len(tup) == 2 {
				return "*"
			} else {
				return tup[2]
			}
		})()
		fh, err := url.Parse("https://api.flakehub.com")
		if err != nil {
			return nil, err
		}

		// "f" is from API call reverse eng
		s.FlakeRef = fh.JoinPath("f", org, repo, strings.Join([]string{
			semver, "tar", "gz",
		}, ".")).String()
	}

	if s.Registry == "" {
		s.Registry = flake.BrandedRegistry
	}
	return s, nil
}

func ParseRunActionCmd(specRef string) (*flake.RunActionCmd, error) {
	s, err := parseSpec(specRef)
	if err != nil {
		return nil, err
	}
	return &flake.RunActionCmd{
		PaisanoRegistry: flake.PaisanoRegistry{
			FlakeRef:        s.FlakeRef,
			Registry:        s.Registry,
			PaisanoRegistry: fmt.Sprintf("%s#%s", s.FlakeRef, s.Registry),
		},
		System: forSystem,
		Cell:   s.Cell,
		Block:  s.Block,
		Target: s.Target,
		Action: s.Action,
	}, nil
}

func RunEActionCmd(f func(cmd *cobra.Command, args []string, a *flake.RunActionCmd) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		act, err := ParseRunActionCmd(args[0])
		if err != nil {
			return fmt.Errorf("error parsing action cmd: %w", err)
		}
		return f(cmd, args, act)
	}
}

var forSystem string

var rootCmd = &cobra.Command{
	Use:                   Use,
	DisableFlagsInUseLine: true,
	Version:               fmt.Sprintf("%s (%s)", buildVersion, buildCommit),
	Short:                 fmt.Sprintf("%[1]s is the CLI / TUI companion for %[2]s", argv0, project),
	Long: fmt.Sprintf(`%[1]s is the CLI / TUI companion for %[2]s.

- Invoke without any arguments to start the TUI.
- Invoke with a target spec and action to run a known target's action directly.

Enable autocompletion via '%[1]s _carapace <shell>'.
For more instructions, see: https://rsteube.github.io/carapace/carapace/gen/hiddenSubcommand.html
`, argv0, project),
	Args: RunEActionCmd(func(_ *cobra.Command, _ []string, _ *flake.RunActionCmd) error {
		return nil
	}),
	RunE: RunEActionCmd(func(_ *cobra.Command, args []string, command *flake.RunActionCmd) error {
		if err := command.Exec(args[1:]); err != nil {
			return err
		}
		return nil

	}),
}

func ParsePaisanoRegistry(reg string) (*flake.PaisanoRegistry, error) {
	rv := flake.PaisanoRegistry{}
	if err := paisanoRegistryRe.MatchToTarget(reg, &rv); err != nil {
		return nil, fmt.Errorf("invalid argument format: %s, should follow %s: %w", reg, Use, err)
	}
	if rv.FlakeRef == "" {
		rv.FlakeRef = "."
	}

	if strings.HasPrefix(rv.FlakeRef, data.FlakeHubProto+":") {
		// NOTE: handle the case of fh:org/repo/semver since not all Nix versions support this
		// https://flakehub.com/docs
		rv.FlakeRef = rv.FlakeRef[len(data.FlakeHubProto)+1:]
		tup := strings.SplitN(rv.FlakeRef, "/", 3)
		org, repo, semver := tup[0], tup[1], (func() string {
			if len(tup) == 2 {
				return "*"
			} else {
				return tup[2]
			}
		})()
		fh, err := url.Parse("https://api.flakehub.com")
		if err != nil {
			return nil, err
		}

		// "f" is from API call reverse eng
		rv.FlakeRef = fh.JoinPath("f", org, repo, strings.Join([]string{
			semver, "tar", "gz",
		}, ".")).String()
	}

	if rv.Registry == "" {
		rv.Registry = flake.BrandedRegistry
	}
	rv.PaisanoRegistry = fmt.Sprintf("%s#%s", rv.FlakeRef, rv.Registry)
	return &rv, nil
}

func singleFlakeRefArg(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("%s expects single flake ref, got %d", cmd.CommandPath(), len(args))
	}
	if len(args) == 1 {
		_, err := ParsePaisanoRegistry(args[0])
		return err
	}
	return nil
}

var reCacheCmd = &cobra.Command{
	Use:   "re-cache",
	Short: "Refresh the CLI cache.",
	Long: `Refresh the CLI cache.
Use this command to cold-start or refresh the CLI cache.
The TUI does this automatically, but the command completion needs manual initialization of the CLI cache.`,
	Args: singleFlakeRefArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: support remote caching, maybe in XDG cache
		local := flake.LocalPaisanoRegistry()
		c, key, loadCmd, buf, err := local.LoadFlakeCmd()
		if err != nil {
			return fmt.Errorf("while loading flake (cmd '%v'): %w", loadCmd, err)
		}
		loadCmd.Run()
		c.PutBytes(*key, buf.Bytes())
		return nil
	},
}
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate the repository.",
	Long: fmt.Sprintf(`Validates that the repository conforms to %[1]s.
Returns a non-zero exit code and an error message if the repository is not a valid %[1]s repository.
The TUI does this automatically.`, project),
	Args: singleFlakeRefArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: support remote caching, maybe in XDG cache
		local_reg := flake.LocalPaisanoRegistry()
		reg_ref := &local_reg
		fmt.Printf("args: %v\n", args)
		if len(args) > 0 {
			_reg, err := ParsePaisanoRegistry(args[0])
			if err != nil {
				return fmt.Errorf("invalid registry %s: %v", args[0], err)
			}
			reg_ref = _reg
		}
		_, _, loadCmd, _, err := reg_ref.LoadFlakeCmd()
		loadCmd.Args = append(loadCmd.Args, "--trace-verbose")
		if err != nil {
			return fmt.Errorf("while loading flake (cmd '%v'): %w", loadCmd, err)
		}
		loadCmd.Stderr = os.Stderr
		if err := loadCmd.Run(); err != nil {
			os.Exit(1)
		}
		// NOTE: don't need to handle err since it would have failed earlier
		// under this route
		prefetchFlakeRef, _ := reg_ref.PrefetchFlakeRef()
		fmt.Printf("note: %s -> %s\n", reg_ref.FlakeRef, prefetchFlakeRef)
		fmt.Printf("Valid Paisano repository ✓\n")

		return nil
	},
}
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available targets.",
	Long: `List available targets.
Shows a list of all available targets. Can be used as an alternative to the TUI.
Also loads the CLI cache, if no cache is found. Reads the cache, otherwise.`,
	Args: singleFlakeRefArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: support remote caching, maybe in XDG cache
		local := flake.LocalPaisanoRegistry()
		cache, key, loadCmd, buf, err := local.LoadFlakeCmd()
		if err != nil {
			return fmt.Errorf("while loading flake (cmd '%v'): %w", loadCmd, err)
		}
		cached, _, err := cache.GetBytes(*key)
		var root *data.Root
		if err == nil {
			root, err = data.LoadJson(bytes.NewReader(cached))
			if err != nil {
				return fmt.Errorf("while loading cached json: %w", err)
			}
		} else {
			loadCmd.Run()
			bufA := &bytes.Buffer{}
			r := io.TeeReader(buf, bufA)
			root, err = data.LoadJson(r)
			if err != nil {
				return fmt.Errorf("while loading json (cmd: '%v'): %w", loadCmd, err)
			}
			cache.PutBytes(*key, bufA.Bytes())
		}
		w := tabwriter.NewWriter(os.Stdout, 5, 2, 4, ' ', 0)
		for _, c := range root.Cells {
			for _, o := range c.Blocks {
				for _, t := range o.Targets {
					for _, a := range t.Actions {
						fmt.Fprintf(w, "//%s/%s/%s:%s\t--\t%s:  %s\n", c.Name, o.Name, t.Name, a.Name, t.Description(), a.Description())
					}
				}
			}
		}
		w.Flush()
		return nil
	},
}

func ExecuteCli() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&forSystem, "for", "", "system, for which the target will be built (e.g. 'x86_64-linux')")
	rootCmd.AddCommand(reCacheCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(checkCmd)
	carapace.Gen(rootCmd).Standalone()
	// completes: '//cell/block/target:action'
	data.PaisanoActionCompletion(carapace.Gen(rootCmd), argv0)
}
