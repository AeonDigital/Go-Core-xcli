package xcli

import (
	"fmt"
	"strings"

	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xclifn"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xclistruc"
	"github.com/AeonDigital/Go-Core-xerrors/pkg/xerrors"
)

// Command represents a node in the CLI command tree.
// Each command forms an isolated scope and executes its own business logic.
type Command struct {
	// Name is the string that triggers this command in the terminal.
	Name string

	// ShortDescription is a brief one-line summary used in general help listings.
	ShortDescription string

	// LongDescription is a detailed explanation shown when help is requested
	// specifically for this command. If left empty, ShortDescription will be used.
	LongDescription string

	// Flags is the list of exclusive options accepted strictly by this command.
	Flags []xclistruc.Flag

	// Subcommands holds the next layer of commands, indexed by their execution Name.
	// Due to context isolation, children do not inherit flags from their parents.
	Subcommands map[string]*Command

	// Run is the execution hook containing the command's business logic.
	// It receives the Context containing all parsed, typed, and validated flags.
	Run func(ctx *xclistruc.FlagValues) error
}

// ValidateAndHydrateFlags loops through registered constraints performing types translation and bounds enforcement.
//
// It checks flag mandatory presence, injects fallback defaults allocation blocks, and leverages
// the package-level central registry directory to resolve runtime type translations dynamically.
//
// Arguments:
//   - rawFlags: The extracted command line layout map pairing flag strings to text values.
//
// Returns:
//   - *xclistruc.FlagValues: A populated type-safe context directory containing translated Go instances.
//   - error: Returns an xerrors.IErrorCLI capturing visual telemetry diagnostics if any validation path fails.
//
// Error & Panic Natures:
//   - Complex Errors: Fails immediately if a required token is missing from raw inputs.
//     Queries the global type registry throwing errors if a type metadata token is unmapped.
//     Performs zombie variable tracking blocks, failing if unregistered flag tokens are found.
func (c *Command) ValidateAndHydrateFlags(
	rawFlags map[string]string,
) (
	*xclistruc.FlagValues,
	error,
) {
	ctxValues := xclistruc.NewFlagValues()

	// Track which keys from terminal inputs were actually processed to detect unmapped items later
	processedRawKeys := make(map[string]bool)

	for _, spec := range c.Flags {
		// Enforce visual notation representation matching user input options
		flagLabel := "--" + spec.LongName

		// Locate the raw value token inside the parsed terminal mapping block
		var rawValue string
		var found bool

		if val, ok := rawFlags[spec.LongName]; ok {
			rawValue = val
			found = true
			processedRawKeys[spec.LongName] = true
		} else if spec.ShortName != "" {
			if val, ok := rawFlags[spec.ShortName]; ok {
				rawValue = val
				found = true
				processedRawKeys[spec.ShortName] = true
			}
		}

		// 1. Mandatory Presence enforcement check
		if spec.Required && !found {
			return nil, xerrors.NewErrorCLI().
				SetMessage("[ERR] %s : required", flagLabel)
		}

		// 2. Default values allocation fallback logic
		if !found {
			if spec.DefaultValue != nil {
				ctxValues.SetInternalValue(spec.LongName, spec.DefaultValue)
			}
			continue
		}

		// 3. Dynamic Type-Safe Parsing and Unified Multi-Layered Validation Engine Dispatch via Global Registry
		parserEngine, exists := GlobalTypeRegistry[spec.Type]
		if !exists {
			return nil, xerrors.NewErrorCLI().
				SetMessage("unsupported flag type: '%s'", spec.Type)
		}

		// Execute the complete type transformation and domain boundary enforcement pipeline
		typedValue, parseErr := parserEngine.ParseAndValidate(flagLabel, rawValue, spec)
		if parseErr != nil {
			return nil, parseErr
		}

		ctxValues.SetInternalValue(spec.LongName, typedValue)
	}

	// 4. Security Check: Block unregistered positional arguments or zombie variables sent by mistake
	for k := range rawFlags {
		if !processedRawKeys[k] {
			return nil, xerrors.NewErrorCLI().
				SetMessage("invalid argument: unrecognized flag identifier token provided: '%s'", k)
		}
	}

	return ctxValues, nil
}

