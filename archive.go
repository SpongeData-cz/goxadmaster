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

type Archive interface {
	List() []Entry
	Extract([]entry)
	SetBatch(int, []Entry)
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

type archive struct {
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
  - path - path to the existing archive
Returns:
  - pointer to a new instance of Archive.
*/
func NewArchive(path string) *archive {
	out := &archive{}
	out.path = path

	out.archive = C.ArchiveNew(C.CString(path))
	out.batch = -1

	out.updateArchiveError()

	return out
}

/*
Lists content of an archive in form of arrays.

Entry records must be destroyed by DestroyList() call explicitly.
Alternatively, it is possible to destroy individual entries using the Destroy() function.

Returns:
  - Slice of Entries.
*/
func (ego *archive) List() []Entry {

	ego.entries = C.ArchiveList(ego.archive)
	es := make([]Entry, 0)

	for elem := ego.entries; *elem != nil; elem = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(elem), POINTER_SIZE)) {
		es = append(es, &entry{
			filename:   C.GoString((*elem).filename),
			dirP:       int((*elem).dirP),
			linkP:      int((*elem).linkP),
			resourceP:  int((*elem).resourceP),
			corruptedP: int((*elem).corruptedP),
			encryptedP: int((*elem).encryptedP),
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
Creates a NULL terminated subset of the entries array.

Returns:
  - Newly created subse,
  - number of items in the subset (0 < number <= batch).
*/
func (ego *archive) makeSubSet() (**C.struct_Entry, int) {
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
Extract from the archive. If batch > -1, only batch entries are extracted.

Parameters:
  - entries - Slice of entries.
*/
func (ego *archive) Extract(entries []Entry) {
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
func (ego *archive) updateArchiveError() {
	if ego.archive.error_num != C.NO_ERROR {
		err := C.GoString(ego.archive.error_str)
		ego.Err = errors.New(err)
	}
}

/*
Updates errors in individual entries.

Parameters:
  - entries - Slice of entries for which errors are to be updated.
*/
func (ego *archive) updateEntriesErrors(entries []Entry) {
	for i := 0; i < len(entries); i++ {
		entries[i].updateError()
	}
}

/*
Sets the batch, which specifies how many Entries to extract.
If batch <= -1, everything will be extracted at once.

Parameters:
  - batch - number of Entries,
  - entries - The slice of Entries.
*/
func (ego *archive) SetBatch(batch int, entries []Entry) {
	if batch <= -1 || batch >= len(entries) {
		return
	}
	ego.batch = batch
	ego.setFrame(ego.entries, ego.batch)
}

/*
Updates the frame that shows what will subsequently be extracted.

Parameters:
  - start - Where to move the frame start pointer,
  - step - maximum batch size.
*/
func (ego *archive) setFrame(start **C.struct_Entry, step int) {
	ego.batchStart = start
	ego.batchEnd = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(ego.batchStart), (POINTER_SIZE * step)))
}

/*
Destroys the Archive.

Return:
  - error, if archive has been already destroyed, nil otherwise.
*/
func (ego *archive) Destroy() error {
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
Sets default destination for entries at extraction.
Destination setting is important for not renamed entries only. Otherwise ignored.

Parameters:
  - path - The destination path.
*/
func (ego *archive) SetDestination(path string) {
	C.ArchiveSetDestination(ego.archive, C.CString(path))
}

/*
Sets default password for entries at extraction.

Parameters:
  - password - The password.
*/
func (ego *archive) SetPassword(password string) {
	C.ArchiveSetPassword(ego.archive, C.CString(password))
}

/*
Sets default encoding name for entries at extraction.

Parameters:
  - encodingName - The encoding name.
*/
func (ego *archive) SetEncodingName(encodingName string) {
	C.ArchiveSetEncodingName(ego.archive, C.CString(encodingName))
}

/*
Sets default password encoding name for entries at extraction.

Parameters:
  - passEncodingName - The password encoding name.
*/
func (ego *archive) SetPasswordEncodingName(passEncodingName string) {
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
func (ego *archive) SetAlwaysOverwritesFiles(alwaysOverwriteFiles bool) {

	C.ArchiveSetAlwaysOverwritesFiles(ego.archive, goBoolToCInt(alwaysOverwriteFiles))
}

/*
Sets if always skip files on error.

Parameters:
  - alwaysSkipsFiles
*/
func (ego *archive) SetAlwaysSkipsFiles(alwaysSkipsFiles bool) {
	C.ArchiveSetAlwaysSkipsFiles(ego.archive, goBoolToCInt(alwaysSkipsFiles))
}

/*
Sets if extract also included subarchives. Not recommended set to yes - unsufficient testing.

Parameters:
  - extractsSubArchives
*/
func (ego *archive) SetExtractsSubArchives(extractsSubArchives bool) {
	C.ArchiveSetExtractsSubArchives(ego.archive, goBoolToCInt(extractsSubArchives))
}

/*
Sets if propagate relevant metadata (passwords etc.). Not recommended set to yes - unsufficient testing.

Parameters:
  - propagatesRelevantMetadata
*/
func (ego *archive) SetPropagatesRelevantMetadata(propagatesRelevantMetadata bool) {
	C.ArchiveSetPropagatesRelevantMetadata(ego.archive, goBoolToCInt(propagatesRelevantMetadata))
}

/*
Sets if to set entries' modification time also to the destination files.

Parameters:
  - copiesArchiveModificationTimeToEnclosingDirectory
*/
func (ego *archive) SetCopiesArchiveModificationTimeToEnclosingDirectory(copiesArchiveModificationTimeToEnclosingDirectory bool) {
	C.ArchiveSetCopiesArchiveModificationTimeToEnclosingDirectory(ego.archive, goBoolToCInt(copiesArchiveModificationTimeToEnclosingDirectory))
}

/*
Sets if to use MacOS forking style. Not recommended - not tested on Linux, just for completeness.

Parameters:
  - macResourceForkStyle
*/
func (ego *archive) SetMacResourceForkStyle(macResourceForkStyle bool) {
	C.ArchiveSetMacResourceForkStyle(ego.archive, goBoolToCInt(macResourceForkStyle))
}
