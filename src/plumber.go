package plucker

import (
	"flag"
	"fmt"
	"os"

	"github.com/rminnich/go9p"
)

//* Objects
const (
	OArg   = iota
	OAttr  = iota
	OData  = iota
	ODest  = iota
	OPlumb = iota
	OSrc   = iota
	OType  = iota
	OWdir  = iota
)

//* Verbs
const (
	VAdd     = iota //* apply to OAttr only
	VClient  = iota
	VDelete  = iota //* apply to OAttr only
	VIs      = iota
	VIsDir   = iota
	VIsFile  = iota
	VMatches = iota
	VSet     = iota
	VStart   = iota
	VTo      = iota
)

type Rule struct {
	obj   int
	verb  int
	arg   string //* unparsed string of all arguments
	qarg  string //* quote-processed arg string
	regex Reprog
}

type Ruleset struct {
	npat int    //?
	nact int    //?
	pat  **Rule //?
	act  **Rule //?
	port string
}

type Exec struct {
	msg           Message
	match         [10]byte //This might not be []byte
	matchBegin    int      //p0?
	matchEnd      int      //p1?
	clearClick    int      //* click was expanded; remove attribute //bool?
	setData       int      //* data should be set to $0 //bool?
	holdForClient int      //* exec'ing client; keep message until port is opened //bool?
	file          string
	dir           string
}

var (
	debug       int     // bool?
	foreground  int = 0 // bool?
	plumbfile   string
	user        string
	home        string
	progName    string
	rules       **Ruleset
	printErrors int    = 1 // bool?
	parsejmp    string     //original type: jmp_buf, TODO
	lastError   string     //error? unnecessary?
)

func makePorts(rules *[]Ruleset) {
	for i := 0; rules[i]; i++ {
		addPort(rules[i].port)
	}
}

func main() {
	var buffer [512]byte

	progname = "plumber"

	debugOpt := flag.Bool("d", false, "debug")
	foregroundOpt := flag.Bool("f", false, "run in foreground")
	plumbfileStr := flag.String("p", "", "plumbfile")

	user := "foo" //TODO
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	if *plumbfileStr == "" {
		*plumbfileStr = fmt.Sprintf("%s/.config/plumbing")
	}

	fd, err := go9p.Open(plumbfile, go9p.OREAD)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	//if setjmp(parsejmp) {
	//	fmt.Printf("parse error\n")
	//	return
	//}

	rules = readRules(*plumbfileStr, fd)
	go9p.Close(fd)

	printErrors = 0
	makePorts(rules)
	startFSys(*foregroundOpt)
}
