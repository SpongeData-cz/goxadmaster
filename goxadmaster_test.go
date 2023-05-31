package goxadmaster_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	. "github.com/SpongeData-cz/goxadmaster"
)

func TestXADD(t *testing.T) {

	t.Run("example", func(t *testing.T) {

		archive := NewArchive("./fixtures/easy.zip")

		if archive.Err != nil {
			t.Error(archive.Err.Error())
		}

		pathToExtract := "./fixtures/extracted/easy/"

		archive.SetDestination(pathToExtract)
		archive.SetAlwaysOverwritesFiles(true)

		entries := archive.List()

		if archive.Err != nil {
			t.Error(archive.Err.Error())
		}

		// Optional entries rename
		for i := 0; i < len(entries); i++ {
			curr := entries[i]
			newName := fmt.Sprintf("binary%d.bin", i)
			curr.SetRenaming(pathToExtract + newName)
		}

		archive.Extract(entries)

		if archive.Err != nil {
			t.Error(archive.Err.Error())
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

		err = archive.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if archive.Err != nil {
			t.Error(archive.Err.Error())
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

		for i := 0; i < len(entries); i++ {
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

		DestroyList(entries)

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

		ar.SetBatch(-1, entries)

		for i := 0; i < len(entries); i++ {
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

		DestroyList(entries)

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

		ar.SetBatch(0, entries)

		for i := 0; i < len(entries); i++ {
			ar.Extract(entries)
			if ar.Err != nil {
				t.Error(ar.Err)
			}
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		DestroyList(entries)

		err = checkFiles(0, false)
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

		DestroyList(entries)

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
				println(entryErr)
				// fmt.Printf("WARNING: %s, WARNING MSG: %s\n", curr.GetFilename(), err.Error())
				// t.Error(entryErr)
			}
		}

		err := ar.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if ar.Err != nil {
			t.Error(ar.Err.Error())
		}

		DestroyList(entries)

		err = checkFiles(5, true)
		if err != nil {
			t.Error(err.Error())
		}
		err = removeExtracted()
		if err != nil {
			t.Error(err.Error())
		}

	})

	t.Run("getters", func(t *testing.T) {

		archive := NewArchive("./fixtures/easy.zip")

		if archive.Err != nil {
			t.Error(archive.Err.Error())
		}

		archive.SetDestination("./fixtures/extracted")
		archive.SetAlwaysOverwritesFiles(true)

		entries := archive.List()
		if archive.Err != nil {
			t.Error(archive.Err.Error())
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

		archive.Extract(entries)
		if archive.Err != nil {
			t.Error(archive.Err.Error())
		}

		err := DestroyList(entries)
		if err != nil {
			t.Error(err.Error())
		}

		err = archive.Destroy()
		if err != nil {
			t.Error(err.Error())
		}
		if archive.Err != nil {
			t.Error(archive.Err.Error())
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
