package ui

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestNewPTermUI(t *testing.T) {
	ui := NewPTermUI()
	test.AssertNotEqual(t, ui, nil, "UI should not be nil")
}

func TestPtermDefaultBoxPrinter_Print(t *testing.T) {
	ui := NewPTermUI()
	ui.DefaultBox().Print("test")
}

func TestPtermDefaultBoxPrinter_Println(t *testing.T) {
	ui := NewPTermUI()
	ui.DefaultBox().Println("test")
}

func TestPtermDefaultBoxPrinter_Printf(t *testing.T) {
	ui := NewPTermUI()
	ui.DefaultBox().Printf("test")
}

func TestPtermDefaultBoxPrinter_Printfln(t *testing.T) {
	ui := NewPTermUI()
	ui.DefaultBox().Printfln("test")
}

func TestPtermPrinter_Print(t *testing.T) {
	ui := NewPTermUI()
	ui.Error().Print("test")
	ui.Warning().Print("test")
	ui.Info().Print("test")
	ui.Success().Print("test\n")
}

func TestPtermPrinter_Println(t *testing.T) {
	ui := NewPTermUI()
	ui.Error().Println("test")
	ui.Warning().Println("test")
	ui.Info().Println("test")
	ui.Success().Println("test")
}

func TestPtermPrinter_Printf(t *testing.T) {
	ui := NewPTermUI()
	ui.Error().Printf("test")
	ui.Warning().Printf("test")
	ui.Info().Printf("test")
	ui.Success().Printf("test")
}

func TestPtermPrinter_Printfln(t *testing.T) {
	ui := NewPTermUI()
	ui.Error().Printfln("test")
	ui.Warning().Printfln("test")
	ui.Info().Printfln("test")
	ui.Success().Printfln("test")
}

func TestPtermUI_Error(t *testing.T) {
	ui := NewPTermUI()
	ui.Error().Print("test\n")
}

func TestPtermUI_Info(t *testing.T) {
	ui := NewPTermUI()
	ui.Info().Print("test\n")
}

func TestPtermUI_Success(t *testing.T) {
	ui := NewPTermUI()
	ui.Success().Print("test\n")
}

func TestPtermUI_Warning(t *testing.T) {
	ui := NewPTermUI()
	ui.Warning().Print("test\n")
}

func TestPtermUI_DisplaySelectedItems(t *testing.T) {
	ui := NewPTermUI()
	ui.DisplaySelectedItems([]string{"test1", "test2"})
}

func TestPtermUI_StartSpinner(t *testing.T) {
	ui := NewPTermUI()
	spinner := ui.StartSpinner("test")
	test.AssertNotEqual(t, spinner, nil, "Spinner should not be nil")
	spinner.Fail("test")
}
