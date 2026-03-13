package mock_filesystem

import fs "io/fs"

type DirEntryMock struct {
	CName  string
	CIsDir bool
}

func (m DirEntryMock) Name() string {
	return m.CName
}

func (m DirEntryMock) IsDir() bool {
	return m.CIsDir
}

func (m DirEntryMock) Type() fs.FileMode {
	return 0
}

func (m DirEntryMock) Info() (fs.FileInfo, error) {
	return nil, nil
}
