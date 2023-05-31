# goxadmaster
Golang binding for a [C wrapper](https://github.com/mafiosso/XADMaster) of a forked project of [XADMaster](https://github.com/MacPaw/XADMaster).

# Installation
## Requirements

[C wrapper](https://github.com/mafiosso/XADMaster#objective-c-library-for-archive-and-file-unarchiving-and-extraction) installation.

# Usage

```go
archive := NewArchive("./fixtures/easy.zip")

if archive.Err != nil {
    return archive.Err
}

pathToExtract := "./fixtures/extracted/"

// Destination setting is important for not renamed
// entries only. Otherwise ignored.
archive.SetDestination(pathToExtract)

archive.SetAlwaysOverwritesFiles(true)

// Make slice of Entries
entries := archive.List()

if archive.Err != nil {
    return archive.Err
}

// Optional entries rename
for i := 0; i < len(entries); i++ {
    curr := entries[i]
    newName := fmt.Sprintf("binary%d.bin", i)
    curr.SetRenaming(pathToExtract + newName)
}

// Does extraction over the list of Entries
// Note that you may pass just a subset of the original list using the SetBatch() function
archive.Extract(entries)

if archive.Err != nil {
    return archive.Err
}

// Inspection of per-Entry errors.
for i := 0; i < len(entries); i++ {
    curr := entries[i]
    err := curr.GetError()
    if err != nil {
        fmt.Printf("WARNING: %s, WARNING MSG: %s", curr.GetFilename(), err.Error())
    }
}

// Correct Entries removal
err := DestroyList(entries)
if err != nil {
    return err
}

// Correct Archive record deletion.
err = archive.Destroy()
if err != nil {
    return err
}

if archive.Err != nil {
    return archive.Err
}

```