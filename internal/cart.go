package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

const maxCartSize = 1024 * 1024 // 1MB

var (
	runtimePath     string
	cartridgeData   [maxCartSize]uint8
	cartridgeLoaded bool = false
	cartHeader           = (*CartHeaderStruct)(unsafe.Pointer(&cartridgeData[0x100]))
)

type CartHeaderStruct struct {
	EntryPoint           [4]uint8
	NintendoLogo         [48]uint8
	Title                [15]uint8
	CgbFlag              uint8
	NewLicenseeCode      [2]uint8
	SgbFlag              uint8
	CartridgeType        uint8
	RomSize              uint8
	RamSize              uint8
	DestinationCode      uint8
	OldLicenseeNumber    uint8
	MaskRomVersionNumber uint8
	HeaderChecksum       uint8
	GlobalChecksumHi     uint8
	GlobalChecksumLo     uint8
}

func getRuntimePath() bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}
	runtimePath = dir
	return true
}

func CartOpenFile() bool {
	cmd := exec.Command("zenity", "--file-selection",
		"--title=Select Game Boy ROM",
		"--file-filter=Game Boy ROMs (*.gb, *.gbc, *.bin) | *.gb *.gbc *.bin",
		"--file-filter=All Files | *",
	)

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	selectedPath := strings.TrimSpace(string(out))
	if selectedPath != "" {
		return CartLoad(selectedPath)
	}
	return false
}

func CartPrintInfo() {
	title := strings.TrimRight(string(cartHeader.Title[:]), "\x00")
	fmt.Printf("Entry point: %.2X%.2X%.2X%.2X\n", cartHeader.EntryPoint[0], cartHeader.EntryPoint[1], cartHeader.EntryPoint[2], cartHeader.EntryPoint[3])
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("CGB Flag: %02X\n", cartHeader.CgbFlag)
	fmt.Printf("SGB Flag: %02X\n", cartHeader.SgbFlag)
	fmt.Printf("Cartridge Type: %02X\n", cartHeader.CartridgeType)
	fmt.Printf("ROM Size: %02X\n", cartHeader.RomSize)
	fmt.Printf("RAM Size: %02X\n", cartHeader.RamSize)
	fmt.Printf("Destination Code: %02X\n", cartHeader.DestinationCode)
	fmt.Printf("Old Licensee Code: %02X\n", cartHeader.OldLicenseeNumber)
	fmt.Printf("Mask ROM Version Number: %02X\n", cartHeader.MaskRomVersionNumber)
	fmt.Printf("Header Checksum: %02X\n", cartHeader.HeaderChecksum)
	fmt.Printf("Global Checksum: %02X%02X\n", cartHeader.GlobalChecksumHi, cartHeader.GlobalChecksumLo)
}

func CartLoad(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open ROM: %v\n", err)
		return false
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Failed to stat ROM: %v\n", err)
		return false
	}

	fileSize := fileInfo.Size()
	if fileSize > maxCartSize {
		fmt.Println("File is too big")
		return false
	}

	// Read ROM into cartridgeData; cartHeader automatically maps to 0x100
	_, err = io.ReadFull(file, cartridgeData[:fileSize])
	if err != nil {
		fmt.Printf("Failed to read ROM: %v\n", err)
		return false
	}

	cartridgeLoaded = true
	fmt.Printf("ROM %s loaded, size: %d bytes\n", filename, fileSize)
	return true
}
