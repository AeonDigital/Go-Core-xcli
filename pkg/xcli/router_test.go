package xcli_test

import (
	"errors"
	"testing"

	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xcliconstt"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xcliintfc"
	"github.com/AeonDigital/Go-Core-xcli/pkg/xcli/xclistruc"
	"github.com/AeonDigital/Go-Core-xerrors/pkg/xerrors"
)

// mockParser atende a interface necessária para o GlobalTypeRegistry
type mockParser struct {
	shouldFail bool
}

// Garante que o mock implementa estritamente a interface necessária
var _ xcliintfc.ValueParser = (*mockParser)(nil)

func (m *mockParser) ParseAndValidate(flagLabel, rawValue string, spec xclistruc.Flag) (any, error) {
	if m.shouldFail {
		return nil, xerrors.NewErrorCLI().SetMessage("mock parsing failure")
	}
	return rawValue, nil
}

func TestRouter_Run(t *testing.T) {
	// Inicializa o GlobalTypeRegistry externamente para os testes de flag
	if xcli.GlobalTypeRegistry == nil {
		xcli.GlobalTypeRegistry = make(map[xcliconstt.FlagType]xcliintfc.ValueParser)
	}

	// Registra parsers mockados para os testes de sucesso e falha
	xcli.GlobalTypeRegistry["string"] = &mockParser{shouldFail: false}
	xcli.GlobalTypeRegistry["bad-type"] = &mockParser{shouldFail: false} // para forçar erro de tipo não registrado
	xcli.GlobalTypeRegistry["fail-type"] = &mockParser{shouldFail: true}

	// Remove "bad-type" para garantir que o teste de tipo não suportado falhe propositalmente
	delete(xcli.GlobalTypeRegistry, "bad-type")

	tests := []struct {
		name          string
		router        *xcli.Router
		rawArgs       []string
		wantErr       bool
		expectedError string
	}{
		{
			name:          "Erro: Root command nulo",
			router:        xcli.NewRouter(nil),
			rawArgs:       []string{},
			wantErr:       true,
			expectedError: "root command is not registered",
		},
		{
			name: "Sucesso: Ajuda disparada no root por string 'help'",
			router: xcli.NewRouter(&xcli.Command{
				Name:             "root",
				ShortDescription: "Root cmd",
			}),
			rawArgs: []string{"help"},
			wantErr: false,
		},
		{
			name: "Sucesso: Ajuda disparada no root por flag '--help'",
			router: xcli.NewRouter(&xcli.Command{
				Name:             "root",
				ShortDescription: "Root cmd",
			}),
			rawArgs: []string{"--help"},
			wantErr: false,
		},
		{
			name: "Sucesso: Ajuda disparada no root por flag '-h'",
			router: xcli.NewRouter(&xcli.Command{
				Name:             "root",
				ShortDescription: "Root cmd",
			}),
			rawArgs: []string{"-h"},
			wantErr: false,
		},
		{
			name: "Erro: Comando desconhecido no escopo",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
			}),
			rawArgs:       []string{"unknown-cmd"},
			wantErr:       true,
			expectedError: "unknown command: 'unknown-cmd' for scope 'root'",
		},
		{
			name: "Sucesso: Navegação para Subcomando válido",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Subcommands: map[string]*xcli.Command{
					"sub": {
						Name: "sub",
						Run: func(ctx *xclistruc.FlagValues) error {
							return nil
						},
					},
				},
			}),
			rawArgs: []string{"sub"},
			wantErr: false,
		},
		{
			name: "Sucesso: Ajuda disparada no subcomando pelos argumentos restantes",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Subcommands: map[string]*xcli.Command{
					"sub": {
						Name:             "sub",
						ShortDescription: "Sub cmd",
					},
				},
			}),
			rawArgs: []string{"sub", "--help"},
			wantErr: false,
		},
		{
			name: "Sucesso: Ajuda disparada no subcomando pela flag '-h' nos restantes",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Subcommands: map[string]*xcli.Command{
					"sub": {
						Name:             "sub",
						ShortDescription: "Sub cmd",
					},
				},
			}),
			rawArgs: []string{"sub", "-h"},
			wantErr: false,
		},
		{
			name: "Erro: Falha no ParseRawArgs (Fase 2 - argumento posicional inválido enviado como flag block)",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Run: func(ctx *xclistruc.FlagValues) error {
					return nil
				},
			}),
			rawArgs:       []string{"--flag", "value", "invalid_positional_here"},
			wantErr:       true,
			expectedError: "invalid argument: 'invalid_positional_here'",
		},
		{
			name: "Erro: Falha na validação de flags (Fase 3 - flag obrigatória ausente)",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Flags: []xclistruc.Flag{
					{LongName: "required-flag", Required: true},
				},
			}),
			rawArgs:       []string{"--other-flag=123"},
			wantErr:       true,
			expectedError: "[ERR] --required-flag : required",
		},
		{
			name: "Erro: Falha na validação de flags (Fase 3 - flag com tipo não suportado)",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Flags: []xclistruc.Flag{
					{LongName: "bad-flag", Type: xcliconstt.FlagType("bad-type")},
				},
			}),
			rawArgs:       []string{"--bad-flag=123"},
			wantErr:       true,
			expectedError: "unsupported flag type: 'bad-type'",
		},
		{
			name: "Erro: Falha na validação de flags (Fase 3 - falha interna do parserEngine)",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Flags: []xclistruc.Flag{
					{LongName: "fail-flag", Type: xcliconstt.FlagType("fail-type")},
				},
			}),
			rawArgs:       []string{"--fail-flag=123"},
			wantErr:       true,
			expectedError: "mock parsing failure",
		},
		{
			name: "Sucesso: Execução do hook Run do comando com flags hidratadas",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Flags: []xclistruc.Flag{
					{LongName: "username", Type: xcliconstt.FlagType("string")},
				},
				Run: func(ctx *xclistruc.FlagValues) error {
					return nil
				},
			}),
			rawArgs: []string{"--username=golang"},
			wantErr: false,
		},
		{
			name: "Sucesso: Comando sem Run ativa o TriggerHelp automaticamente",
			router: xcli.NewRouter(&xcli.Command{
				Name:             "root",
				ShortDescription: "Comando descritivo sem hook Run",
			}),
			rawArgs: []string{},
			wantErr: false,
		},
		{
			name: "Erro: O hook Run do comando retorna uma falha de negócio",
			router: xcli.NewRouter(&xcli.Command{
				Name: "root",
				Run: func(ctx *xclistruc.FlagValues) error {
					return errors.New("business logic failed")
				},
			}),
			rawArgs:       []string{},
			wantErr:       true,
			expectedError: "business logic failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.router.Run(tt.rawArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Router.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if err.Error() != tt.expectedError {
					var cliErr xerrors.IErrorCLI
					if errors.As(err, &cliErr) {
						if cliErr.GetDevMessage() != tt.expectedError && cliErr.Error() != tt.expectedError {
							t.Errorf("Router.Run() error message = %q, want %q", cliErr.GetDevMessage(), tt.expectedError)
						}
					} else if err.Error() != tt.expectedError {
						t.Errorf("Router.Run() error = %q, want %q", err.Error(), tt.expectedError)
					}
				}
			}
		})
	}
}