// TriggerHelp intercepts the flow and renders the automatic command documentation.
//
// It evaluates operational descriptions metadata fields and directly serializes aligned visual
// tables mapping usage definitions, child subcommands, and flags guidelines to standard stdout tracks.
//
// Returns:
//   - error: Returns a structured error tracking instance if writing stdout triggers telemetry blocks.
func (c *Command) TriggerHelp() error {
	// Determine the best description to display based on availability
	description := c.LongDescription
	if description == "" {
		description = c.ShortDescription
	}

	// Render Command Usage and Description
	xclifn.PrintStdout("Usage:")
	xclifn.PrintStdout("  %s [subcommand] [--flags]\n", c.Name)

	if description != "" {
		xclifn.PrintStdout("\nDescription:")
		xclifn.PrintStdout("  %s\n", description)
	}

	// Render Available Subcommands if the command has children
	if len(c.Subcommands) > 0 {
		xclifn.PrintStdout("\nAvailable Subcommands:")
		for _, sub := range c.Subcommands {
			xclifn.PrintStdout("  %-15s %s", sub.Name, sub.ShortDescription)
		}
		xclifn.PrintStdout("")
	}

	// Render Flags Configuration if the command accepts any options
	if len(c.Flags) > 0 {
		xclifn.PrintStdout("\nFlags:")

		// Pass 1: Initialize storage slices and evaluate the maximum visual padding required
		totalFlags := len(c.Flags)
		syntaxCollection := make([]string, totalFlags)
		descCollection := make([]string, totalFlags)
		metaCollection := make([]string, totalFlags)

		maxFlagWidth := 0

		for i, flag := range c.Flags {
			// 1.1 Build and format the precise flag trigger syntax line
			syntaxStr := "  --" + flag.LongName
			if flag.ShortName != "" {
				syntaxStr += ", -" + flag.ShortName
			}
			syntaxCollection[i] = syntaxStr

			// Track dynamically the maximum horizontal length occupied by the flag triggers
			if len(syntaxStr) > maxFlagWidth {
				maxFlagWidth = len(syntaxStr)
			}

			// 1.2 Resolve description fallback strings
			flagDesc := flag.LongDescription
			if flagDesc == "" {
				flagDesc = flag.ShortDescription
			}
			descCollection[i] = flagDesc

			// 1.3 Compile the detailed micro-metadata parameter block
			metaStr := fmt.Sprintf("[Type: %s]", flag.Type)
			if flag.Required {
				metaStr += " [Required]"
			} else if flag.DefaultValue != nil {
				metaStr += fmt.Sprintf(" [Default: %v]", flag.DefaultValue)
			}
			metaCollection[i] = metaStr
		}

		// Calculate total horizontal indent padding boundary (+4 spaces minimum separation)
		descriptionIndentSize := maxFlagWidth + 4
		leftPaddingSpaces := strings.Repeat(" ", descriptionIndentSize)

		// Pass 2: Synchronized iteration over all matrix layers to construct the final layout
		for i := 0; i < totalFlags; i++ {
			syntax := syntaxCollection[i]
			desc := descCollection[i]
			meta := metaCollection[i]

			// Create dynamic right padding to align the first line description start
			rightPaddingSize := descriptionIndentSize - len(syntax)
			rightPadding := strings.Repeat(" ", rightPaddingSize)

			// Execute 120-character line constraint wrapping for multi-line descriptions
			maxDescLineLength := 120 - descriptionIndentSize
			descLines := []string{}

			// Simple character-counting split routine to safeguard terminal viewport limits
			words := strings.Fields(desc)
			if len(words) == 0 {
				descLines = append(descLines, "")
			} else {
				currentLine := ""
				for _, word := range words {
					if len(currentLine)+len(word)+1 > maxDescLineLength {
						descLines = append(descLines, currentLine)
						currentLine = word
					} else {
						if currentLine == "" {
							currentLine = word
						} else {
							currentLine += " " + word
						}
					}
				}
				if currentLine != "" {
					descLines = append(descLines, currentLine)
				}
			}

			// Render the organized structure to the output channel
			var buffer strings.Builder

			// Line 1: Primary trigger entry + exact padding separator + first slice of description text
			fmt.Fprintf(&buffer, "%s%s%s\n", syntax, rightPadding, descLines[0])

			// Lines 2+: If description overflows the 120char ceiling, indent trailing entries
			for idx := 1; idx < len(descLines); idx++ {
				fmt.Fprintf(&buffer, "%s%s\n", leftPaddingSpaces, descLines[idx])
			}

			// Final Line: Metadata specifications perfectly nested and aligned on the terminal layer
			fmt.Fprintf(&buffer, "%s%s\n", leftPaddingSpaces, meta)

			// Commit the complete autonomous block down to the system stdout stream
			xclifn.PrintStdout(buffer.String())
		}
	}
	return nil
}
