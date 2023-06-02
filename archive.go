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
	"math/bits"
	"unsafe"
)

const POINTER_SIZE = bits.UintSize / 8

type IArchive interface {
	List() []IEntry
	Extract([]IEntry)
	SetBatch(int, []IEntry)
	Destroy() error
	SetDestination(string)
	SetPassword(string)
	SetEncodingName(string)
	SetPasswordEncodingName(string)
	SetAlwaysOverwritesFiles(bool)
	SetAlwaysSkipsFiles(bool)
	SetExtractsSubArchives(bool)
	SetPropagatesRelevantMetadata(bool)
	SetCopiesArchiveModificationTimeToEnclosingDirectory(bool)
	SetMacResourceForkStyle(bool)
}

type Archive struct {
	path       string
	Err        error
	batch      int
	batchStart **C.struct_Entry
	batchEnd   **C.struct_Entry
	entries    **C.struct_Entry
	archive    *C.struct_Archive
}

/*
Creates a new Archive. Has to be deallocated with Destroy method after use.

Parameters:
  - path - path to the existing Archive.

Returns:
  - pointer to a new instance of Archive.
*/
func NewArchive(path string) *Archive {
	out := &Archive{}
	out.path = path

	out.archive = C.ArchiveNew(C.CString(path))
	out.batch = -1

	out.updateArchiveError()

	return out
}

/*
Lists content of an Archive in form of arrays.

Entry records must be destroyed by DestroyList() call explicitly.
Alternatively, it is possible to destroy individual Entries using the Destroy() function.

Returns:
  - slice of Entries.
*/
func (ego *Archive) List() []IEntry {

	ego.entries = C.ArchiveList(ego.archive)
	es := make([]IEntry, 0)

	for elem := ego.entries; *elem != nil; elem = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(elem), POINTER_SIZE)) {
		es = append(es, &Entry{
			filename:   C.GoString((*elem).filename),
			dirP:       int((*elem).dirP) != 0,
			linkP:      int((*elem).linkP) != 0,
			resourceP:  int((*elem).resourceP) != 0,
			corruptedP: int((*elem).corruptedP) != 0,
			encryptedP: int((*elem).encryptedP) != 0,
			eid:        uint32((*elem).eid),
			encoding:   C.GoString((*elem).encoding),
			renaming:   C.GoString((*elem).renaming),
			err:        (*elem).error,
			size:       uint((*elem).size),
			entryC:     (*elem),
		})
	}

	ego.updateArchiveError()
	return es
}

/*
Creates a NULL-terminated subset of the Entries array.

Returns:
  - newly created subset,
  - number of items in the subset (0 < number <= batch).
*/
func (ego *Archive) makeSubSet() (**C.struct_Entry, int) {
	subSet := (**C.struct_Entry)(C.calloc(C.ulong(ego.batch+1), C.ulong(unsafe.Sizeof(*ego.entries))))
	tmp := subSet
	i := 0
	for elem := ego.batchStart; *elem != nil && elem != ego.batchEnd; elem = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(elem), POINTER_SIZE)) {
		*subSet = *elem
		subSet = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(subSet), POINTER_SIZE))
		i++
	}
	return tmp, i
}

/*
Extract from the Archive. If batch > -1, only batch Entries are extracted.

Parameters:
  - entries - slice of Entries.
*/
func (ego *Archive) Extract(entries []IEntry) {
	if ego.batch == 0 {
		return
	} else if ego.batch == -1 {
		C.ArchiveExtract(ego.archive, ego.entries)
	} else {
		subSet, read := ego.makeSubSet()

		if read != 0 {
			ego.setFrame(ego.batchEnd, read)
			C.ArchiveExtract(ego.archive, subSet)
		}
		C.free(unsafe.Pointer(subSet))
	}

	ego.updateArchiveError()
	ego.updateEntriesErrors(entries)
}

/*
Updates Archive error.
*/
func (ego *Archive) updateArchiveError() {
	if ego.archive.error_num != C.NO_ERROR {
		err := C.GoString(ego.archive.error_str)
		ego.Err = errors.New(err)
	}
}

/*
Updates errors in individual Entries.

Parameters:
  - entries - slice of Entries for which errors are to be updated.
*/
func (ego *Archive) updateEntriesErrors(entries []IEntry) {
	for i := 0; i < len(entries); i++ {
		entries[i].updateError()
	}
}

