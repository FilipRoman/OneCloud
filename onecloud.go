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

	fmt.Println(Onedrive)
	fmt.Println("By Virmox")

	var dstDir string
	var srcFile string
	var game string

	dirdata, errRead := os.ReadFile("Game_save_path.txt")
	oddata, errRead := os.ReadFile("OneDrive_save_path.txt")
	gamedata, errRead := os.ReadFile("Game_executable.txt")
	if errRead != nil {
		fmt.Println("Set your game's save folder path: \n(ex. C:/Users/admin/AppData/LocalLow/Subnautica/GameSaves)")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		Dpath := scanner.Text()
		errWrite := os.WriteFile("Game_save_path.txt", []byte(Dpath), 0666)
		if errWrite != nil {
			log.Fatal(errWrite)
		}
		fmt.Print("\033[H\033[2J")
		fmt.Println("Put your game's save folder inside OneDrive folder")
		time.Sleep(time.Second * 2)
		fmt.Println("Set your OneDrive game's save folder path: \n(ex. C:/Users/admin/OneDrive/Subnautica/GameSaves)")
		scanner.Scan()
		Odpath := scanner.Text()
		errWriteOD := os.WriteFile("OneDrive_save_path.txt", []byte(Odpath), 0666)
		if errWriteOD != nil {
			log.Fatal(errWriteOD)
		}
		fmt.Print("\033[H\033[2J")
		fmt.Println("Enter the name of the game's executable: \n(ex. Subnatica.exe)")
		scanner.Scan()
		GameExe := scanner.Text()
		errWriteG := os.WriteFile("Game_executable.txt", []byte(GameExe), 0666)
		if errWriteG != nil {
			log.Fatal(errWriteG)
		}
		dstDir = Dpath
		srcFile = Odpath
		game = GameExe
		fmt.Print("\033[H\033[2J")
		fmt.Println("Path set, exiting - reopen the app")
	} else {
		GameExe := string(gamedata)
		Dpath := string(dirdata)
		Odpath := string(oddata)
		fmt.Println("[1] - Set new game's save folder path")
		fmt.Println("[2] - Set new OneDrive save folder path")
		fmt.Println("[3] - Set new game's executable name")
		fmt.Println("[4] - Sync saves and lunch")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		option := scanner.Text()

		switch {
		case option == "1":
			fmt.Print("\033[H\033[2J")
			fmt.Println("Enter your new game's save folder path:\n(ex. C:/Users/admin/AppData/LocalLow/Subnatica/GameSaves)")
			scanner.Scan()
			dstDir = scanner.Text()
			errWriteOp1 := os.WriteFile("Game_save_path.txt", []byte(dstDir), 0666)
			if errWriteOp1 != nil {
				log.Fatal(errWriteOp1)
			}
		case option == "2":
			fmt.Print("\033[H\033[2J")
			fmt.Println("Enter your new OneDrive save folder path:\n(ex. C:/Users/admin/OneDrive/Subnautica/GameSaves)")
			scanner.Scan()
			dstDir = scanner.Text()
			errWriteOp2 := os.WriteFile("OneDrive_save_path.txt", []byte(dstDir), 0666)
			if errWriteOp2 != nil {
				log.Println(errWriteOp2)

			}
		case option == "3":
			fmt.Print("\033[H\033[2J")
			fmt.Println("Enter your new game's executable name:\n(ex. Subnatica.exe)")
			scanner.Scan()
			GameExe = scanner.Text()
			errWriteOp3 := os.WriteFile("Game_executable.txt", []byte(GameExe), 0666)
			if errWriteOp3 != nil {
				log.Println(errWriteOp3)

			}
		case option == "4":
			fmt.Print("\033[H\033[2J")
			dstDir = Dpath
			srcFile = Odpath
			game = GameExe
			dstFile := dstDir
			fmt.Println("Syncing game saves...")
			errcopy := copyFile(srcFile, dstFile)
			if errcopy != nil {
				log.Println(errcopy)

			}
			fmt.Println("Lunching game...")
			exePath, err := filepath.Abs(game)
			if err != nil {
				log.Println(err)

			}
			cmd := exec.Command(exePath)
			fmt.Println("Waiting for game to close...")
			err = cmd.Run()
			if err != nil {
				log.Println(err)

			}

			fmt.Println("Uploading saves to OneDrive...")
			time.Sleep(time.Second * 3)
			errcopy = copyFile(dstFile, srcFile)
			if errcopy != nil {
				log.Println(errcopy)
			}
			fmt.Println("Done, closing...")
			time.Sleep(time.Second * 1)
			os.Exit(0)

		default:
			log.Println("Invalid option. Restart and type a number")

		}
	}

	fmt.Println("Closing in 3 seconds")
	time.Sleep(time.Second * 3)
	//bye bye
	os.Exit(0)
}
