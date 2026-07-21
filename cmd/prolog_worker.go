package main

import (
	"context"
	"fmt"

	"github.com/trealla-prolog/go/trealla"
	"strings"
)

type OutResult struct {
	trealla.Functor `prolog:"/1"`
	Result          trealla.Atom
}

type Result struct {
	Out OutResult `prolog:"Out"`
}

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
		cleanedBytes := rawBytes[0 : len(rawBytes)-1] //Remove trailing =
		cleanedStr := strings.TrimSpace(string(cleanedBytes))
		if len(cleanedStr) == 0 {
			continue
		}

		// Keep your working Go downcasing logic intact
		normalizedStr := strings.ToLower(cleanedStr)

		csvStr := strings.ReplaceAll(normalizedStr, " ", ",")
		queryStr := fmt.Sprintf("process_command([%s,3,'Out'], Out).", csvStr)

		fmt.Printf("Processing queryStr: '%s'\n", queryStr)

		q, err11 := pl.QueryOnce(ctx, queryStr)
		if err11 != nil {
			panic(err11)
		}

		// val := q.Solution["Out"]
		// fmt.Printf("Once query valid: %+v ; %T\n", val, val)

		var r Result

		if errv := q.Solution.Scan(&r); errv != nil {
			panic(errv)
		}

		// fmt.Printf("Cast %s; %T \n", r.Out.Result.String(), r.Out.Result)
		// if r.Out.Result == trealla.Atom("connect_true") {
		// 	fmt.Printf("Matches atom")
		// }

		g.prologOutput <- r.Out.Result

		// if query.Next(ctx) {
		// 	answer := query.Current()

		// 	// Extract variables using the safe trealla.Atom mapping type
		// 	actionAtom, ok1 := answer.Solution["Action"].(trealla.Atom)
		// 	argResult, ok2 := answer.Solution["Arg"].(int64)

		// 	if ok1 && ok2 {
		// 		actionStr := string(actionAtom)
		// 		fmt.Printf("🎯 SUCCESS! Action: %s, Arg: %d\n", actionStr, argResult)

		// 		g.prologOutput <- CommandPayload{
		// 			Action: actionStr,
		// 			Value:  int(argResult),
		// 		}
		// 	} else {
		// 		fmt.Printf("❌ Failed to parse output mapping: Action=%T, Arg=%T\nFor: %v\n",
		// 			answer.Solution["Action"], answer.Solution["Arg"], answer)
		// 	}
		// } else {
		// 	fmt.Printf("❌ Prolog: Embedded rules rejected the command format.\n")
		// }
		// query.Close()
	}
}
