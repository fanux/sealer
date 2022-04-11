package merge

import (
	"fmt"

	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
)

func MergeImages(images []string, target string) error {
	options, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	options.GraphDriverOptions = []string{}
	if err != nil {
		fmt.Println("err ", err)
		return err
	}
	store, err := storage.GetStore(options)
	if err != nil {
		fmt.Printf("err: %v", err)
		return err
	}

	imageList, err := store.Images()
	for _, i := range imageList {
		fmt.Printf("imagename: %v", i.Names)
	}

	metadata, err := store.Metadata("2fb6fc2d97e1")
	fmt.Println("metadata is:", metadata)

	dir, err := store.Mount("2fb6fc2d97e1", "")
	fmt.Println("mount dir is: ", dir, err)

	return nil
}
