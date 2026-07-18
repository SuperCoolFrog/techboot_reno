package main

import (
	"context"
	"fmt"

	"github.com/trealla-prolog/go/trealla"
	"strings"
)

func (g *Game) prologWorker() {
	// Initialize a clean, zero-dependency Trealla WebAssembly instance
	pl, err := trealla.New()
	if err != nil {
		fmt.Printf("Failed to boot Trealla: %v\n", err)
		return
	}

	ctx := context.Background()

	if err := pl.ConsultText(ctx, "user", g.parserpl); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Embedded rules compiled successfully and ready for pre-query loading!\n")

	for rawBytes := range g.prologInput {
		cleanedBytes := rawBytes[0:len(rawBytes)-1] //Remove trailing =
		cleanedStr := strings.TrimSpace(string(cleanedBytes))
		if len(cleanedStr) == 0 {
			continue
		}


		// Keep your working Go downcasing logic intact
		normalizedStr := strings.ToLower(cleanedStr)
		fmt.Printf("Processing Input: '%s'\n", normalizedStr)

		// only working like this; bind does not work, possibly expects atom
		query := pl.Query(ctx,
			// "parse_command('"+normalizedStr+"', Action, Arg).",
			// "phrase(lex(Tokens),  [h,e,l,l,o,' ', w,o,r,l,d]).",
			// "lex2([h,e,l,l,o,' ', w,o,r,l,d], Tokens).",
			"parse_command('"+normalizedStr+"', Action, Arg).",
		)

		q, err11 := pl.QueryOnce(ctx, "sentence([the, cat, chases, a, mouse], Out).")
		if err11 != nil {
			panic(err11)
		} else {
			fmt.Printf("Once query valid: %v\n", q)
		}

		if query.Next(ctx) {
			answer := query.Current()

			// Extract variables using the safe trealla.Atom mapping type
			actionAtom, ok1 := answer.Solution["Action"].(trealla.Atom)
			argResult, ok2 := answer.Solution["Arg"].(int64)

			if ok1 && ok2 {
				actionStr := string(actionAtom)
				fmt.Printf("🎯 SUCCESS! Action: %s, Arg: %d\n", actionStr, argResult)

				g.prologOutput <- CommandPayload{
					Action: actionStr,
					Value:  int(argResult),
				}
			} else {
				fmt.Printf("❌ Failed to parse output mapping: Action=%T, Arg=%T\nFor: %v\n",
					answer.Solution["Action"], answer.Solution["Arg"], answer)
			}
		} else {
			fmt.Printf("❌ Prolog: Embedded rules rejected the command format.\n")
		}
		query.Close()
	}
}
