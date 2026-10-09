package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRuntimePath(t *testing.T) {
	ok := getRuntimePath()
	if !ok {
		t.Fatalf("expected getRuntimePath to succeed")
	}
	if runtimePath == "" {
		t.Fatalf("expected runtimePath to not be empty")
	}
}

func TestCartHeaderPointerOverlay(t *testing.T) {
	// Modifying cartridgeData directly should instantly reflect in cartHeader without any parsing
	copy(cartridgeData[0x134:0x134+4], []byte("TEST"))
	cartridgeData[0x147] = 0x01

	if string(cartHeader.Title[:4]) != "TEST" {
		t.Fatalf("expected title TEST, got %s", string(cartHeader.Title[:4]))
	}
	if cartHeader.CartridgeType != 0x01 {
		t.Fatalf("expected cartridge type 0x01, got 0x%02x", cartHeader.CartridgeType)
	}
}

func TestCartLoad(t *testing.T) {
	tetrisPath := filepath.Join("..", "tetris.gb")
	if _, err := os.Stat(tetrisPath); err == nil {
		if !CartLoad(tetrisPath) {
			t.Fatalf("CartLoad failed on %s", tetrisPath)
		}
		title := string(cartHeader.Title[:6])
		if title != "TETRIS" {
			t.Fatalf("expected title TETRIS, got %s", title)
		}
	}
}
