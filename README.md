Go Core xcli
================================

![Go Test Coverage](https://raw.githubusercontent.com/AeonDigital/Go-Core-xcli/badges/.badges/main/coverage.svg)

> [Aeon Digital](http://aeondigital.com.br)  
> rianna@aeondigital.com.br

&nbsp;

> High-performance, strongly-typed Command Line Interface (CLI) routing engine for Go

&nbsp;

This package delivers an enterprise-grade command-tree routing framework designed under an absolute isolation-of-scope architecture. 
It provides automated multi-layered data parsing, strict boundary validation tracks, and native Unix-like help interfaces with zero third-party dependencies.


&nbsp;
&nbsp;


________________________________________________________________________________

## INSTALLATION

Use `go get` to install the package repository directly into your environment:

```shell
go get github.com/AeonDigital/Go-Core-xcli@latest
```


&nbsp;
&nbsp;


________________________________________________________________________________

## PURPOSE & ARCHITECTURE

The core mission of `xcli` is to completely separate raw command-line text inputs from your application's downstream domain logic. 
To guarantee that your codebase never touches corrupted data, execution runs through an isolated two-phase cycle.


&nbsp;


### Phase 1: Command Tree Routing

The router traverses registered commands strictly via hierarchy trees (`maincmd subcmd subsubcmd`). 
Context isolation is absolute: child subcommands never inherit flags from their parent nodes. 
Every command node acts as an independent execution sandbox.


&nbsp;


### Phase 2: Isolation, Parsing, and Validation

Once the target command node is resolved, the router strips away command names and passes raw flag tokens to a functional pipeline. 
This layer automatically cleans strings, handles types transformation, enforces constraint limits, and performs active filesystem privilege auditing before triggering the execution hook.


&nbsp;
&nbsp;


________________________________________________________________________________

## QUICK START EXAMPLE

The following complete example demonstrates how to initialize the router, register a command branch with structural constraints, and boot the execution phase using terminal inputs.


&nbsp;


```go
package main

import (
	"os"

	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xcliconstt"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xclifn"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xclistruc"
	"github.com/AeonDigital/Go-Core-xerrors/pkg/xerrors"
)

func main() {
	// 1. Declare target operational command nodes configurations
	buildCmd := &xcli.Command{
		Name:             "build",
		ShortDescription: "Compiles the active project artifacts.",
		Flags: []xclistruc.Flag{
			{
				LongName:         "output",
				ShortName:        "out",
				ShortDescription: "Target physical location path for compilation artifacts.",
				Required:         true,
				Type:             xcliconstt.TypeFilepath,
				MustExist:        false,
				Access:           xcliconstt.ModeWrite,
			},
			{
				LongName:         "tags",
				ShortDescription: "Strict JSON metadata matrix array collection.",
				Required:         false,
				Type:             xcliconstt.TypeStringArray,
				DefaultValue:     []string{"production"},
			},
		},
		Run: func(ctx *xclistruc.FlagValues) error {

			outputPath := ctx.GetString("output")
			buildTags := ctx.GetStringSlice("tags")

			xclifn.PrintStdout("Compiling target deployment payload to: %s", outputPath)
			xclifn.PrintStdout("Active build metadata tags: %v", buildTags)
			return nil
		},
	}

	rootCmd := &xcli.Command{
		Name:             "app",
		ShortDescription: "Main application service gateway interface.",
		Subcommands: map[string]*xcli.Command{
			"build": buildCmd,
		},
	}

	// 2. Initialize the orchestration motor router component
	cliRouter := xcli.NewRouter(rootCmd)

	// 3. Forward terminal prompt inputs to trigger the execution matrix layers
	err := cliRouter.Run(os.Args[1:])
	if err != nil {
		cliErr, ok := err.(xerrors.IErrorCLI)
		if ok {
			xclifn.PrintStderr(cliErr.GetUserMessage())
		} else {
			xclifn.PrintError(err)
		}
		os.Exit(1)
	}
}
```


&nbsp;
&nbsp;


________________________________________________________________________________

## DATA TYPE SUPPORT MATRIX

The framework embeds a centralized, imuttable type engine mapping primitives, structured layouts, and system resources. 
Every data format listed below is natively supported as an individual flag value or as a comma-separated array instance (`[]type`).


&nbsp;


### Primitive Types

* `string`: Text variables automatically stripped of surrounding blank spaces using Unicode runes.
* `bool`: Boolean activations supporting explicit `true`, `false`, `1`, or `0` keywords.
* `int`: Concrete native integer numeric primitives.
* `float`: High-precision floating-point coordinates.


&nbsp;


### Structured Data & Temporal Layouts

* `json`: Validated raw JSON text sequences scrutinized under internal format checkers.
* `duration`: Chronological time ranges parsed through native duration behaviors.
* `date`: Calendar dates managed under standard layout strings.
* `time`: Daily timeline configurations.
* `datetime`: Full timestamps combining chronological points.
* `email`: Electronic mail format validated against RFC 5322 specifications.


&nbsp;


### System & Network Resources

* `path`: Structural syntax verification checking directory structural path arrangements.
* `filename`: Isolated file naming protecting descriptors from illegal slash sequences.
* `filepath`: Real physical validation verifying file path presence on the machine disk.
* `dirname`: Semantic directory structure validations.
* `dirpath`: Real physical folder verification requiring directory confirmation on disk.
* `url`: Uniform resource locator paths.
* `fullurl`: Complete internet resource addresses requiring scheme and host blocks.
* `relativeurl`: Clean system target sub-paths.


&nbsp;
&nbsp;


________________________________________________________________________________

## CORE DECLARATION STRUCTURES

Developers declare commands and flags through declarative object schemas. The internal routing engine enforces these rules at runtime.

&nbsp;

### Flag Specification Schema

```go
type Flag struct {
  LongName         string     // e.g., "output" invoked via --output
  ShortName        string     // Shorthand up to 3 chars (no bundling) e.g., -out
  ShortDescription string     // Brief inline summary for help lists
  LongDescription  string     // Detailed operational usage explanation
  Required         bool       // If true, command execution fails if absent
  DefaultValue     any        // Fallback value if omitted (ignored if Required)
  Type             FlagType   // Targeted conversion descriptor mapping token
  Min              any        // Inclusive lower bound for quantifiable types
  Max              any        // Inclusive upper bound for quantifiable types
  MinLength        *int       // Unicode string character lower bound count
  MaxLength        *int       // Unicode string character upper bound count
  Regex            string     // Regular expression string pattern verification
  MinItems         *int       // Minimum item capacity count for arrays
  MaxItems         *int       // Maximum item capacity count for arrays
  MustExist        bool       // Enforces that file paths must reside on disk
  Access           AccessMode // Permission auditing level (read, write, readwrite)
}
```


&nbsp;


### Command Tree Specification Schema

```go
type Command struct {
  Name             string
  ShortDescription string
  LongDescription  string
  Flags            []xclistruc.Flag
  Subcommands      map[string]*Command
  Run              func(ctx *xclistruc.FlagValues) error
}
```


&nbsp;
&nbsp;


________________________________________________________________________________

## SYSTEM PIPELINE MECHANICS

The router coordinates data parsing and constraints application using a three-tier processing workflow.


&nbsp;


### 1. Token Cleansing & Structural Parsing

Raw terminal input text sequences are gathered and isolated. 
If a target flag matches an array descriptor (`[]type`), the entry point splits the text block using comma tokenization loops before iterating elements.


&nbsp;


### 2. Mathematical & Constraint Boundary Evaluation

Primitives pass through strict mathematical limits. 
Strings are inspected via Unicode counts to enforce size restrictions accurately, followed by targeted regular expression matching.


&nbsp;


### 3. Physical Disk Resource & Capability Verification

Filesystem targets (`filepath`, `dirpath`) are tested directly against physical OS components. 
The system audits permissions based on the `Access` configuration and checks existing items based on the `MustExist` toggle.


&nbsp;
&nbsp;


________________________________________________________________________________

## SYSTEM DEPENDENCIES

`xcli` maintains a lightweight production footprint, keeping its internal design decoupled from massive external vendor dependencies.

* **Language Runtime:** Go 1.18+ (utilizing compile-time generic descriptors).
* **External Packages:** Relying entirely on Go standard library modules (`os`, `strings`, `time`, `net/url`, `net/mail`, `path/filepath`).


&nbsp;
&nbsp;


________________________________________________________________________________

## LICENSE

This project is offered under the [MIT license](LICENSE.md).
