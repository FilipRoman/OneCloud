package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func copyFile(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			err = os.MkdirAll(dstPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			err = copyFileInternal(path, dstPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFileInternal(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func main() {
	Onedrive := `
  ____           _______             __
 / __ \___  ___ / ___/ /__  __ _____/ /
/ /_/ / _ \/ -_) /__/ / _ \/ // / _  / 
\____/_//_/\__/\___/_/\___/\_,_/\_,_/  
                                       
`
	Author := `
███ ▀▄    ▄         ▄   ▄█ █▄▄▄▄ █▀▄▀█ ████▄     ▄  
█  █  █  █           █  ██ █  ▄▀ █ █ █ █   █ ▀▄   █ 
█ ▀ ▄  ▀█       █     █ ██ █▀▀▌  █ ▄ █ █   █   █ ▀  
█  ▄▀  █         █    █ ▐█ █  █  █   █ ▀████  ▄ █   
███  ▄▀           █  █   ▐   █      █        █   ▀▄ 
                   █▐       ▀      ▀          ▀     
                   ▐                                
`
	fmt.Println(Onedrive)
	fmt.Println(Author)

	var dstDir string
	var srcFile string

	dirdata, errReadDir := os.ReadFile("Dredge_save_path.txt")
	oddata, errReadOd := os.ReadFile("OneDrive_save_path.txt")
	if errReadDir != nil {
		fmt.Println("Set your DREDGE folder path: \n(ex. C:/Users/admin/AppData/LocalLow/Black Salt Games/DREDGE)\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		Dpath := scanner.Text()
		errWrite := os.WriteFile("Dredge_save_path.txt", []byte(Dpath), 0666)
		if errWrite != nil {
			log.Fatal(errWrite)
		}
		fmt.Print("\033[H\033[2J")
		fmt.Println("Set your OneDrive folder path: \n(ex. C:/Users/admin/OneDrive/DREDGE)\n")
		scanner.Scan()
		Odpath := scanner.Text()
		errWriteOD := os.WriteFile("OneDrive_save_path.txt", []byte(Odpath), 0666)
		if errWriteOD != nil {
			log.Fatal(errWriteOD)
		}
		dstDir = Dpath
		srcFile = Odpath
	} else {

		Dpath := string(dirdata)
		Odpath := string(oddata)
		fmt.Println("[1] - Set new DREDGE folder path")
		fmt.Println("[2] - Set new OneDrive folder path")
		fmt.Println("[3] - Sync saves and lunch")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		option := scanner.Text()
		switch {
		case option == "1":
			fmt.Print("\033[H\033[2J")
			fmt.Println("Enter your new DREDGE folder path:\n(ex. C:/Users/admin/AppData/LocalLow/Black Salt Games/DREDGE)\n")
			scanner.Scan()
			dstDir = scanner.Text()
			errWriteOp1 := os.WriteFile("Dredge_save_path.txt", []byte(dstDir), 0666)
			if errWriteOp1 != nil {
				log.Fatal(errWriteOp1)
			}
		case option == "2":
			fmt.Print("\033[H\033[2J")
			fmt.Println("Enter your new OneDrive folder path:\n(ex. C:/Users/admin/OneDrive/DREDGE)\n")
			scanner.Scan()
			dstDir = scanner.Text()
			errWriteOp2 := os.WriteFile("OneDrive_save_path.txt", []byte(dstDir), 0666)
			if errWriteOp2 != nil {
				log.Println(errWriteOp2)

			}
		case option == "3":
			fmt.Print("\033[H\033[2J")
			dstDir = Dpath
			srcFile = Odpath
			dstFile := dstDir
			fmt.Println("Syncing DREDGE saves...")
			errcopy := copyFile(srcFile, dstFile)
			if errcopy != nil {
				log.Println(errcopy)

			}
			fmt.Println("Lunching DREDGE...")
			exePath, err := filepath.Abs("DREDGE.exe")
			if err != nil {
				log.Println(err)

			}

			cmd := exec.Command(exePath)
			err = cmd.Run()
			if err != nil {
				log.Println(err)

			}

			fmt.Println("Uploading saves to OneDrive...")
			errcopy = copyFile(dstFile, srcFile)
			if errcopy != nil {
				log.Println(errcopy)
			}

		default:
			log.Println("Invalid option")

		}
	}
	fmt.Print("\033[H\033[2J")
	fmt.Println("If first use, path set - open app again.\nIf not - closing in 3 seconds")
	time.Sleep(time.Second * 3)
	fmt.Println("Yes im going to leave that code working like that")
	time.Sleep(time.Second * 1)

	os.Exit(0)
	if errReadOd != nil {
		log.Fatal(errReadOd)
	}

}
