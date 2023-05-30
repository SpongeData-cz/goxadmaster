package goxadmaster

/*
#cgo LDFLAGS: -lXADMaster
#include <libXADMaster.h>
#include <stdlib.h>
inline void EntrySetRenaming(Entry * self, const char * renaming);
*/
import "C"
import (
	"errors"
	"fmt"
)

type Entry interface {
	GetFilename() string
	GetDir() bool
	GetLink() bool
	GetResource() bool
	GetCorrupted() bool
	GetEncrypted() bool
	GetEid() uint32
	GetEncoding() string
	GetRenaming() string
	GetError() error
	GetSize() uint
	SetRenaming(string)
	Destroy() error
}

type entry struct {
	filename   string
	dirP       int
	linkP      int
	resourceP  int
	corruptedP int
	encryptedP int
	eid        uint32
	encoding   string
	renaming   string
	err        *C.EntryError
	size       uint
	entryC     *C.struct_Entry
}

// Getters
func (ego *entry) GetFilename() string {
	return ego.filename
}

func (ego *entry) GetDir() bool {
	return ego.dirP != 0
}

func (ego *entry) GetLink() bool {
	return ego.linkP != 0
}

func (ego *entry) GetResource() bool {
	return ego.resourceP != 0
}

func (ego *entry) GetCorrupted() bool {
	return ego.corruptedP != 0
}

func (ego *entry) GetEncrypted() bool {
	return ego.encryptedP != 0
}

func (ego *entry) GetEid() uint32 {
	return ego.eid
}

func (ego *entry) GetEncoding() string {
	return ego.encoding
}

func (ego *entry) GetRenaming() string {
	return ego.renaming
}

func (ego *entry) GetError() error {

	ego.err = ego.entryC.error

	if ego.err == nil {
		return nil
	}

	err := C.GoString(ego.err.error_str)
	return errors.New(err)
}

func (ego *entry) GetSize() uint {
	return ego.size
}

func (ego *entry) SetRenaming(renaming string) {
	C.EntrySetRenaming(ego.entryC, C.CString(renaming))
	ego.renaming = renaming
}

func (ego *entry) Destroy() error {
	if ego.entryC == nil {
		return fmt.Errorf("Entry has been already destroyed.")
	}

	C.EntryDestroy(ego.entryC)
	ego.entryC = nil
	return nil
}

func DestroyList(entries []Entry) error {
	for j := 0; j < len(entries); j++ {
		curr := entries[j]
		err := curr.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
}
