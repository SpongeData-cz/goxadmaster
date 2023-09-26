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
)

type IEntry interface {
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
	updateError()
}

type Entry struct {
	filename   string
	dirP       bool
	linkP      bool
	resourceP  bool
	corruptedP bool
	encryptedP bool
	eid        uint32
	encoding   string
	renaming   string
	err        *C.EntryError
	size       uint
	entryC     *C.struct_Entry
}

/*
Get full path with filename within the archive.
*/
func (ego *Entry) GetFilename() string {
	return ego.filename
}

/*
Predicate - is directory?
*/
func (ego *Entry) GetDir() bool {
	return ego.dirP
}

/*
Predicate - is link?
*/
func (ego *Entry) GetLink() bool {
	return ego.linkP
}

/*
Predicate - is a resource?
*/
func (ego *Entry) GetResource() bool {
	return ego.resourceP
}

/*
Predicate - is corrupted?
*/
func (ego *Entry) GetCorrupted() bool {
	return ego.corruptedP
}

/*
Predicate - is encrypted by using of password?
*/
func (ego *Entry) GetEncrypted() bool {
	return ego.encryptedP
}

/*
Get Entry unique identifier.
*/
func (ego *Entry) GetEid() uint32 {
	return ego.eid
}

/*
Get Entry detected encoding.
*/
func (ego *Entry) GetEncoding() string {
	return ego.encoding
}

/*
Get Entry renaming.

You may set Entry destination by hand by setting this.
*/
func (ego *Entry) GetRenaming() string {
	return ego.renaming
}

func (ego *Entry) updateError() {
	ego.err = ego.entryC.error
}

/*
Get error record.
*/
func (ego *Entry) GetError() error {

	ego.updateError()

	if ego.err == nil {
		return nil
	}

	err := C.GoString(ego.err.error_str)
	return errors.New(err)
}

/*
Get unpacked size.
*/
func (ego *Entry) GetSize() uint {
	return ego.size
}

/*
Sets renaming for the Entry from constant string and allocates copy.

Parameters:
  - renaming - path with a new name.
*/
func (ego *Entry) SetRenaming(renaming string) {
	C.EntrySetRenaming(ego.entryC, C.CString(renaming))
	ego.renaming = renaming
}

/*
Destroys individual Entry.

Returns:
  - error if Entry has been already destroyed.
*/
func (ego *Entry) Destroy() error {
	if ego.entryC == nil {
		return errors.New("entry has been already destroyed")
	}

	C.EntryDestroy(ego.entryC)
	ego.entryC = nil
	return nil
}

/*
Destroys a slice of Entries.

Parameters:
  - entries - slice of Entries to be destroyed.

Returns:
  - error if any of the Entry has already been destroyed.
*/
func DestroyList(entries []IEntry) error {
	for j := 0; j < len(entries); j++ {
		curr := entries[j]
		err := curr.Destroy()
		if err != nil {
			return err
		}
	}
	return nil
}
