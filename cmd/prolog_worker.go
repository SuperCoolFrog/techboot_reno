package main

import (
	"context"
	"fmt"

	"github.com/trealla-prolog/go/trealla"
	"strings"
)

const parserpl2 = `
parse_command(InputStr, Action, Arg) :-
    atom_string(InputAtom, InputStr),
    atom_concat(CleanedCommand, '=', InputAtom),
    atomic_list_concat([ActionAtom, ArgAtom], ' ', CleanedCommand),
    downcase_atom(ActionAtom, Action),
    atom_number(ArgAtom, Arg).

`

func (g *Game) prologWorker() {
	// 1. Boot up a clean, zero-dependency Trealla WebAssembly sandbox instance
	pl, err := trealla.New()
	if err != nil {
		fmt.Printf("Failed to boot Trealla: %v\n", err)
		return
	}

	fmt.Printf("✅ Trealla Prolog engine ready for high-throughput string parsing!\n")

	// 2. Main consumer channel loop reading raw bytes from your Ebitengine render thread
	for rawBytes := range g.prologInput {
		cleanedStr := strings.TrimSpace(string(rawBytes))
		if len(cleanedStr) == 0 {
			continue
		}

		// CRITICAL FIX: Downcase the string in Go first!
		// This translates "CONNECT 42=" directly to "connect 42="
		normalizedStr := strings.ToLower(cleanedStr)

		fmt.Printf("Processing Input: '%s'\n", normalizedStr)
		ctx := context.Background()

		// THE FINAL COMPACT ISO QUERY METHOD:
		// 1. append/3 slices the incoming lowercased character array seamlessly
		// 2. atom_chars/2 and atom_number/2 re-bundle the components into standard output fields
		query := pl.Query(ctx,
			`append(ActionChars, [' '|ValAndEquals], RawInput),
			 append(ArgChars, ['='], ValAndEquals),
			 atom_chars(Action, ActionChars),
			 atom_chars(ArgAtom, ArgChars),
			 atom_number(ArgAtom, Arg).`,
			trealla.WithBind("RawInput", normalizedStr), // Trealla maps this straight to a lowercase character list!
		)

		if query.Next(ctx) {
			answer := query.Current()

			// FIX: Extract the Action variable as a trealla.Atom instead of a primitive string
			actionAtom, ok1 := answer.Solution["Action"].(trealla.Atom)
			argResult, ok2 := answer.Solution["Arg"].(int64)

			if ok1 && ok2 {
				// Convert the custom Trealla Atom type cleanly back into a standard Go string
				actionStr := string(actionAtom)
				fmt.Printf("🎯 SUCCESS! Action: %s, Arg: %d\n", actionStr, argResult)

				// Route structural execution tokens straight back into your main update thread loop
				g.prologOutput <- CommandPayload{
					Action: actionStr,
					Value:  int(argResult),
				}
			} else {
				fmt.Printf("❌ Failed to parse output mapping: Action=%T, Arg=%T\n",
					answer.Solution["Action"], answer.Solution["Arg"])
			}
		}

		// Catch driver or internal runtime thread failures cleanly
		if err := query.Err(); err != nil {
			fmt.Printf("⚠️ Engine Runtime Error: %v\n", err)
		}
		query.Close()
	}
}
