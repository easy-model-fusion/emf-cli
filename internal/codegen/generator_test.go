package codegen

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

func TestPythonCodeGenerator_Generate(t *testing.T) {
	gen := NewPythonCodeGenerator(true)
	code, err := gen.Generate(&File{
		Name: "test.py",
		HeaderComments: []string{
			"Code generated by EMF",
			"DO NOT EDIT!",
		},
		Imports: []Import{
			{
				What: []ImportWhat{
					{
						Name: "os",
					},
				},
			},
			{
				What: []ImportWhat{
					{
						Name: "List",
					},
				},
				From: "typing",
			},
		},
		Functions: []*Function{
			{
				Name: "main",
				Params: []Parameter{
					{
						Name: "args",
						Type: "List[str]",
					},
				},
				Body: []Statement{
					&Assignment{
						Variable: "a",
						Value:    "1",
					},
				},
			},
		},
		Classes: []*Class{
			{
				Name: "Test",
				Fields: []Field{
					{
						Name: "a",
						Type: "int",
					},
				},
				Methods: []*Function{
					{
						Name: "test",
						Params: []Parameter{
							{
								Name: "self",
							},
						},
						Body: []Statement{
							&Assignment{
								Variable: "self.a",
								Value:    "1",
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	t.Logf("\n%s", code)

	_, err = gen.Generate(&File{
		Name: "",
		Classes: []*Class{
			{
				Name: "",
			},
		},
	})

	if err == nil {
		t.Error("expected error")
	}

}

func TestPythonCodeGenerator_reset(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	cg.sb.WriteString("blabla")
	cg.currentLine = 2
	cg.currentColumn = 3
	cg.indentLevel = 5
	cg.reset()
	test.AssertEqual(t, cg.sb.String(), "", "StringBuilder should be reset")
	test.AssertEqual(t, cg.currentLine, 1)
	test.AssertEqual(t, cg.currentColumn, 0)
	test.AssertEqual(t, cg.indentLevel, 0)
}

func TestPythonCodeGenerator_NewPythonCodeGenerator(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	test.AssertEqual(t, cg.currentLine, 1)
	test.AssertEqual(t, cg.currentColumn, 0)
	test.AssertEqual(t, cg.indentLevel, 0)
	test.AssertEqual(t, cg.indentFourSpaces, true)
	test.AssertEqual(t, cg.sb.String(), "")
}

func TestPythonCodeGenerator_append(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	cg.append("test")
	test.AssertEqual(t, cg.sb.String(), "test")

	test.AssertEqual(t, cg.currentLine, 1)
	test.AssertEqual(t, cg.currentColumn, len("test"))

	cg.append("test\n")

	test.AssertEqual(t, cg.currentLine, 2)
	test.AssertEqual(t, cg.currentColumn, 0)
	test.AssertEqual(t, cg.sb.String(), "testtest\n")
}

func TestPythonCodeGenerator_appendIndented(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	// test with indentLevel = 0
	test.AssertEqual(t, cg.indentLevel, 0)
	cg.appendIndented("test")
	test.AssertEqual(t, cg.sb.String(), "test")

	test.AssertEqual(t, cg.currentLine, 1)
	test.AssertEqual(t, cg.currentColumn, len("test"))

	cg.appendIndented("test\n")

	test.AssertEqual(t, cg.currentLine, 2)
	test.AssertEqual(t, cg.currentColumn, 0)
	test.AssertEqual(t, cg.sb.String(), "testtest\n")

	// test with indentLevel = 1
	cg.up()
	test.AssertEqual(t, cg.indentLevel, 1)
	cg.appendIndented("test")
	test.AssertEqual(t, cg.sb.String(), "testtest\n    test")
	test.AssertEqual(t, cg.currentLine, 2)
	test.AssertEqual(t, cg.currentColumn, len("    test"))

	cg.up()

	cg.appendIndented("test\n")
	test.AssertEqual(t, cg.sb.String(), "testtest\n    test        test\n")
	test.AssertEqual(t, cg.currentLine, 3)
	test.AssertEqual(t, cg.currentColumn, 0)

	// test with indentFourSpaces = false
	cg = NewPythonCodeGenerator(false)
	cg.appendIndented("test")
	test.AssertEqual(t, cg.sb.String(), "test")
	test.AssertEqual(t, cg.currentLine, 1)
	test.AssertEqual(t, cg.currentColumn, len("test"))

	// test with indentLevel = 1
	cg.up()

	test.AssertEqual(t, cg.indentLevel, 1)
	cg.appendIndented("test\n")

	test.AssertEqual(t, cg.currentLine, 2)
	test.AssertEqual(t, cg.currentColumn, 0)
	test.AssertEqual(t, cg.sb.String(), "test        test\n")

	// test with indentLevel = 2
	cg.up()

	test.AssertEqual(t, cg.indentLevel, 2)
	cg.appendIndented("test")

	test.AssertEqual(t, cg.sb.String(), "test        test\n                test")
	test.AssertEqual(t, cg.currentLine, 2)
	test.AssertEqual(t, cg.currentColumn, len("                test"))
}

func TestPythonCodeGenerator_up(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	cg.up()
	test.AssertEqual(t, cg.indentLevel, 1)
	cg.up()
	test.AssertEqual(t, cg.indentLevel, 2)
}

func TestPythonCodeGenerator_down(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	cg.down()
	test.AssertEqual(t, cg.indentLevel, 0)
	cg.down()
	test.AssertEqual(t, cg.indentLevel, 0)
}

func TestPythonCodeGenerator_newLine(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	cg.newLine()
	test.AssertEqual(t, cg.currentLine, 2)
	test.AssertEqual(t, cg.currentColumn, 0)
	test.AssertEqual(t, cg.sb.String(), "\n")
}

func TestPythonCodeGenerator_VisitClass_WithEmptyName(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	class := &Class{
		Name: "",
	}

	test.AssertNotEqual(t, cg.VisitClass(class), nil, "expected error")
}

func TestPythonCodeGenerator_VisitClass_WithoutBody(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	class := &Class{
		Name:   "test",
		Extend: "test",
	}

	test.AssertEqual(t, cg.VisitClass(class), nil)
	test.AssertEqual(t, cg.sb.String(), "class test(test):\n    pass\n")
}

func TestPythonCodeGenerator_VisitClass_WithFieldsError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	class := &Class{
		Name: "test",
		Fields: []Field{
			{
				Name: "",
				Type: "int",
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitClass(class), nil, "expected error")
}

func TestPythonCodeGenerator_VisitClass_WithMethodsError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	class := &Class{
		Name: "test",
		Methods: []*Function{
			{
				Name: "",
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitClass(class), nil, "expected error")
}

func TestPythonCodeGenerator_VisitClass(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	class := &Class{
		Name: "Test",
		Fields: []Field{
			{
				Name: "a",
				Type: "int",
			},
		},
		Methods: []*Function{
			{
				Name: "test",
				Params: []Parameter{
					{
						Name: "self",
					},
				},
				Body: []Statement{
					&Assignment{
						Variable: "self.a",
						Value:    "1",
					},
				},
			},
		},
	}

	test.AssertEqual(t, cg.VisitClass(class), nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "class Test:\n    a: int\n\n    def test(self):\n        self.a = 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitField(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	field := &Field{
		Name: "",
	}

	test.AssertNotEqual(t, cg.VisitField(field), nil, "expected error")

	field = &Field{
		Name: "test",
		Type: "",
	}
	test.AssertNotEqual(t, cg.VisitField(field), nil, "expected error")

	field = &Field{
		Name: "test",
		Type: "test",
	}

	test.AssertEqual(t, cg.VisitField(field), nil)

	cg.sb.Reset()

	field = &Field{
		Name: "a",
		Type: "int",
	}

	test.AssertEqual(t, cg.VisitField(field), nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "a: int\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitFunction_WithEmptyName(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	function := &Function{
		Name: "",
	}

	test.AssertNotEqual(t, cg.VisitFunction(function), nil, "error expected")
}

func TestPythonCodeGenerator_VisitFunction_WithEmptyParamsAndEmptyBody(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	function := &Function{
		Name: "test",
	}

	test.AssertEqual(t, cg.VisitFunction(function), nil, "no error expected")

	// the output string must contains pass (added automatically by the generator)
	test.AssertEqual(t, cg.sb.String(), "def test():\n    pass\n")
}

func TestPythonCodeGenerator_VisitFunction_WithBodyAndParams(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	function := &Function{
		Name: "test",
		Params: []Parameter{
			{
				Name: "args",
				Type: "List[str]",
			},
		},
		Body: []Statement{
			&Assignment{
				Variable: "a",
				Value:    "1",
			},
		},
	}

	test.AssertEqual(t, cg.VisitFunction(function), nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "def test(args: List[str]):\n    a = 1\n")
}

func TestPythonCodeGenerator_VisitFunction_WithMultipleParams(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	function := &Function{
		Name: "test",
		Params: []Parameter{
			{
				Name: "args",
				Type: "List[str]",
			},
			{
				Name: "kwargs",
				Type: "Dict[str, str]",
			},
		},
	}

	test.AssertEqual(t, cg.VisitFunction(function), nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "def test(args: List[str], kwargs: Dict[str, str]):\n    pass\n")
}

func TestPythonCodeGenerator_VisitFunction_WithMultipleParamsAndError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	// test with multiple parameters and error
	function := &Function{
		Name: "test",
		Params: []Parameter{
			{
				Name: "args",
				Type: "List[str]",
			},
			{
				Name: "",
				Type: "Dict[str, str]",
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitFunction(function), nil, "expected error")
}

func TestPythonCodeGenerator_VisitFunction_WithImportError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	// test with import error
	function := &Function{
		Name: "test",
		Imports: []Import{
			{
				What: []ImportWhat{
					{
						Name: "",
					},
				},
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitFunction(function), nil, "expected error")
}

func TestPythonCodeGenerator_VisitFunction_WithBodyError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	// test with body error
	function := &Function{
		Name: "test",
		Body: []Statement{
			&Assignment{
				Variable: "a",
				Type:     "int",
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitFunction(function), nil, "expected error")
}

func TestPythonCodeGenerator_VisitFunction_WithReturnType(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	function := &Function{
		Name:       "test",
		ReturnType: "int",
	}

	test.AssertEqual(t, cg.VisitFunction(function), nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "def test() -> int:\n    pass\n")
}

func TestPythonCodeGenerator_VisitImport(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	imp := &Import{}

	err := cg.VisitImport(imp)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	imp = &Import{
		What: []ImportWhat{
			{
				Name: "os",
			},
		},
	}

	test.AssertEqual(t, cg.VisitImport(imp), nil)

	cg.sb.Reset()

	imp = &Import{
		What: []ImportWhat{
			{
				Name:  "os",
				Alias: "o",
			},
		},
	}

	err = cg.VisitImport(imp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "import os as o\n")

	cg.sb.Reset()

	imp = &Import{
		What: []ImportWhat{
			{
				Name:  "os",
				Alias: "o",
			},
			{
				Name: "sys",
			},
			{
				Name:  "List",
				Alias: "L",
			},
		},
	}

	err = cg.VisitImport(imp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "import os as o, sys, List as L\n")
}

func TestPythonCodeGenerator_VisitImportWhat(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	imp := &ImportWhat{}

	err := cg.VisitImportWhat(imp)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	imp = &ImportWhat{
		Name: "os",
	}

	test.AssertEqual(t, cg.VisitImportWhat(imp), nil)

	cg.sb.Reset()

	imp = &ImportWhat{
		Name:  "os",
		Alias: "o",
	}

	err = cg.VisitImportWhat(imp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "os as o")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitParameter(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	param := &Parameter{
		Name: "",
	}

	err := cg.VisitParameter(param)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	param = &Parameter{
		Name: "test",
	}

	test.AssertEqual(t, cg.VisitParameter(param), nil)

	cg.sb.Reset()

	param = &Parameter{
		Name: "args",
		Type: "List[str]",
	}

	err = cg.VisitParameter(param)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "args: List[str]")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitAssignment(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	assign := &Assignment{}
	err := cg.VisitAssignment(assign)

	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	assign = &Assignment{
		Variable: "a",
		Value:    "1",
	}

	test.AssertEqual(t, cg.VisitAssignment(assign), nil)

	assign = &Assignment{
		Variable: "a",
		Type:     "int",
		Value:    "",
	}

	test.AssertNotEqual(t, cg.VisitAssignment(assign), nil, "expected error")

	cg.sb.Reset()

	assign = &Assignment{
		Variable: "a",
		Type:     "int",
		Value:    "1",
	}

	err = cg.VisitAssignment(assign)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "a: int = 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitFile(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	file := &File{
		Name: "",
		Imports: []Import{
			{
				What: []ImportWhat{
					{
						Name: "",
					},
				},
			},
		},
	}

	if err := cg.VisitFile(file); err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	file = &File{
		Name: "test.py",
		Functions: []*Function{
			{
				Name: "",
			},
		},
	}

	if err := cg.VisitFile(file); err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	err := cg.VisitFile(file)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	file = &File{
		Name: "test.py",
		HeaderComments: []string{
			"Code generated by EMF",
			"DO NOT EDIT!",
		},
		Imports: []Import{
			{
				What: []ImportWhat{
					{
						Name: "os",
					},
				},
			},
			{
				What: []ImportWhat{
					{
						Name: "List",
					},
				},
				From: "typing",
			},
		},
		Functions: []*Function{
			{
				Name: "main",
				Params: []Parameter{
					{
						Name: "args",
						Type: "List[str]",
					},
				},
				Body: []Statement{
					&Assignment{
						Variable: "a",
						Value:    "1",
					},
				},
			},
		},
		Classes: []*Class{
			{
				Name: "Test",
				Fields: []Field{
					{
						Name: "a",
						Type: "int",
					},
				},
				Methods: []*Function{
					{
						Name: "test",
						Params: []Parameter{
							{
								Name: "self",
							},
						},
						Body: []Statement{
							&Assignment{
								Variable: "self.a",
								Value:    "1",
							},
						},
					},
				},
			},
		},
	}

	err = cg.VisitFile(file)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("\n%s", cg.sb.String())
}