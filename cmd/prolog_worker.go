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

	_, _ = pl.QueryOnce(ctx, "dynamic(parse_command/3).")

	// Normalize line endings safely
	// rulesStr := strings.ReplaceAll(g.parserpl, "\r\n", "\n")

	//query := fmt.Sprintf("assertz(parse_command(test, connect, 42)).", g.parserpl)
	// query := fmt.Sprintf("%s", g.parserpl)
	_, err = pl.QueryOnce(ctx, g.parserpl)
	if err != nil {
		fmt.Printf("[X] Embedded rules failed: %v\n", err)
		panic(err)
	}

	fmt.Printf("✅ Embedded rules compiled successfully and ready for pre-query loading!\n")

	for rawBytes := range g.prologInput {
		cleanedStr := strings.TrimSpace(string(rawBytes))
		if len(cleanedStr) == 0 {
			continue
		}

		// Keep your working Go downcasing logic intact
		normalizedStr := strings.ToLower(cleanedStr)
		fmt.Printf("Processing Input: '%s'\n", normalizedStr)

		// only working like this; bind does not work, possibly expects atom
		query := pl.Query(ctx,
			// "parse_command(test, Action, Arg).",
			"parse_command('"+normalizedStr+"', Action, Arg).",
			// "parse_command(RawInput, Action, Arg).",
			// trealla.WithBind("RawInput", "'"+normalizedStr+"'"),
			// "parse_command('connect', Action, Arg).",
		)

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
