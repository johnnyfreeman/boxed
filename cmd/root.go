package cmd

import (
	"fmt"
	"io"
	"os"

	"boxed/internal/box"
	boxio "boxed/internal/io"
	"boxed/internal/parser"
	"boxed/internal/render"

	"github.com/spf13/cobra"
)

// Executor encapsulates the dependencies needed to execute a box render command.
// This structure enables dependency injection, making the CLI logic testable by
// allowing test code to inject mock renderers and custom writers instead of
// always using os.Stdout.
type Executor struct {
	renderer render.Renderer
	writer   io.Writer
}

// NewExecutor creates an executor with the given dependencies.
// The renderer handles box-to-string conversion (injected to allow mock renderers
// in tests), and the writer receives the final output (injected to capture output
// in tests rather than always printing to os.Stdout).
func NewExecutor(renderer render.Renderer, writer io.Writer) *Executor {
	return &Executor{
		renderer: renderer,
		writer:   writer,
	}
}

// Execute performs the complete flow: parse → validate → render → output.
// This method coordinates the entire pipeline but remains simple because each
// step is handled by dedicated, well-tested modules. The method itself contains
// no business logic, just composition of validated components.
func (e *Executor) Execute(boxType string, opts parser.Options, useStdin bool, useJSON bool, jsonFile string, exitOnError bool, exitOnWarning bool) error {
	// JSON input takes precedence over other options
	if useJSON || jsonFile != "" {
		var reader *boxio.JSONReader
		if jsonFile != "" {
			file, err := os.Open(jsonFile)
			if err != nil {
				return fmt.Errorf("failed to open JSON file: %w", err)
			}
			defer file.Close()
			reader = boxio.NewJSONReader(file)
		} else {
			reader = boxio.NewJSONReader(os.Stdin)
		}

		jsonOpts, err := reader.ReadBox()
		if err != nil {
			return fmt.Errorf("failed to read JSON: %w", err)
		}

		// Merge JSON options with CLI options (CLI takes precedence for overrides)
		if opts.Title == "" {
			opts.Title = jsonOpts.Title
		}
		if opts.Subtitle == "" {
			opts.Subtitle = jsonOpts.Subtitle
		}
		if opts.Footer == "" {
			opts.Footer = jsonOpts.Footer
		}
		if opts.Width == 0 {
			opts.Width = jsonOpts.Width
		}
		if opts.BorderStyle == "" {
			opts.BorderStyle = jsonOpts.BorderStyle
		}
		opts.KVFlags = append(opts.KVFlags, jsonOpts.KVFlags...)
	} else if useStdin {
		reader := boxio.NewStdinKVReader(os.Stdin)
		stdinKVs, err := reader.ReadKVPairs()
		if err != nil {
			return fmt.Errorf("failed to read KV pairs from stdin: %w", err)
		}

		for _, kv := range stdinKVs {
			opts.KVFlags = append(opts.KVFlags, kv.String())
		}
	}

	b, err := parser.ParseBox(boxType, opts)
	if err != nil {
		return err
	}

	output := e.renderer.RenderBox(b)

	_, err = fmt.Fprintln(e.writer, output)
	if err != nil {
		return err
	}

	// Exit with non-zero code based on box type if flags are set
	if exitOnError && b.Type == box.Error {
		os.Exit(1)
	}
	if exitOnWarning && b.Type == box.Warning {
		os.Exit(2)
	}

	return nil
}

// NewRootCmd creates the root cobra command with all subcommands configured.
// Each box type (success, error, info, warning) gets its own subcommand sharing
// the same flag definitions. This design makes the CLI intuitive: users type
// "boxed success ..." rather than "boxed --type success ...", which reads more
// naturally and follows common CLI patterns.
func NewRootCmd(executor *Executor) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "boxed",
		Short: "Render beautiful bordered boxes for terminal output",
		Long: `boxed renders bordered status boxes in the terminal with support for
titles, subtitles, key-value pairs, and footers. Perfect for deployment scripts,
CI/CD pipelines, and any command-line tool that needs clear visual status output.`,
		SilenceUsage: true,
	}

	makeBoxCmd := func(boxType box.BoxType) *cobra.Command {
		var title, subtitle, footer, borderStyle string
		var kvFlags []string
		var width int
		var useStdin, useJSON, exitOnError, exitOnWarning bool
		var jsonFile string

		cmd := &cobra.Command{
			Use:   string(boxType),
			Short: fmt.Sprintf("Render a %s box", boxType),
			Long:  fmt.Sprintf("Render a %s box with %s border color", boxType, getColorName(boxType)),
			Example: fmt.Sprintf(`  boxed %s --title "Deploy Complete"
  boxed %s --title "Build v2.1.0" --kv "Duration=2m 34s" --kv "Commit=abc1234"
  boxed %s --title "Status" --subtitle "Production" --footer "Updated 2025-10-19"
  echo -e "env=prod\nregion=us-east-1" | boxed %s --title "Config" --stdin-kv
  echo '{"title":"Status","kv":{"CPU":"45%%"}}' | boxed %s --json
  boxed %s --json-file status.json`,
				boxType, boxType, boxType, boxType, boxType, boxType),
			RunE: func(cmd *cobra.Command, args []string) error {
				opts := parser.Options{
					Title:       title,
					Subtitle:    subtitle,
					KVFlags:     kvFlags,
					Footer:      footer,
					Width:       width,
					BorderStyle: borderStyle,
				}

				return executor.Execute(string(boxType), opts, useStdin, useJSON, jsonFile, exitOnError, exitOnWarning)
			},
		}

		cmd.Flags().StringVarP(&title, "title", "t", "", "Box title (bold, centered)")
		cmd.Flags().StringVarP(&subtitle, "subtitle", "s", "", "Box subtitle (italic, centered)")
		cmd.Flags().StringArrayVarP(&kvFlags, "kv", "k", nil, "Key-value pairs (repeatable, format: key=value)")
		cmd.Flags().StringVarP(&footer, "footer", "f", "", "Box footer (faint, centered)")
		cmd.Flags().IntVarP(&width, "width", "w", 0, "Box width (0 for auto-size)")
		cmd.Flags().StringVarP(&borderStyle, "border-style", "b", "rounded", "Border style (normal, rounded, thick, double)")
		cmd.Flags().BoolVar(&useStdin, "stdin-kv", false, "Read additional KV pairs from stdin (one per line)")
		cmd.Flags().BoolVar(&useJSON, "json", false, "Read box definition from JSON stdin")
		cmd.Flags().StringVar(&jsonFile, "json-file", "", "Read box definition from JSON file")
		cmd.Flags().BoolVar(&exitOnError, "exit-on-error", false, "Exit with code 1 when rendering an error box")
		cmd.Flags().BoolVar(&exitOnWarning, "exit-on-warning", false, "Exit with code 2 when rendering a warning box")
		cmd.MarkFlagsMutuallyExclusive("stdin-kv", "json", "json-file")

		return cmd
	}

	rootCmd.AddCommand(
		makeBoxCmd(box.Success),
		makeBoxCmd(box.Error),
		makeBoxCmd(box.Info),
		makeBoxCmd(box.Warning),
	)

	return rootCmd
}

func getColorName(t box.BoxType) string {
	switch t {
	case box.Success:
		return "green"
	case box.Error:
		return "red"
	case box.Info:
		return "blue"
	case box.Warning:
		return "yellow"
	default:
		return "default"
	}
}
