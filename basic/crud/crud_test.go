package crud

import (
	"fmt"
	"testing"
)

func TestCRUD(t *testing.T) {
	RunClientsetCRUD()
}

func RunClientsetCRUD() {
	fmt.Println("=== createPod ===")
	CreatePod()
	fmt.Println("=== listPod ===")
	ListPod()
	fmt.Println("=== getPod ===")
	GetPod()
	fmt.Println("=== deletePod ===")
	UpdatePod()
	fmt.Println("=== updatePod ===")
	DeletePod()
}
