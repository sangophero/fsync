package cmd

import (
	"fsync/filesystem"
	filesystemMock "fsync/filesystem/mock"
	"io/fs"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
)

func TestReadingDirectory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileManager := filesystemMock.NewMockFileManagerInterface(ctrl)
	rootCommandHandler := &RootCommandHandler{
		SourceDirectoryPath:      "/src",
		DestinationDirectoryPath: "/dst",
		DeleteMissingFlag:        true,
		fileManager:              fileManager,
	}

	fileManager.EXPECT().ReadDirectory("/src").Return(&filesystem.Directory{
		RelativePath: "",
		Directories: map[string]*filesystem.Directory{
			"dir1": {
				RelativePath: "dir1",
				Files: map[string]*filesystem.File{
					"dir1/file1": {
						RelativePath: "dir1/file1",
						Hash:         "adIb7zKdU6D2K_bxn6q5nqs9rE1T5D-7slf-5NFPYKs=",
					},
				},
			},
		},
		Files: make(map[string]*filesystem.File),
	}, nil)

	fileManager.EXPECT().ReadDir("/dst/dir1").Return([]fs.DirEntry{
		filesystemMock.DirEntryMock{
			CName:  "file1",
			CIsDir: false,
		},
		filesystemMock.DirEntryMock{
			CName:  "file2",
			CIsDir: false,
		},
	}, nil)

	// one dir should be created
	fileManager.EXPECT().Mkdir("/dst/dir1").Return(nil)

	// file2 should be removed
	fileManager.EXPECT().Remove("/dst/dir1/file2").Return(nil)

	fileManager.EXPECT().ReadFile("/dst/dir1/file1").Return([]byte("file1 content"), nil)
	fileManager.EXPECT().ReadFile("/src/dir1/file1").Return([]byte("file1 content"), nil)

	fileManager.EXPECT().WriteFile("/dst/dir1/file1", []byte("file1 content")).Return(nil)

	rootCommandHandler.Run(&cobra.Command{}, []string{})
}