/*
Sets the batch, which specifies how many Entries to extract.
If batch <= -1, everything will be extracted at once.

Parameters:
  - batch - number of Entries,
  - entries - the slice of Entries.
*/
func (ego *Archive) SetBatch(batch int, entries []IEntry) {
	if batch <= -1 || batch >= len(entries) {
		return
	}
	ego.batch = batch
	ego.setFrame(ego.entries, ego.batch)
}

/*
Updates the frame that shows what will subsequently be extracted.

Parameters:
  - start - where to move the frame start pointer,
  - step - maximum batch size.
*/
func (ego *Archive) setFrame(start **C.struct_Entry, step int) {
	ego.batchStart = start
	ego.batchEnd = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(ego.batchStart), (POINTER_SIZE * step)))
}

/*
Destroys the Archive.

Returns:
  - error, if Archive has been already destroyed, nil otherwise.
*/
func (ego *Archive) Destroy() error {
	if ego.archive == nil {
		return fmt.Errorf("Archive has been already destroyed.")
	}

	C.ArchiveDestroy(ego.archive)
	C.free(unsafe.Pointer(ego.entries))

	ego.entries = nil
	ego.archive = nil
	ego.batchStart = nil
	ego.batchEnd = nil

	return nil
}

/*
Sets default destination for Entries at extraction.
Destination setting is important for not renamed Entries only. Otherwise ignored.

Parameters:
  - path - the destination path.
*/
func (ego *Archive) SetDestination(path string) {
	C.ArchiveSetDestination(ego.archive, C.CString(path))
}

/*
Sets default password for Entries at extraction.

Parameters:
  - password - the password.
*/
func (ego *Archive) SetPassword(password string) {
	C.ArchiveSetPassword(ego.archive, C.CString(password))
}

/*
Sets default encoding name for Entries at extraction.

Parameters:
  - encodingName - the encoding name.
*/
func (ego *Archive) SetEncodingName(encodingName string) {
	C.ArchiveSetEncodingName(ego.archive, C.CString(encodingName))
}

/*
Sets default password encoding name for Entries at extraction.

Parameters:
  - passEncodingName - the password encoding name.
*/
func (ego *Archive) SetPasswordEncodingName(passEncodingName string) {
	C.ArchiveSetPasswordEncodingName(ego.archive, C.CString(passEncodingName))
}

func goBoolToCInt(param bool) C.int {
	if param {
		return C.int(1)
	}
	return C.int(0)
}

/*
Sets if always overwrite files if they are present on the destination path.

Parameters:
  - alwaysOverwriteFiles
*/
func (ego *Archive) SetAlwaysOverwritesFiles(alwaysOverwriteFiles bool) {

	C.ArchiveSetAlwaysOverwritesFiles(ego.archive, goBoolToCInt(alwaysOverwriteFiles))
}

/*
Sets if always skip files on error.

Parameters:
  - alwaysSkipsFiles
*/
func (ego *Archive) SetAlwaysSkipsFiles(alwaysSkipsFiles bool) {
	C.ArchiveSetAlwaysSkipsFiles(ego.archive, goBoolToCInt(alwaysSkipsFiles))
}

/*
Sets if extract also included subarchives. Not recommended set to yes - unsufficient testing.

Parameters:
  - extractsSubArchives
*/
func (ego *Archive) SetExtractsSubArchives(extractsSubArchives bool) {
	C.ArchiveSetExtractsSubArchives(ego.archive, goBoolToCInt(extractsSubArchives))
}

/*
Sets if propagate relevant metadata (passwords etc.). Not recommended set to yes - unsufficient testing.

Parameters:
  - propagatesRelevantMetadata
*/
func (ego *Archive) SetPropagatesRelevantMetadata(propagatesRelevantMetadata bool) {
	C.ArchiveSetPropagatesRelevantMetadata(ego.archive, goBoolToCInt(propagatesRelevantMetadata))
}

/*
Sets if to set Entries modification time also to the destination files.

Parameters:
  - copiesArchiveModificationTimeToEnclosingDirectory
*/
func (ego *Archive) SetCopiesArchiveModificationTimeToEnclosingDirectory(copiesArchiveModificationTimeToEnclosingDirectory bool) {
	C.ArchiveSetCopiesArchiveModificationTimeToEnclosingDirectory(ego.archive, goBoolToCInt(copiesArchiveModificationTimeToEnclosingDirectory))
}

/*
Sets if to use MacOS forking style. Not recommended - not tested on Linux, just for completeness.

Parameters:
  - macResourceForkStyle
*/
func (ego *Archive) SetMacResourceForkStyle(macResourceForkStyle bool) {
	C.ArchiveSetMacResourceForkStyle(ego.archive, goBoolToCInt(macResourceForkStyle))
}
