package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

func ReadMBR(path string) ([]byte, error) {
	b1 := make([]byte, 512)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	n, err := f.Read(b1)
	if err != nil {
		return nil, err
	}
	if n != 512 {
		return nil, fmt.Errorf("read %d bytes, expected 512", n)
	}
	f.Close()
	return b1, nil
}

func WriteMBR(mbr []byte, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	n, err := f.Write(mbr)
	if err != nil {
		return err
	}
	if n != 512 {
		return fmt.Errorf("wrote %d bytes, expected 512", n)
	}
	f.Close()
	return nil
}

func DumpMBR(mbr []byte) error {
	f, err := os.OpenFile("/cleanmbr.bin", os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	n, err := f.Write(mbr)
	if err != nil {
		return err
	}
	if n != 512 {
		return fmt.Errorf("wrote %d bytes, expected 512", n)
	}
	return nil
}

func main() {

	disclaimerReader := bufio.NewReader(os.Stdin)
	fmt.Printf("[!] This program will read directly from and write directly to the hard disk\n")
	fmt.Printf("[!] Master Boot Record (MBR)\n\n")
	fmt.Printf("[!] If it crashes, or behaves badly, or was tampered with, \n[!] your system will ALMOST CERTAINLY break\n\n")
	fmt.Printf("[!] It also might just save your bacon ;)\n\n")
	fmt.Printf("[!] If you understand these risks and wish to continue, \n[!] type \"yes\" at the prompt\n\n")
	fmt.Print("[?] Continue: ")
	response, _ := disclaimerReader.ReadString('\n')
	if response != "yes\n" {
		fmt.Printf("[!] Exiting...\n")
		os.Exit(0)
	}

	files, err := ioutil.ReadDir("/dev")
	if err != nil {
		panic(err)
	}
	r, err := regexp.Compile("[sh]d[a-z]$")
	if err != nil {
		panic(err)
	}
	disks := make([]string, 0)
	for _, f := range files {
		name := f.Name()
		if r.MatchString(name) {
			fmt.Printf("[+] Found likely disk: /dev/%s\n", name)
			disks = append(disks, "/dev/"+name)
		}
	}

	var bootDisk string
	if len(disks) == 0 {
		fmt.Printf("[x] Could not automatically detect hard disks.\n")
		fmt.Printf("[x] Cannot continue\n")
		os.Exit(1)
	}

	if len(disks) > 1 {
		var selection int
		fmt.Printf("[!] Multiple things that might be disks were found.\n")
		fmt.Printf("[!] Which one of these contains your MBR? (probably ends with \"a\")\n")
		for i, name := range disks {
			if i == 0 {
				fmt.Printf("%d	%s (default)\n", i+1, name)
			} else {
				fmt.Printf("%d	%s\n", i+1, name)
			}

		}
		for {
			fmt.Printf("[?] Enter 1-%d:", len(disks))
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if input == "\n" {
				selection = 1
				break
			} else {
				selection, err = strconv.Atoi(string(input[0]))
				if err != nil {
					fmt.Printf("error occurred during input processing: %s\n", err.Error())
					fmt.Printf("make sure you're entering a digit!\n")
					continue
				}
				if selection > len(disks) || selection < 1 {
					fmt.Printf("invalid selection number\n")
					continue
				}
				break
			}
		}
		bootDisk = disks[selection-1]
	} else {
		bootDisk = disks[0]
	}

	fmt.Printf("[+] Proceeding with boot disk %s...\n", bootDisk)

	cleanMBR, err := ReadMBR(bootDisk)
	if err != nil {
		fmt.Printf("[x] Failed to read boot disk mbr with error %s\n", err.Error())
		fmt.Printf("[x] Cannot continue\n")
		os.Exit(1)
	}
	sleepDuration, _ := time.ParseDuration("500ms")
	fmt.Printf("[+] MBR read into memory, beginning write loop...\n")
	for {
		time.Sleep(sleepDuration)
		err := WriteMBR(cleanMBR, bootDisk)
		if err != nil {
			fmt.Printf("[!] Write failed with error %s\n", err.Error())
			fmt.Printf("[!] Dumping clean MBR at /cleanmbr.bin and exiting...\n")
			err = DumpMBR(cleanMBR)
			if err != nil {
				fmt.Printf("[!] Dump failed with error %s\n", err.Error())
			}
			os.Exit(1)
		}
		fmt.Printf("[+] nyann't\n")
	}
}
