# goxadmaster
Golang binding for a [C wrapper](https://github.com/mafiosso/XADMaster) of a forked project of [XADMaster](https://github.com/MacPaw/XADMaster).

# Installation
## Requirements

[C wrapper](https://github.com/mafiosso/XADMaster#objective-c-library-for-archive-and-file-unarchiving-and-extraction) installation.

# Usage
Using goxadmaster is very similar to using the [C wrapper](https://github.com/mafiosso/XADMaster).

### Creating an Archive
Has to be deallocated with *Destroy* method after use.

```go
archive := goxadmaster.NewArchive("./fixtures/easy.zip")
if archive.Err != nil {
    return archive.Err
}
```
### Archive Set methods
After creating the Archive, it is possible to use the following *Set* methods:
```go
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
```

### List
The next step is the *List* method.
This method lists content of an archive in form of arrays.
Entry records must be destroyed by *DestroyList* call explicitly.
Alternatively, it is possible to destroy individual entries using the *Destroy* function.
```go
entries := archive.List()
if archive.Err != nil {
    return archive.Err
}
```
### Entry set method
Alternatively, Entries can be renamed using the *SetRenaming* method. The full path with the new name must be passed as a parameter.

If a Destination was set before the *SetRenaming*, it is ignored.
```go
for i := 0; i < len(entries); i++ {
    curr := entries[i]
    newName := fmt.Sprintf("binary%d.bin", i)
    curr.SetRenaming(pathToExtract + newName)
}
```

### Batch
After listing there is an option to use method *SetBatch*, which specifies how many Entries to extract.

If **batch** <= -1, everything will be extracted at once.
```go
archive.SetBatch(2, entries)
```

### Extraction
Here comes the time for extraction. This is done using the *Extract* function.
If **batch** > -1, it is necessary to perform the extraction repeatedly.
```go
archive.Extract(entries)
if archive.Err != nil {
    return archive.Err
}
```

### Errors
Next, it is a good idea to iterate through the items to see if any of them have errors. 

This can be done, for example, with the following code:
```go
for i := 0; i < len(entries); i++ {
    curr := entries[i]
    err := curr.GetError()
    if err != nil {
        fmt.Printf("WARNING: %s, WARNING MSG: %s", curr.GetFilename(), err.Error())
    }
}
```

### Destroy
Firstly, Entry records must be destroyed by *DestroyList* call explicitly or alternatively, it is possible to destroy individual entries using the *Destroy* function.
```go
err := goxadmaster.DestroyList(entries)
if err != nil {
    return err
}
```

Finally, the Archive is destroyed using the Destroy method.
```go
err = archive.Destroy()
if err != nil {
    return err
}
if archive.Err != nil {
    return archive.Err
}
```

# Example
```go
import (
	"fmt"

	"github.com/SpongeData-cz/goxadmaster"
)

func example() error {
    // Creates a new Archive.
    archive := goxadmaster.NewArchive("./fixtures/easy.zip")
    if archive.Err != nil {
        return archive.Err
    }

    // Destination setting is important for not renamed
    // entries only. Otherwise ignored.
    pathToExtract := "./fixtures/extracted/"    
    archive.SetDestination(pathToExtract)

    // Programmer may call archive Set methods here.
    archive.SetAlwaysOverwritesFiles(true)

    // Make slice of Entries.
    entries := archive.List()
    if archive.Err != nil {
        return archive.Err
    }

    // Optional entries rename.
    for i := 0; i < len(entries); i++ {
        curr := entries[i]
        newName := fmt.Sprintf("binary%d.bin", i)
        curr.SetRenaming(pathToExtract + newName)
    }

    // Does extraction over the list of Entries.
    // Note that you may pass just a subset of 
    // the original list using the SetBatch() function.
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

    // Correct Entries removal.
    err := goxadmaster.DestroyList(entries)
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

    return nil
}

```