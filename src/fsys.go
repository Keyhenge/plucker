package plucker

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rminnich/go9p"
)

const (
	Stack = 32 * 1024
)

type Dirtab struct {
	name     string
	mimetype uint8    //might need better name/type; originally "type uchar"
	qid      uint     //?
	perm     uint     //?
	nopen    int      //* #fids open on this port
	fopen    *Fid     //?
	holdq    *Holdq   //?
	readq    *Readreq //?
	sendq    *Sendreq //?
}

type Fid struct {
	fid         int //file identifier
	busy        int //bool?
	open        int //bool?
	mode        int //uint?
	qid         Qid
	dir         *Dirtab //is this necessary?
	offset      int64   //* zeroed at beginning of each message, read or write
	writeBuffer []byte  //* partial message written so far; offset tells how much
	next        *Fid    //looks like linked list, replace with slice?
	nextOpen    *Fid    //?
}

type Readreq struct {
	fid    *Fid
	fcall  Fcall
	buffer []byte
	next   *Readreq //looks like linked list, replace with slice?
}

type Sendreq struct {
	nfid  int   //* number of fids that should receive this message
	nleft int   //* number left that haven't received it
	fid   []Fid //* fid[nfid]
	msg   *Message
	pack  string   //* packMessage()'d message'
	npack int      //* length of pack //unnecessary
	next  *Sendreq //looks like linked list, replace with slice?
}

type Holdq struct {
	msg  *Message
	next *Holdq //looks like linked list, replace with slice?
}

//struct	/* needed because incref() doesn't return value */
//{
//	Lock	lk;
//	int	ref;
//} rulesref;

const (
	NDIR  = 50
	Nhash = 16

	Qdir   = iota
	Qrules = iota
	Qsend  = iota
	Qport  = iota
	NQID   = Qport //unnecessary?

	//TODO
	//static Dirtab dir[NDIR] =
	//{
	//	{ ".",			QTDIR,	Qdir,			0500|DMDIR },
	//	{ "rules",		QTFILE,	Qrules,		0600 },
	//	{ "send",		QTFILE,	Qsend,		0200 }
	//};
)

var (
	srvfd       int
	clock       int
	fids        *[Nhash]Fid
	readLock    QLock                  //This is probably replaced by sync.Mutex
	queue       QLock                  //This is probably replaced by sync.Mutex
	messageSize       = 8192 + IOHDRSZ //const? IOHDRSZ?
	ndir        uint  = NQID           //number of entries in dir? might be unnecessary, just len(dir)
	nports      uint  = 0              //unnecessary
	ports       []string

	dir []Dirtab
)

var (
	ErrBadFCall     error = errors.New("bad fcall type")
	ErrPerm         error = errors.New("permission denied")
	ErrNoMemory     error = errors.New("malloc failed for buffer") //probably unnecessary
	ErrNotDir       error = errors.New("not a directory")
	ErrNoExist      error = errors.New("plumb file does not exist")
	ErrIsDir        error = errors.New("file is a directory")
	ErrBadMessage   error = errors.New("bad plumb message format")
	ErrNoSuchPort   error = errors.New("no such plumb port") //rename these 2?
	ErrNoPort       error = errors.New("couldn't find destination for message")
	ErrInUse        error = errors.New("file already open")
	ErrTooManyPorts error = errors.New(fmt.Sprintf("plumb: too many ports; max %d\n", NDIR)) //This might bake a number into the error too early
)

//* Add a new port. A no-op if port already exists or is an empty string
//TODO: Return error? Silent fail seems bad
func addPort(port string) {
	var index uint

	if port == "" {
		return
	}
	for index = NQID; index < ndir; index++ {
		if port == dir[index].name {
			return
		}
	}
	if index == NDIR {
		fmt.Fprint(os.Stderr, ErrTooManyPorts.Error())
		return
	}

	//ndir += 1
	dir[index].name = port
	dir[index].qid = index
	dir[index].perm = 0400 // Need more info on this, 0400 seems like it should be in a different base
	//nports += 1
	ports = append(ports, dir[index].name)
}

func startFsys(foreground bool) error {
	in, out := io.Pipe() //* In is server end, Out is client end

	//fmtinstall("F", fcallfmt) // TODO
	//clock = time(0) //Why?
	if err := go9p.post9pservice(out, "plumb", nil); err != nil {
		return err
	}
	err := out.Close()
	if err != nil {
		return err
	}
	if foreground {
		fsysProc(nil)
	} else {
		procCreate(fsysProc, nil, Stack)
	}
}

//originally static void fsysproc(void *v), current parameter seems wrong
func fsysProc(v func()) {
	//USED(v) //(void)(v)
	//initFCall() // TODO

	var t go9p.Fcall
	for { // REFAC: Go std library probably has something

	}
}
