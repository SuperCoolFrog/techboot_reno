package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	// "github.com/trealla-prolog/go/trealla"
	"image/color"
	"techboot_reno/cmd/assets"
)

const (
	// Grid ends up being 26,19
	screenWidth   = 640 * 2
	screenHeight  = 480 * 2
	fontSize      = 24.0 * 2.0
	screenPadding = fontSize / 4
)

// var fontColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // this is pure green
// var fontColor = color.RGBA{R: 100, G: 255, B: 150, A: 255} // Minty green
// var fontColor = color.RGBA{R: 112, G: 238, B: 238, A: 255} // Muted Cyan
var fontColor = color.RGBA{R: 51, G: 255, B: 51, A: 255} // Softened phospor green
var fontColorComment = color.RGBA{R: 150, G: 150, B: 150, A: 255 * 0.6}

//go:embed fonts/Courierprime_1OVL.ttf
var fontBytes []byte
var fontSrc *text.GoTextFaceSource

//go:embed parser.pl
var parserpl []byte

//go:embed puzzles.pl
var puzzlespl []byte

func init() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}
	fontSrc = src
}

type Game struct {
	State                  GameState
	inputRunes             []rune
	ans                    *AnimationSystem
	bs                     *BufferSystem
	gs                     *GridSystem
	pz                     *PuzzleSystem
	ps                     *PathSystem
	jxsp                   *JunctionSystem //Scene -> PuzzleIds
	jxpp                   *JunctionSystem //PuzzleId -> PathIds
	jxgp                   *CompoundSystem //GateId+GateType -> []Paths
	pcn                    *FlagSystem     // PuzzleConnected
	pac                    *FlagSystem     // PuzzleAnimationComplete
	gap                    *FlagSystem     // GateAnimationPlaying
	gac                    *FlagSystem     // GateAnimationCompleted
	b                      Bucket
	MouseMoved             bool
	LastMouseX, LastMouseY int
	Exit                   bool
	parserpl               string
	puzzlespl              string
	prologInput            chan []byte          // Channel sending raw bytes to Prolog thread
	prologOutput           chan CommandResponse // Channel receiving parsed commands from Prolog thread
}

func (g *Game) Update() error {

	mx, my := ebiten.CursorPosition()
	g.MouseMoved = mx != g.LastMouseX || my != g.LastMouseY
	g.LastMouseX = mx
	g.LastMouseY = my

	// Update logic if needed
	g.inputRunes = ebiten.AppendInputChars(g.inputRunes[:0])

	// if len(g.inputRunes) > 0 {
	// 	fmt.Printf("rune: %v\n", g.inputRunes)
	// }

	g.UpdateState()
	g.ans.Update()

	if g.Exit {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	screen.Fill(color.RGBA{R: 18, G: 18, B: 18, A: 255}) // matte black
	// screen.Fill(color.RGBA{R: 26, G: 27, B: 38, A: 255}) // storm blue grey 26, 27, 38
	// screen.Fill(color.RGBA{R: 20, G: 24, B: 30, A: 255}) // a little off black
	// screen.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	// screen.Fill(color.RGBA{R: 0, G: 0, B: 50, A: 255})

	// err := g.Grid.RenderDebug(screen)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	g.gs.Render(screen)
	g.ans.Render(screen, g.gs)

	// Test
	// op := &ebiten.DrawImageOptions{}
	// screen.DrawImage(assets.Images[assets.SpriteIDCircle], op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	windowSizeW, windowSizeH := screenWidth, screenHeight
	ebiten.SetWindowSize(windowSizeW, windowSizeH)
	ebiten.SetWindowTitle("Techboot Reno - Cyber Crawler")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60) // Locks Update cycles to 60Hz natively

	// This is a lot probably could tweak it once I have an idea of total grids
	// const MaxTotalCells = 100_000
	// const MaxGrids = 50

	// Tweaking values
	const MaxTotalCells = 150_000
	const MaxGrids = 15

	// MaxMarkers calculated by 2(cols)+2(rows) of output grid

	puzzleSystem := NewPuzzleSystem(10, 26*2+36*2, 3, 0, 0, 0)

	maxIdJxgp := puzzleSystem.TotalGates
	if int(GateTypeCount) > maxIdJxgp {
		maxIdJxgp = int(GateTypeCount)
	}

	game := &Game{
		State:        Scene3_Init, //Scene1_Init,
		gs:           NewGridSystem(MaxTotalCells, MaxGrids),
		ans:          NewAnimationSystem(),
		bs:           NewBufferSystem(500_000, 10),
		ps:           NewPathSystem(200),
		pz:           puzzleSystem,
		jxsp:         NewJunctionSystem(int(GameStateCount), 5),
		jxpp:         NewJunctionSystem(puzzleSystem.TotalPuzzles, 20),
		jxgp:         NewCompoundSystem(maxIdJxgp, 5),
		pcn:          NewFlagSystem(puzzleSystem.TotalPuzzles),
		pac:          NewFlagSystem(puzzleSystem.TotalPuzzles),
		gap:          NewFlagSystem(puzzleSystem.TotalGates),
		gac:          NewFlagSystem(puzzleSystem.TotalGates),
		parserpl:     string(parserpl),
		puzzlespl:    string(puzzlespl),
		prologInput:  make(chan []byte, 128), // Buffered to prevent blocking input
		prologOutput: make(chan CommandResponse, 128),
	}

	assets.Load() // Load Assets before init bucket

	game.b = InitBucketItems(game)

	GeneratePuzzles(game)

	// Initialize parser
	go game.prologWorker()

	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		if err != nil {
			log.Fatal(err)
		}
	}
}
