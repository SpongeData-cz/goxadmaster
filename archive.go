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

const STEP = bits.UintSize / 8

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
	SetAlwaysRenamesFiles(bool)
	SetAlwaysSkipsFiles(bool)
	SetExtractsSubArchives(bool)
	SetPropagatesRelevantMetadata(bool)
	SetCopiesArchiveModificationTimeToEnclosingDirectory(bool)
	SetMacResourceForkStyle(bool)
	SetPerIndexRenamedFiles(bool)
}

type archive struct {
	path       string
	Err        error
	batch      int // is valid 0
	batchStart **C.struct_Entry
	batchEnd   **C.struct_Entry
	entries    **C.struct_Entry
	archive    *C.struct_Archive
}

/*
Creates a new Archive. Has to be deallocated with Destroy method after use.
Parameters:
  - path - path to the existing archive
Returns: pointer to a new instance of Archive.
*/
func NewArchive(path string) *archive {
	out := &archive{}
	out.path = path

	out.archive = C.ArchiveNew(C.CString(path))
	out.batch = -1

	out.checkError()

	return out
}

/*
 */
func (ego *archive) List() []Entry {

	ego.entries = C.ArchiveList(ego.archive)
	es := make([]Entry, 0)

	for elem := ego.entries; *elem != nil; elem = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(elem), STEP)) {
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

	ego.checkError()

	return es
}

func (ego *archive) makeSubSet() (**C.struct_Entry, int) {
	subSet := (**C.struct_Entry)(C.calloc(C.ulong(ego.batch+1), C.ulong(unsafe.Sizeof(*ego.entries))))
	tmp := subSet
	i := 0
	for elem := ego.batchStart; *elem != nil && elem != ego.batchEnd; elem = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(elem), STEP)) {
		*subSet = *elem
		subSet = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(subSet), STEP))
		i++
	}
	return tmp, i
}

func (ego *archive) Extract(entries []Entry) {
	if ego.batch == 0 {
		return
	} else if ego.batch == -1 {
		C.ArchiveExtract(ego.archive, ego.entries)
	} else {
		subSet, read := ego.makeSubSet()

		ego.batchStart = ego.batchEnd
		ego.batchEnd = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(ego.batchEnd), (STEP * read)))
		C.ArchiveExtract(ego.archive, subSet)
		C.free(unsafe.Pointer(subSet))
	}

	ego.checkError()
	ego.fillErrors(entries)
}

func (ego *archive) checkError() {
	if ego.archive.error_num != C.NO_ERROR {
		err := C.GoString(ego.archive.error_str)
		ego.Err = errors.New(err)
	}
}

func (ego *archive) fillErrors(entries []Entry) {
	for i := 0; i < len(entries); i++ {
		entries[i].GetError()
	}
}

func (ego *archive) SetBatch(batch int, entries []Entry) {
	if batch == -1 || batch >= len(entries) {
		return
	}
	ego.batch = batch
	ego.batchStart = ego.entries
	ego.batchEnd = (**C.struct_Entry)(unsafe.Add(unsafe.Pointer(ego.batchStart), (STEP * ego.batch)))
}

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

func (ego *archive) SetDestination(path string) {
	C.ArchiveSetDestination(ego.archive, C.CString(path))
}

func (ego *archive) SetPassword(password string) {
	C.ArchiveSetPassword(ego.archive, C.CString(password))
}

func (ego *archive) SetEncodingName(encodingName string) {
	C.ArchiveSetEncodingName(ego.archive, C.CString(encodingName))
}

func (ego *archive) SetPasswordEncodingName(passEncodingName string) {
	C.ArchiveSetPasswordEncodingName(ego.archive, C.CString(passEncodingName))
}

func goBoolToCInt(param bool) C.int {
	if param {
		return C.int(1)
	}
	return C.int(0)
}

func (ego *archive) SetAlwaysOverwritesFiles(alwaysOverwriteFiles bool) {

	C.ArchiveSetAlwaysOverwritesFiles(ego.archive, goBoolToCInt(alwaysOverwriteFiles))
}

func (ego *archive) SetAlwaysRenamesFiles(alwaysRenamesFiles bool) {
	C.ArchiveSetAlwaysRenamesFiles(ego.archive, goBoolToCInt(alwaysRenamesFiles))
}

func (ego *archive) SetAlwaysSkipsFiles(alwaysSkipsFiles bool) {
	C.ArchiveSetAlwaysSkipsFiles(ego.archive, goBoolToCInt(alwaysSkipsFiles))
}

func (ego *archive) SetExtractsSubArchives(extractsSubArchives bool) {
	C.ArchiveSetExtractsSubArchives(ego.archive, goBoolToCInt(extractsSubArchives))
}

func (ego *archive) SetPropagatesRelevantMetadata(propagatesRelevantMetadata bool) {
	C.ArchiveSetPropagatesRelevantMetadata(ego.archive, goBoolToCInt(propagatesRelevantMetadata))
}

func (ego *archive) SetCopiesArchiveModificationTimeToEnclosingDirectory(copiesArchiveModificationTimeToEnclosingDirectory bool) {
	C.ArchiveSetCopiesArchiveModificationTimeToEnclosingDirectory(ego.archive, goBoolToCInt(copiesArchiveModificationTimeToEnclosingDirectory))
}

func (ego *archive) SetMacResourceForkStyle(macResourceForkStyle bool) {
	C.ArchiveSetMacResourceForkStyle(ego.archive, goBoolToCInt(macResourceForkStyle))
}

func (ego *archive) SetPerIndexRenamedFiles(perIndexRenamedFiles bool) {
	C.ArchiveSetPerIndexRenamedFiles(ego.archive, goBoolToCInt(perIndexRenamedFiles))
}
