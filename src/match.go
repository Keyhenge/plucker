package plucker

import (
	"bytes"
	"fmt"
)

// TODO: rename to matchObj?
func verbIs(obj int, msg Message, r Rule) (bool, error) {
	switch obj {
	case OData:
		return bytes.Compare(msg.data, []byte(r.qarg)) == 0, nil
	case OSrc:
		return msg.source == r.qarg, nil
	case ODest:
		return msg.dest == r.qarg, nil
	case OType:
		return msg.mimetype == r.qarg, nil
	case OWdir:
		return msg.wdir == r.qarg, nil
	default:
		return false, fmt.Errorf("unimplemented 'is' object %d\n", obj)
	}
}

//TODO
func setVar(rs [10]Resub, match [10]byte) {
	for i := 0; i < 10; i += 1 {
		match[i] = 0
	}
	for i := 0; i < 10 && rs[i].s.sp != ""; i += 1 {
		match[i] = rs[i].s.sp
	}
}

//TODO
func clickMatch(re *Reprog, text string, rs [10]Resub, click int) {
	//* Click is in characters, not bytes
	for i := 0; i < click; i += 1 {

	}
}

// TODO
func verbMatches(obj int, msg Message, r Rule, e Exec) (bool, error) { return false, nil }

// TODO, seems unnecessary
func isFile(file string, maskon uint64, maskoff uint64)

// TODO, seems wrong
func absolute(dir string, file string) string {
	if file[0] == '/' {
		return file
	}
	return dir + file
}

// TODO
func verbIsFile(obj int, msg Message, r Rule, e Exec, maskon uint64, maskoff uint64, vars []string) (bool, error) {
	return false, nil
}

// TODO
func verbSet(obj int, msg Message, r Rule, e Exec) (bool, error) { return false, nil }

// TODO
func verbAdd(obj int, msg Message, r Rule, e Exec) (bool, error) { return false, nil }

// TODO
func verbDelete(obj int, msg Message, r Rule, e Exec) (bool, error) { return false, nil }

// TODO, should rename to matchVerb
func matchPat(msg Message, r Rule, e Exec) (bool, error) {
	switch r.verb {
	case VAdd:
		return verbAdd(r.obj, msg, r, e)
	case VDelete:
		return verbDelete(r.obj, msg, r, e)
	case VIs:
		return verbIs(r.obj, msg, r, e)
	case VIsDir:
		return verbIsDir(r.obj, msg, r, e, DMDIR, 0, e.dir)
	case VIsFile:
		return verbIsFile(r.obj, msg, r, e, ^DMDIR, DMDIR, e.file)
	case VMatches:
		return verbMatches(r.obj, msg, r, e)
	case VSet:
		return verbSet(r.obj, msg, r, e)
	default:
		return false, fmt.Errorf("unimplemented verb %d\n", r.verb)
	}
}

// TODO
func rewrite(msg Message, e Exec)

// TODO?
func buildargv(s string, e Exec)

// TODO
func matchRuleset(msg Message, rs Ruleset) Exec

const (
	NARGS		= 100
	NARGCHAR	= 8*1024
	EXECSTACK 	= 32*1024+(NARGS+1)*sizeof(char*)+NARGCHAR
)

//* copy argv to stack and free the incoming strings, so we don't leak argument vectors
// TODO?
func stackargv(inargv []string, argv [NARGS+1]string, args string)

// TODO: probably unnecessary
//func execProc(v *void)

// TODO
func startup(rs Ruleset, e Exec) string
