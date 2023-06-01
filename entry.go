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
	/*
		Get full path with filename within the archive.
	*/
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
	/*
		Destroys individual entry.

		Returns:
		  - error if entry has been already destroyed.
	*/
	Destroy() error
	updateError()
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

func (ego *entry) GetFilename() string {
	return ego.filename
}

/*
Predicate - is directory?
*/
func (ego *entry) GetDir() bool {
	return ego.dirP != 0
}

/*
Predicate - is link?
*/
func (ego *entry) GetLink() bool {
	return ego.linkP != 0
}

/*
Predicate - is a resource?
*/
func (ego *entry) GetResource() bool {
	return ego.resourceP != 0
}

/*
Predicate - is corrupted?
*/
func (ego *entry) GetCorrupted() bool {
	return ego.corruptedP != 0
}

/*
Predicate - is encrypted by using of password?
*/
func (ego *entry) GetEncrypted() bool {
	return ego.encryptedP != 0
}

/*
Get Entry unique identifier.
*/
func (ego *entry) GetEid() uint32 {
	return ego.eid
}

/*
Get Entry detected encoding.
*/
func (ego *entry) GetEncoding() string {
	return ego.encoding
}

/*
Get Entry renaming.

You may set entry destination by hand by setting this.
*/
func (ego *entry) GetRenaming() string {
	return ego.renaming
}

func (ego *entry) updateError() {
	ego.err = ego.entryC.error
}

/*
Get error record.
*/
func (ego *entry) GetError() error {

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
func (ego *entry) GetSize() uint {
	return ego.size
}

/*
Sets renaming for the entry from constant string and allocates copy.

Parameters:
  - renaming - path with a new name.
*/
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

/*
Destroys a slice of Entries.

Parameters:
  - entries - slice of Entries to be destroyed.

Returns:
  - error if any of the entry has already been destroyed.
*/
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
