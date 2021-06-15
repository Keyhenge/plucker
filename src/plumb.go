package plucker

import (
	"flag"
	"fmt"
	"os"

	"github.com/rminnich/go9p"
)

//TODO: Get STDIN here
func gather() []byte {
	return []byte{}
}

func main() {
	plumbfileStr := flag.String("p", "", "plumbfile")
	attrStr := flag.String("a", "", "'attr=value ...'")
	sourceStr := flag.String("s", "", "source")
	destStr := flag.String("d", "", "destination")
	typeStr := flag.String("t", "text", "type")
	wdirStr := flag.String("w", "~", "working directory")
	stdinOpt := flag.Bool("i", false, "read from STDIN")

	var fd int
	var err error

	if *plumbfileStr != "" {
		fd, err = go9p.Open(plumbfile, go9p.OWRITE)
	} else {
		fd, err = open("send", go9p.OWRITE)
	}
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	attrList, err := unpackAttributeList(*attrStr)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	msg := Message{
		source:   *sourceStr,
		dest:     *destStr,
		mimetype: *typeStr,
		wdir:     *wdirStr,
		attrList: attrList,
	}

	if *stdinOpt {
		//TODO: Get STDIN here
		msg.data = gather()

		_, err := lookup(msg.attrList, "action")
		if err != nil {
			actionAttr, err := unpackAttribute("action=showdata")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			msg.attrList = append(msg.attrList, actionAttr)
		}
		if err = send(fd, msg); err != nil {
			fmt.Printf("%s\n", err)
		}
		return
	}
	for index := 0; index < len(os.Args); index++ {
		msg.data = []byte(os.Args[index])
		if err = send(fd, msg); err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
