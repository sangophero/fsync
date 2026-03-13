package filesystem

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type FileManagerInterface interface {
	ReadDirectory(dirPath string) (*Directory, error)

	Mkdir(name string) error
	ReadFile(name string) ([]byte, error)
	ReadDir(name string) ([]fs.DirEntry, error)
	WriteFile(name string, data []byte) error
	Remove(name string) error
}

type FileManager struct{}

func NewFileManager() FileManagerInterface {
	return new(FileManager)
}

func (f FileManager) Mkdir(name string) error {
	return os.Mkdir(name, 0777)
}

func (f FileManager) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (f FileManager) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (f FileManager) WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0777)
}

func (f FileManager) Remove(name string) error {
	return os.Remove(name)
}

func (f *FileManager) ReadDirectory(dirPath string) (*Directory, error) {
	mainDirectory := new(Directory)

	err := readDirectory(f, dirPath, dirPath, mainDirectory)
	return mainDirectory, err
}

type Directory struct {
	RelativePath string
	Directories  map[string]*Directory
	Files        map[string]*File
}

func (d *Directory) Update(fileManager FileManagerInterface, dstDirPath string) error {
	dstFilePath := filepath.Join(dstDirPath, d.RelativePath)
	log.Info().Str("directory", dstFilePath).Msg("creating directory")
	if err := fileManager.Mkdir(dstFilePath); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	return nil
}

type File struct {
	RelativePath string
	ModTime      time.Time
	Size         int64
	Hash         string
}

func (f *File) Update(fileManager FileManagerInterface, srcDirpath, dstDirPath string, dstFileModTime time.Time, dstFileSize int64) error {
	srcFilePath := filepath.Join(srcDirpath, f.RelativePath)
	dstFilePath := filepath.Join(dstDirPath, f.RelativePath)

	srcFileContent, err := fileManager.ReadFile(srcFilePath)
	if err != nil {
		return err
	}

	dstFileContent, err := fileManager.ReadFile(dstFilePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// If it doesn't exist -> create a new one
		log.Info().Str("file", dstFilePath).Msg("file doesn't exist, creating a new one")
		return fileManager.WriteFile(dstFilePath, srcFileContent)
	} else if err != nil {
		return err
	}

	if f.ModTime.Unix() != dstFileModTime.Unix() {
		log.Info().Str("file", dstFilePath).Msg("files mod times don't match, updating...")
		return fileManager.WriteFile(dstFilePath, srcFileContent)
	}

	if f.Size != dstFileSize {
		log.Info().Str("file", dstFilePath).Msg("files sizes don't match, updating...")
		return fileManager.WriteFile(dstFilePath, srcFileContent)
	}

	hashedDstFile, err := hashFile(dstFileContent)
	if err != nil {
		return err
	}

	if hashedDstFile == f.Hash {
		log.Debug().Str("file", dstFilePath).Msg("skipping updating file, nothing changed")
		return nil
	}
	log.Info().Str("file", dstFilePath).Msg("updating file")
	return fileManager.WriteFile(dstFilePath, srcFileContent)
}

func readDirectory(fileManager FileManagerInterface, rootDirPath, dirPath string, parentDirectory *Directory) error {
	parentDirectory.Directories = make(map[string]*Directory)
	parentDirectory.Files = make(map[string]*File)
	log.Debug().Str("directory", dirPath).Msg("reading directory")
	dir, err := fileManager.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range dir {
		if err != nil {
			return err
		}

		fullpath := filepath.Join(dirPath, file.Name())
		relativepath := filepath.Join(strings.TrimPrefix(dirPath, rootDirPath), file.Name())

		if file.IsDir() {
			dir := &Directory{
				RelativePath: relativepath,
			}

			parentDirectory.Directories[relativepath] = dir
			if err = readDirectory(fileManager, rootDirPath, fullpath, dir); err != nil {
				return err
			}
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			return err
		}

		fileContent, err := fileManager.ReadFile(fullpath)
		if err != nil {
			return err
		}

		hashedFile, err := hashFile(fileContent)
		if err != nil {
			return err
		}

		parentDirectory.Files[relativepath] = &File{
			RelativePath: relativepath,
			ModTime:      fileInfo.ModTime(),
			Size:         fileInfo.Size(),
			Hash:         hashedFile,
		}
	}
	return nil
}

func hashFile(fileContent []byte) (string, error) {
	h := sha256.New()

	if _, err := h.Write(fileContent); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}
