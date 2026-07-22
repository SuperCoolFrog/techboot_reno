package main

import "github.com/trealla-prolog/go/trealla"

const AtomInvalid = trealla.Atom("invalid")
const AtomConnectTrue = trealla.Atom("connect_true")
const AtomConnectFalse = trealla.Atom("connect_false")
const AtomConnectList = trealla.Atom("list")
const AtomConnectRoy1Fn = trealla.Atom("roy_1_fn")
const AtomConnectRoy2Fn = trealla.Atom("roy_2_fn")
const AtomConnectRoy3Fn = trealla.Atom("roy_3_fn")

type CommandId uint32

var AtomCommandIds = map[trealla.Atom]CommandId{
	AtomInvalid:       0,
	AtomConnectTrue:   1,
	AtomConnectFalse:  2,
	AtomConnectRoy1Fn: 3,
	AtomConnectRoy2Fn: 4,
	AtomConnectRoy3Fn: 5,
}

var CommandBytesIdx = []int{0, 0, 0, 0, 28, 0}
var CommandBytesLen = []int{0, 0, 0, 28, 0, 0}

var CommandBytes = []byte(
	`Email: Sorry About Breakfast`,
)
