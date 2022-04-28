package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	api "github.com/noarainstorm/uploadgramApiGo"

	"github.com/spf13/cobra"
)

func main() {
	conf := api.New("", "")
	mainCmd := &cobra.Command{
		Use:   "./uploadcli [Upload | Download | Delete | Rename]",
		Short: "Very fast uploadgram client for cli",
	}

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload files to uploadgram",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Few arguments to run..")
			}
			upload(args, conf)
		},
	}

	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Download files from uploadgram",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Few arguments to run..")
			}
			download(args, conf)
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete files from uploadgram",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Few arguments to run..")
			}
			delete(args, conf)
		},
	}

	renameCmd := &cobra.Command{
		Use:   "rename",
		Short: "Reaneme files in uploadgram",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatal("Few arguments to run..")
			}
			rename(args[0], args[1], conf)
		},
	}

	mainCmd.AddCommand(uploadCmd, downloadCmd, deleteCmd, renameCmd)
	mainCmd.Execute()
}

func upload(fileName []string, ug api.All) {
	for cou, elem := range fileName {
		log.Println("Uploading file... ", cou+1)
		err := ug.Upload(elem)
		if err != nil {
			if errors.Is(err, api.ErrServer) {
				log.Fatal("An error occurred on the server!")
			}
			if errors.Is(err, api.ErrOpenFile) {
				log.Fatal("There is no such file on the disk!")
			}
			log.Fatal(err)
		}
		fmt.Printf("Done!\nLink: %s\nToken: %s\n", ug.Response.Url, ug.Response.Token)
	}
}

func download(urls []string, ug api.All) {
	log.Println("Start downloading files..")
	for cou, elem := range urls {
		data, err := ug.Download(elem)
		if err != nil {
			if errors.Is(err, api.ErrNotFound) {
				log.Fatal("The file was not found on the server!")
			}
			if errors.Is(err, api.ErrLink) {
				log.Fatal("Invalid link!")
			}
			log.Fatal(err)
		}
		file, err := os.OpenFile(ug.DownloadInfo.Filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Done! File %d downloaded...", cou+1)
	}
}

func delete(urls []string, ug api.All) {
	log.Println("Trying to delete the file...")
	for cou, elem := range urls {
		err := ug.Delete(elem)
		if err != nil {
			if errors.Is(err, api.ErrToken) {
				log.Fatal("Invalid token!")
			}
		}
		if ug.DeleteStat == 200 {
			log.Println("Successfully deleted ", cou+1)
		}
		if errors.Is(err, api.ErrNotFound) {
			log.Fatal("Something went wrong... Check the token")
		}
	}
}

func rename(token string, newName string, ug api.All) {
	log.Println("Trying to rename the file...")
	err := ug.Rename(token, newName)
	if errors.Is(err, api.ErrToken) {
		log.Fatal("Invalid token!")
	}
	if errors.Is(err, api.ErrFileName) {
		log.Fatal("The file name is too short!")
	}
	if errors.Is(err, api.ErrNotFound) {
		log.Fatal("Something went wrong... Check the token")
	}
	log.Printf("Successfully renamed to %s\n", newName)
}
