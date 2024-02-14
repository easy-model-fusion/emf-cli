package model

import (
	"fmt"
	"testing"

	"github.com/easy-model-fusion/client/test"
)

// getModel initiates a basic model with an id as suffix
func getModel(suffix int) Model {
	idStr := fmt.Sprint(suffix)
	return Model{
		Name:          "model" + idStr,
		Config:        Config{ModuleName: "module" + idStr, ClassName: "class" + idStr},
		DirectoryPath: "/path/to/model" + idStr,
		AddToBinary:   true,
	}
}

// TestEmpty_True tests the Empty function with an empty models slice.
func TestEmpty_True(t *testing.T) {
	// Init
	var models []Model

	// Execute
	isEmpty := Empty(models)

	// Assert
	test.AssertEqual(t, isEmpty, true, "Expected true.")
}

// TestEmpty_False tests the Empty function with a non-empty models slice.
func TestEmpty_False(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	isEmpty := Empty(models)

	// Assert
	test.AssertEqual(t, isEmpty, false, "Expected false.")
}

// TestContains_True tests the Contains function with an element contained by the slice.
func TestContains_True(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	contains := Contains(models, models[0])

	// Assert
	test.AssertEqual(t, contains, true, "Expected true.")
}

// TestContains_False tests the Contains function with an element not contained by the slice.
func TestContains_False(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	contains := Contains(models, getModel(2))

	// Assert
	test.AssertEqual(t, contains, false, "Expected false.")
}

// TestContainsByName_True tests the ContainsByName function with an element's name contained by the slice.
func TestContainsByName_True(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	contains := ContainsByName(models, models[0].Name)

	// Assert
	test.AssertEqual(t, contains, true, "Expected true.")
}

// TestContainsByName_False tests the ContainsByName function with an element's name not contained by the slice.
func TestContainsByName_False(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	contains := ContainsByName(models, getModel(2).Name)

	// Assert
	test.AssertEqual(t, contains, false, "Expected false.")
}

// TestDifference tests the Difference function to return the correct difference.
func TestDifference(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1), getModel(2), getModel(3), getModel(4)}
	index := 2
	sub := models[:index]
	expected := models[index:]

	// Execute
	difference := Difference(models, sub)

	// Assert
	test.AssertEqual(t, len(expected), len(difference), "Lengths should be equal.")
}

// TestUnion tests the Union function to return the correct union.
func TestUnion(t *testing.T) {
	// Init
	index := 2
	models1 := []Model{getModel(0), getModel(1), getModel(2), getModel(3), getModel(4)}
	models2 := models1[:index]
	expected := models2

	// Execute
	union := Union(models1, models2)

	// Assert
	test.AssertEqual(t, len(expected), len(union), "Lengths should be equal.")
}

// TestGetModelsByNames tests the GetModelsByNames function to return the correct models.
func TestGetModelsByNames(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}
	names := []string{models[0].Name, models[1].Name}

	// Execute
	result := GetModelsByNames(models, names)

	// Assert
	test.AssertEqual(t, len(models), len(result), "Lengths should be equal.")
}

// TestGetNames tests the GetNames function to return the correct names.
func TestGetNames(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	names := GetNames(models)

	// Assert
	test.AssertEqual(t, len(models), len(names), "Lengths should be equal.")
}
