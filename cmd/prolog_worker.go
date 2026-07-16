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
	pl, err := trealla.New()
	if err != nil {
		fmt.Printf("Failed to boot Trealla: %v\n", err)
		return
	}

	// 1. THE DEFINITIVE FIX: Inline assert the rule.
	// This completely avoids ConsultText filing/parsing issues.
	ctxInit := context.Background()
	initQuery := pl.Query(ctxInit,
		`asserta((
			parse_command(InputStr, Action, Arg) :-
				atom_string(InputAtom, InputStr),
				atom_concat(CleanedCommand, '=', InputAtom),
				atomic_list_concat([ActionAtom, ArgAtom], ' ', CleanedCommand),
				downcase_atom(ActionAtom, Action),
				atom_number(ArgAtom, Arg)
		)).`,
	)
	if initQuery.Next(ctxInit) {
		fmt.Printf("✅ Prolog Parser Rule committed to memory successfully!\n")
	} else {
		fmt.Printf("❌ Failed to register Prolog rules.\n")
	}
	initQuery.Close()

	// 2. Resume your production channel processing block
	for rawBytes := range g.prologInput {
		cleanedStr := strings.TrimSpace(string(rawBytes))
		if len(cleanedStr) == 0 {
			continue
		}

		fmt.Printf("Processing Input: '%s'\n", cleanedStr)
		ctx := context.Background()

		query := pl.Query(ctx,
			"parse_command(RawInput, Action, Arg).",
			trealla.WithBind("RawInput", cleanedStr),
		)

		if query.Next(ctx) {
			answer := query.Current()

			actionResult, ok1 := answer.Solution["Action"].(string)
			argResult, ok2 := answer.Solution["Arg"].(int64)

			if ok1 && ok2 {
				fmt.Printf("🎯 SUCCESS! Action: %s, Arg: %d\n", actionResult, argResult)

				g.prologOutput <- CommandPayload{
					Action: actionResult,
					Value:  int(argResult),
				}
			} else {
				fmt.Printf("❌ Type mapping mismatch: Action=%T, Arg=%T\n",
					answer.Solution["Action"], answer.Solution["Arg"])
			}
		} else {
			fmt.Printf("Prolog: Command input layout did not match your syntax definitions.\n")
		}

		if err := query.Err(); err != nil {
			fmt.Printf("⚠️ Engine Runtime Error: %v\n", err)
		}
		query.Close()
	}
}
