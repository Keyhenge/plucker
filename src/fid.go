package plucker

import (
	"errors"

	"github.com/rminnich/go9p"
)

var (
	fsplumb *go9p.Clnt
	pfd     int = -1
	pfid    *go9p.Fid
)

func unmount() {
	if fsplumb != nil {
		fsplumb.Unmount()
	}
}

func open(name string, omode int) (int, error) {
	var err error

	if fsplumb == nil {
		fsplumb, err = go9p.Mount("plumb", "") //TODO: fsplumb = nsmount("plumb", "");
		if err != nil {
			return -1, err
		}
	}
	/*
	* It's important that when we send something,
	* we find out whether it was a valid plumb write.
	* (If it isn't, the client might fall back to some
	* other mechanism or indicate to the user what happened.)
	* We can't use a pipe for this, so we have to use the
	* fid interface.  But we need to return a fd.
	* Return a fd for /dev/null so that we return a unique
	* file descriptor.  In plumbsend we'll look for pfd
	* and use the recorded fid instead.
	 */

	if (omode & 3) == go9p.OWRITE {
		if pfd != -1 {
			return -1, errors.New("already have plumb send open")
		}
		pfd, err = go9p.Open("/dev/null", go9p.OWRITE) //TODO: pfd = open("/dev/null", OWRITE);
		if err != nil {
			return -1, err
		}
		pfid, err = go9p.FSOpen(fsplumb, name, omode) //TODO: pfid = fsopen(fsplumb, name, omode);
		if err != nil {
			go9p.Close(pfd) //TODO: close(pfd)
			pfd = -1        //REFAC?
			return -1, err
		}
		return pfd, nil
	}

	return go9p.FSOpenFD(fsplumb, name, omode) //TODO: return fsopenfd(fsplumb, name, omode);
}

func openFid(name string, mode int) (*go9p.Fid, error) {
	var err error

	if fsplumb == nil {
		fsplumb, err = go9p.Mount("plumb", "") //TODO: fsplumb = nsmount("plumb", "");
		if err != nil {
			return nil, err
		}
	}

	return go9p.FSOpen(fsplumb, name, mode) //TODO: pfid = fsopen(fsplumb, name, mode);
}

func sendToFid(fid *go9p.Fid, msg Message) error {
	if fid == nil {
		return errors.New("invalid fid")
	}
	buffer := packMessage(msg)
	return go9p.FSWrite(fid, buffer) //TODO: n = fswrite(fid, buf, n)
}

func send(fd int, msg Message) error {
	if fd == -1 {
		return errors.New("invalid fid")
	}
	if fd != pfd {
		return errors.New("fd is not the plumber")
	}

	return sendToFid(pfid, msg)
}

//TODO

func receive(fd int) (*Message, error) {
	return nil, nil
}

func receiveFid(fid *go9p.Fid) (*Message, error) {
	return nil, nil
}
