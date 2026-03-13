package cmd

import (
	"fsync/filesystem"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type RootCommandHandler struct {
	SourceDirectoryPath      string
	DestinationDirectoryPath string
	DeleteMissingFlag        bool

	Command     *cobra.Command
	fileManager filesystem.FileManagerInterface
}

type RootCommandHandlerInterface interface {
	Run(cmd *cobra.Command, args []string)
}

func NewRootCommandHandler() *RootCommandHandler {
	rch := &RootCommandHandler{fileManager: filesystem.NewFileManager()}
	rch.Command = &cobra.Command{
		Use:     "fsync [src filepath] [dst filepath]",
		Args:    cobra.ExactArgs(2),
		PreRunE: rch.PreRunE,
		Run:     rch.Run,
	}
	rch.Command.PersistentFlags().BoolVarP(&rch.DeleteMissingFlag, "delete-missing", "", false, "--delete-missing")
	return rch
}

func (rch *RootCommandHandler) PreRunE(cmd *cobra.Command, args []string) error {
	log.Debug().Str("source directory", args[0]).Str("destination directory", args[1]).Send()
	src, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	dst, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	rch.SourceDirectoryPath, rch.DestinationDirectoryPath = src, dst
	return nil
}

func (rch *RootCommandHandler) Run(cmd *cobra.Command, args []string) {
	log.Info().Str("source directory", rch.SourceDirectoryPath).Msg("starting reading files from source")

	mainDirectory, err := rch.fileManager.ReadDirectory(rch.SourceDirectoryPath)
	if err != nil {
		log.Err(err).Msg("error while reading source directory")
		return
	}

	rch.verifyFiles(mainDirectory)
	rch.verifyDirectories(mainDirectory)
}

func (rch *RootCommandHandler) verifyFiles(parentDirectory *filesystem.Directory) {
	for _, f := range parentDirectory.Files {
		log.Info().Str("file", f.RelativePath).Msg("updating file")
		if rch.DeleteMissingFlag {
			rch.deleteMissingFiles(parentDirectory)
		}
		if err := f.Update(rch.fileManager, rch.SourceDirectoryPath, rch.DestinationDirectoryPath, f.ModTime, f.Size); err != nil {
			log.Err(err).Str("file", f.RelativePath).Msg("error while updating directory")
		}
	}
}

func (rch *RootCommandHandler) deleteMissingFiles(directory *filesystem.Directory) {
	directoryFullpath := filepath.Join(rch.DestinationDirectoryPath, directory.RelativePath)
	log.Debug().Str("directory fullpath", directoryFullpath).Str("directory", directory.RelativePath).Msg("starting deleting missing files")

	files, err := rch.fileManager.ReadDir(directoryFullpath)
	if err != nil {
		log.Err(err).Str("directory", directoryFullpath).Msg("error while reading directory - deleting missing files")
		return
	}

	for _, file := range files {
		fullpath := filepath.Join(directoryFullpath, file.Name())
		relativePath := filepath.Join(strings.TrimPrefix(directory.RelativePath, rch.DestinationDirectoryPath), file.Name())
		exist := false

		if file.IsDir() {
			_, exist = directory.Directories[relativePath]
		} else {
			_, exist = directory.Files[relativePath]
		}

		if !exist {
			if err = rch.fileManager.Remove(fullpath); err != nil {
				log.Err(err).Str("file", relativePath).Msg("error while removing file")
				continue
			}

			log.Info().Str("file", relativePath).Msg("successfully deleted file")
		}
	}
}

func (rch *RootCommandHandler) verifyDirectories(directory *filesystem.Directory) {
	for _, d := range directory.Directories {
		log.Info().Str("directory", d.RelativePath).Msg("updating directory")
		if err := d.Update(rch.fileManager, rch.DestinationDirectoryPath); err != nil {
			log.Err(err).Str("directory", d.RelativePath).Msg("error while updating directory")
		}

		rch.verifyDirectories(d)
		rch.verifyFiles(d)
	}
}
