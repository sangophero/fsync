package filesystem

//filesystemMock "fsync/filesystem/mock"

/*
func TestReadDirectory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockedFileSystemInterface := NewMockFileSystemInterface(ctrl)

	fileManager := filesystemMock.NewMockFileManagerInterface(ctrl)

	fileManager.EXPECT().ReadDir("./src").Return([]fs.DirEntry{
		filesystemMock.DirEntryMock{
			CName:  "dir1",
			CIsDir: true,
		},
	}, nil)

	fileManager.EXPECT().ReadDir("src/dir1").Return([]fs.DirEntry{
		filesystemMock.DirEntryMock{
			CName:  "file1",
			CIsDir: false,
		},
		filesystemMock.DirEntryMock{
			CName:  "file2",
			CIsDir: false,
		},
	}, nil)

	fileManager.EXPECT().ReadFile("src/dir1/file1").Return([]byte("file1 content"), nil)
	fileManager.EXPECT().ReadFile("src/dir1/file2").Return([]byte("file2 content"), nil)

	directory, err := fileManager.ReadDirectory("./src")
	assert.Nil(t, err)

	expectedResponse := &Directory{
		RelativePath: "",
		Directories: map[string]*Directory{
			"dir1": {
				RelativePath: "dir1",
				Directories:  make(map[string]*Directory),
				Files: map[string]*File{
					"src/dir1/file1": {
						RelativePath: "src/dir1/file1",
						Hash:         "FwV4nTgO4RC8CSMd-K9CoMxWShUQ69IWhRbUmFxAomM=",
					},
					"src/dir1/file2": {
						RelativePath: "src/dir1/file2",
						Hash:         "adIb7zKdU6D2K_bxn6q5nqs9rE1T5D-7slf-5NFPYKs=",
					},
				},
			},
		},
		Files: make(map[string]*File),
	}
	assert.Equal(t, expectedResponse, directory)
}
*/
