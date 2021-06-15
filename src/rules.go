package plucker

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rminnich/go9p"
)

type Input struct {
	file string //* name of file
	fd   BioBuf //* input buffer, if from real file //REFAC
	s    []byte //* input string, if from /mnt/plumb/rules
	//end []byte //* end of input string //Unnecessary probably
	lineno int
	next   *Input //Linked list behavior, consider REFAC
}

type Var struct {
	name   string
	value  string //probably string?
	qvalue string //probably string?
}

var (
	parsing int
	nvars   int
	vars    *Var   //this might need to be an array
	input   *Input //linked list
	ebuf    [4096]byte

	badports []string = []string{
		".",
		"..",
		"send",
		"",
	}

	objects []string = []string{
		"arg",
		"attr",
		"data",
		"dst",
		"plumb",
		"src",
		"type",
		"wdir",
		"", //This might not be correct
	}

	verbs []string = []string{
		"add",
		"client",
		"delete",
		"is",
		"isdir",
		"isfile",
		"matches",
		"set",
		"start",
		"to",
		"", //This might not be correct
	}
)

// Consider REFAC that iterates over slice instead of linked list
func printInputStackRev(in *Input) {
	if in == nil {
		return
	}
	printInputStackRev(in.next)
	fmt.Printf("%s:%d: ", in.file, in.lineno)
}

func printInputStack() {
	printInputStackRev(input)
}

func pushInput(name string, fd int, str []byte) error {
	depth := 0
	for in := input; in != nil; in = input.next {
		depth += 1
		if depth >= 10 { //* prevent deep C stack in plumber and bad include structure
			return errors.New("include stack too deep; max 10") //Would this be necessary with a slice?
		}
	}

	newInput := &Input{
		file: name,
		next: input,
	}
	input = newInput

	if len(str) > 0 {
		newInput.s = str
	} else {
		return Binit(newInput.fd, fd, go9p.OREAD) //TODO
	}

	return nil
}

func popInput() {

}
