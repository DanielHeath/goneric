package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"os"
	"strings"
)

func main() {
	gfs := GonericFs{}
	path, done, err := gfs.Serve()
	if err != nil {
		panic(err)
	}
	_ = path
	select {
	case <-done:
		break
	}
}

type GonericFs struct {
	pathfs.FileSystem // Optional, defaults to pathfs.NewDefaultFileSystem
}

const BASEPATH = "src/goneric4"

func (me *GonericFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Println("getAttr")
	log.Println(name)
	if strings.HasPrefix(name, BASEPATH) {
		return nil, fuse.ENOENT
	}
	fi, err := os.Stat("src/" + name)
	if err != nil {
		log.Println("cant getattr there")
		log.Println(err)
		return nil, fuse.ENOENT
	}
	if fi.IsDir() {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}
	return &fuse.Attr{
		Mode: fuse.S_IFREG | 0644, Size: uint64(fi.Size()),
	}, fuse.OK
}

func (me *GonericFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Println("openDir")
	log.Println(name)
	if strings.HasPrefix(name, BASEPATH) {
		return nil, fuse.ENOENT
	}
	file, err := os.Open("src/" + name)
	if err != nil {
		log.Println("no opening that")
		log.Println(err)
		return nil, fuse.ENOENT
	}

	fi, err := file.Readdir(0) // All the files! Don't try this on a wide FS.
	if err != nil {
		log.Println("cant read dir listing")
		log.Println(err)
		return nil, fuse.ENOENT
	}

	entries := make([]fuse.DirEntry, len(fi))
	for i, file := range fi {
		var mode uint32
		mode = fuse.S_IFREG
		if file.IsDir() {
			mode = fuse.S_IFDIR
		}
		entries[i] = fuse.DirEntry{Name: file.Name(), Mode: mode}
	}
	return entries, fuse.OK
}

func (me *GonericFs) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	log.Println("open")
	log.Println(name)
	if strings.HasPrefix(name, BASEPATH) {
		return nil, fuse.ENOENT
	}
	f, err := os.Open("src/" + name)
	if err != nil {
		log.Println("cant read that file")
		log.Println(err)
		return nil, fuse.ENOENT
	}
	return nodefs.NewLoopbackFile(f), fuse.OK
}

func (me *GonericFs) Serve() (path string, done chan bool, err error) {
	if me.FileSystem == nil {
		me.FileSystem = pathfs.NewDefaultFileSystem()
	}

	nfs := pathfs.NewPathNodeFs(me, nil)
	err = os.MkdirAll(BASEPATH, os.ModePerm)
	if err != nil {
		return
	}
	path = BASEPATH
	server, _, err := nodefs.MountRoot(path, nfs.Root(), nil)
	if err != nil {
		return
	}
	done = make(chan bool)
	go func() {
		server.Serve()
		os.RemoveAll(path) // Tidy up once finished.
		done <- true
	}()
	return
}
