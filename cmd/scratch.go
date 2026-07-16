package main

// import (
// 	"context"
// 	"log"
//
// 	"://github.com"
// 	"://github.com/ebitenutil"
// 	"://github.com/engine"
// )
//
// // CommandPayload bridges the thread gap.
// // It carries parsed logical conclusions out of Prolog back to the game thread.
// type CommandPayloadx struct {
// 	Action string
// 	Value  int
// }
//
// type Gamex struct {
// 	inputBuffer  []byte
// 	prologInput  chan []byte         // Channel sending raw bytes to Prolog thread
// 	gameCommands chan CommandPayload // Channel receiving parsed commands from Prolog thread
//
// 	// Game state variables updated by parsed commands
// 	playerSpeed  int
// 	playerHealth int
// }
//
// func NewGame() *Game {
// 	g := &Game{
// 		inputBuffer:  make([]byte, 0, 128),
// 		prologInput:  make(chan []byte, 10), // Buffered to prevent blocking input
// 		gameCommands: make(chan CommandPayload, 10),
// 		playerSpeed:  100,
// 		playerHealth: 100,
// 	}
//
// 	// Spin up the background thread immediately
// 	go g.prologWorker()
//
// 	return g
// }
//
// // Update handles logic at exactly 60 Ticks Per Second
// func (g *Game) Updatex() error {
// 	// 1. Capture keyboard text input from Ebitengine
// 	// (Ebitengine tracks typed runes every frame natively)
// 	runes := ebiten.AppendInputChars(nil)
// 	for _, r := range runes {
// 		b := byte(r)
// 		g.inputBuffer = append(g.inputBuffer, b)
//
// 		// 2. Check for the terminal character trigger
// 		if b == '=' {
// 			// Allocate or copy a standalone slice for the background worker
// 			// to prevent data races on the local inputBuffer
// 			commandCopy := make([]byte, len(g.inputBuffer))
// 			copy(commandCopy, g.inputBuffer)
//
// 			// Ship the bytes off the render thread instantly
// 			select {
// 			case g.prologInput <- commandCopy:
// 			default:
// 				// Dropped if worker queue is completely choked
// 			}
//
// 			// Clear the local input buffer for the next player command
// 			g.inputBuffer = g.inputBuffer[:0]
// 		}
// 	}
//
// 	// 3. Process completed Prolog actions safely on the main thread
// 	// Non-blocking poll reads everything ready this tick
// loop:
// 	for {
// 		select {
// 		case cmd := <-g.gamecommands:
// 			switch cmd.action {
// 			case "set_speed":
// 				g.playerspeed = cmd.value
// 			case "set_health":
// 				g.playerhealth = cmd.value
// 			}
// 		default:
// 			break loop // nothing left in the queue for this frame
// 		}
// 	}
//
// 	return nil
// }
//
// // prologWorker runs continuously on an isolated OS thread
// func (g *Game) prologWorkerx() {
// 	// Initialize interpreter completely separated from Ebitengine's execution thread
// 	p := prolog.New(nil, nil)
// 	_ = p.Exec(`
// 		parse_command(Action, Arg) --> command_name(Action), " ", integer_chars(ArgChars), "=", { number_codes(Arg, ArgChars) }.
// 		command_name(set_speed)   --> "SET_SPEED".
// 		command_name(set_health)  --> "SET_HEALTH".
// 		integer_chars([C|Cs])     --> [C], { C >= 48, C =< 57 }, integer_chars(Cs).
// 		integer_chars([])         --> [].
// 	`)
//
// 	termBuffer := make([]engine.Term, 0, 128)
//
// 	// Block until raw bytes are forwarded from the game loop
// 	for rawBytes := range g.prologInput {
// 		termBuffer = termBuffer[:0]
// 		for _, b := range rawBytes {
// 			termBuffer = append(termBuffer, engine.Integer(b))
// 		}
//
// 		sols, err := p.Query(context.Background(), "parse_command(Action, Arg, ?, []).", engine.List(termBuffer...))
// 		if err != nil {
// 			continue
// 		}
//
// 		if sols.Next() {
// 			var result struct {
// 				Action engine.Atom
// 				Arg    int
// 			}
// 			if err := sols.Scan(&result); err == nil {
// 				// Ship structural instructions back into the game engine queue
// 				g.gameCommands <- CommandPayload{
// 					Action: result.Action.String(),
// 					Value:  result.Arg,
// 				}
// 			}
// 		}
// 		sols.Close()
// 	}
// }
//
// // func main() {
// // 	ebiten.SetWindowSize(640, 480)
// // 	ebiten.SetWindowTitle("Ebitengine + Async Prolog Tokenizer")
// // 	if err := ebiten.RunGame(NewGame()); err != nil {
// // 		log.Fatal(err)
// // 	}
// // }
