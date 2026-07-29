package main

import (
	// "fmt"
	"github.com/trealla-prolog/go/trealla"
)

const AtomInvalid = trealla.Atom("invalid")
const AtomConnectTrue = trealla.Atom("connect_true")
const AtomConnectFalse = trealla.Atom("connect_false")
const AtomList = trealla.Atom("list")
const AtomFiles = trealla.Atom("files")
const AtomPrograms = trealla.Atom("programs")
const AtomNetworks = trealla.Atom("networks")
const AtomRoy1Fn = trealla.Atom("roy_1_fn")
const AtomRoy2Fn = trealla.Atom("roy_2_fn")
const AtomRoy3Fn = trealla.Atom("roy_3_fn")
const AtomBreach = trealla.Atom("breach")
const AtomMem = trealla.Atom("mem")
const AtomMemCopy = trealla.Atom("memcopy")
const AtomLobby = trealla.Atom("lobby")

// var CommandBytes = make([][]byte, len(AtomCommandIds))
var CommandBytes = [][]byte{
	[]byte{},
	[]byte{},
	[]byte{},
	[]byte{},

	[]byte("[Files]:"),
	[]byte("[Programs]: "),
	[]byte("[Networks]:"),

	[]byte("[Email]: Sorry About Breakfast"),
	[]byte("[Email]: RE: Sorry About Breakfast"),
	[]byte("[Email]: RE: RE: Sorry About Breakfast"),

	[]byte("`BREACH {name}`: Hack into Network"),
	[]byte("[Lobby]: connection"),
}

//CommandBytes[AtomCommandIds[AtomConnectRoy1Fn]] = []byte{}

// var CommandBytesIdx = []int{0, 0, 0, 0, 29, 62}
// var CommandBytesLen = []int{0, 0, 0, 28, 33, 0}
//
// var CommandBytes = []byte(
// 	"Email: Sorry About Breakfast" +
// "Email: RE: Sorry About Breakfast" +
// 		"Email: RE: RE: Sorry About Breakfast",
// )
