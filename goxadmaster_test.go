package goxadmaster_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	. "github.com/SpongeData-cz/goxadmaster"
)

func TestXAD(t *testing.T) {

	t.Run("example", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		pathToExtract := "./fixtures/extracted/easy/"
		ar.SetDestination(pathToExtract)
		ar.SetAlwaysOverwritesFiles(true)

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			newName := fmt.Sprintf("binary%d.bin", i)
			curr.SetRenaming(pathToExtract + newName)
		}

		ar.Extract(entries)
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			err := curr.GetError()
			if err != nil {
				fmt.Printf("WARNING: %s, WARNING MSG: %s", curr.GetFilename(), err.Error())
			}
		}

		err := DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = checkFiles(5, true)
		if err != nil {
			t.Error(err.Error())
		}
		err = removeExtracted()
		if err != nil {
			t.Error(err.Error())
		}

	})

	t.Run("rename", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.SetAlwaysOverwritesFiles(true)

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			newName := fmt.Sprintf("./fixtures/extracted/easy/binary%d.bin", i)
			curr.SetRenaming(newName)
		}

		ar.Extract(entries)
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = checkFiles(5, true)
		if err != nil {
			t.Error(err.Error())
		}
		err = removeExtracted()
		if err != nil {
			t.Error(err.Error())
		}

	})

	t.Run("passwordArchive", func(t *testing.T) {

		ar := NewArchive("./fixtures/Passworded.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.SetDestination("./fixtures/extracted")
		ar.SetAlwaysOverwritesFiles(true)
		ar.SetPassword("1234")

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.Extract(entries)
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			err := curr.GetError()
			if err != nil {
				entryErr := fmt.Sprintf("WARNING: %s, WARNING MSG: %s\n", curr.GetFilename(), err.Error())
				println(entryErr) // or t.Error(entryErr)
			}
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = checkFiles(5, true)
		if err != nil {
			t.Error(err.Error())
		}
		err = removeExtracted()
		if err != nil {
			t.Error(err.Error())
		}

	})

	t.Run("gettersSetters", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		// For the test
		ar.SetDestination("./fixtures/extracted")
		ar.SetAlwaysOverwritesFiles(true)
		ar.SetEncodingName("")
		ar.SetPasswordEncodingName("")
		ar.SetAlwaysSkipsFiles(false)
		ar.SetExtractsSubArchives(false)
		ar.SetPropagatesRelevantMetadata(false)
		ar.SetCopiesArchiveModificationTimeToEnclosingDirectory(false)
		ar.SetMacResourceForkStyle(false)

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		curr := entries[0]

		if curr.GetFilename() != "script1.sh" {
			t.Error("File name doesn't match.")
		}
		if curr.GetDir() || curr.GetLink() || curr.GetResource() || curr.GetCorrupted() || curr.GetEncrypted() {
			t.Error("It is true in some case, though it should be false in all.")
		}
		if curr.GetEid() != 0 {
			t.Error("Eid is not 0.")
		}
		if curr.GetEncoding() != "" || curr.GetRenaming() != "" {
			t.Error("It has encoding or it has renaming, even if it's not supposed to.")
		}
		if curr.GetError() != nil {
			t.Error("An error is found when it shouldn't be.")
		}
		// Possibility of different size
		if curr.GetSize() != 8 {
			t.Error("Wrong size. (Possibility of different size)")
		}

		err := DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

	})

}

func TestBatch(t *testing.T) {

	t.Run("batch", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.SetDestination("./fixtures/extracted")
		ar.SetAlwaysOverwritesFiles(true)

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.SetBatch(2, entries)

		// len(entries) / batch
		iteration := (len(entries) / 2) + 1

		for i := 0; i < iteration; i++ {
			ar.Extract(entries)
			if ar.Err != nil {
				t.Error(ar.Err.Error())
			}
		}

		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			err := curr.GetError()
			if err != nil {
				fmt.Printf("WARNING: %s, WARNING MSG: %s", curr.GetFilename(), err.Error())
			}
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = checkFiles(5, true)
		if err != nil {
			t.Error(err.Error())
		}
		err = removeExtracted()
		if err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("negativeBatch", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.SetDestination("./fixtures/extracted")
		ar.SetAlwaysOverwritesFiles(true)

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		// Everything will be extracted at once.
		ar.SetBatch(-1, entries)

		ar.Extract(entries)
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			err := curr.GetError()
			if err != nil {
				fmt.Printf("WARNING: %s, WARNING MSG: %s", curr.GetFilename(), err.Error())
			}
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = checkFiles(5, true)
		if err != nil {
			t.Error(err.Error())
		}
		err = removeExtracted()
		if err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("emptyBatch", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		ar.SetDestination("./fixtures/extracted")
		ar.SetAlwaysOverwritesFiles(true)

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		// Nothing will be extracted.
		ar.SetBatch(0, entries)

		ar.Extract(entries)
		if ar.Err != nil {
			t.Error(ar.Err)
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = checkFiles(0, false)
		if err != nil {
			t.Error(err.Error())
		}

	})
}

func TestErr(t *testing.T) {

	t.Run("nonExistentArchive", func(t *testing.T) {

		ar := NewArchive("./fixtures/nonExist.zip")
		if ar.Err == nil {
			t.Error("Should throw an error.")
		}

		arErr := ar.Err

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != arErr {
			t.Error(ar.Err.Error())
		}
	})

	t.Run("DoubleDestroy", func(t *testing.T) {

		ar := NewArchive("./fixtures/easy.zip")
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		entries := ar.List()
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err := DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}
		err = DestroyList(entries)
		if err == nil {
			t.Error("Shoud throw an error, the entries have already been destroyed once.")
		}

		err = ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		err = ar.Destroy()
		if err == nil {
			t.Error("Shoud throw an error, the archive has already been destroyed once.")
		}
	})
}

func checkFiles(numOfFiles int, easy bool) error {

	path := "./fixtures/extracted"
	if easy {
		path += "/easy"
	}

	dir, err := os.Open(path)
	if err != nil {
		return err
	}

	files, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	if len(files) != numOfFiles {
		return errors.New("Wrong number of files extracted!")
	}
	return nil
}

func removeExtracted() error {

	path := "./fixtures/extracted/easy"

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}
