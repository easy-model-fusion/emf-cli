package codegen

import (
	"github.com/easy-model-fusion/emf-cli/test"
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
					&AssignmentStmt{
						Variable:    "a",
						StringValue: "1",
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
							&AssignmentStmt{
								Variable:    "self.a",
								StringValue: "1",
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

func TestPythonCodeGenerator_VisitClass_WithStatementsError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	class := &Class{
		Name: "test",
		Statements: []Statement{
			&AssignmentStmt{
				Variable:    "",
				StringValue: "1",
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitClass(class), nil, "expected error")
}

func TestPythonCodeGenerator_VisitClass_WithStatements(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	class := &Class{
		Name: "test",
		Statements: []Statement{
			&AssignmentStmt{
				Variable:    "a",
				StringValue: "1",
			},
		},
	}

	test.AssertEqual(t, cg.VisitClass(class), nil, "no error expected")

	test.AssertEqual(t, cg.sb.String(), "class test:\n    a = 1\n\n")

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
					&AssignmentStmt{
						Variable:    "self.a",
						StringValue: "1",
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
			&AssignmentStmt{
				Variable:    "a",
				StringValue: "1",
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
			&AssignmentStmt{
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

func TestPythonCodeGenerator_VisitFunction_WithUnorderedDefaultParams(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	function := &Function{
		Name: "test",
		Params: []Parameter{
			{
				Name:    "args",
				Type:    "List[str]",
				Default: "[]",
			},
			{
				Name: "kwargs",
				Type: "Dict[str, str]",
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitFunction(function), nil, "expected error")
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

func TestPythonCodeGenerator_VisitParameter_WithDefault(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	param := &Parameter{
		Name:    "args",
		Type:    "List[str]",
		Default: "[]",
	}

	err := cg.VisitParameter(param)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "args: List[str] = []")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitAssignment(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	assign := &AssignmentStmt{}
	err := cg.VisitAssignmentStmt(assign)

	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	assign = &AssignmentStmt{
		Variable:    "a",
		StringValue: "1",
	}

	test.AssertEqual(t, cg.VisitAssignmentStmt(assign), nil)

	assign = &AssignmentStmt{
		Variable:    "a",
		Type:        "int",
		StringValue: "",
	}

	test.AssertNotEqual(t, cg.VisitAssignmentStmt(assign), nil, "expected error")

	cg.sb.Reset()

	assign = &AssignmentStmt{
		Variable:    "a",
		Type:        "int",
		StringValue: "1",
	}

	err = cg.VisitAssignmentStmt(assign)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "a: int = 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitAssignment_WithFunctionCall(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	assign := &AssignmentStmt{
		Variable: "a",
		FunctionCallValue: &FunctionCall{
			Name: "test",
			Params: []FunctionCallParameter{
				{
					Name:  "a",
					Value: "1",
				},
			},
		},
	}

	err := cg.VisitAssignmentStmt(assign)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "a = test(\n    a = 1\n)\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitAssignment_WithFunctionCallError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	assign := &AssignmentStmt{
		Variable: "a",
		FunctionCallValue: &FunctionCall{
			Name: "test",
			Params: []FunctionCallParameter{
				{
					Name:  "a",
					Value: "",
				},
			},
		},
	}

	err := cg.VisitAssignmentStmt(assign)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}
}

func TestPythonCodeGenerator_VisitAssignment_WithValueAndFunctionCall(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	assign := &AssignmentStmt{
		Variable:    "a",
		StringValue: "aa",
		FunctionCallValue: &FunctionCall{
			Name: "test",
		},
	}
	test.AssertNotEqual(t, cg.VisitAssignmentStmt(assign), nil, "error expected")
}

func TestPythonCodeGenerator_VisitComment(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	comment := &CommentStmt{}
	err := cg.VisitCommentStmt(comment)

	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	comment = &CommentStmt{
		Lines: []string{},
	}

	err = cg.VisitCommentStmt(comment)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	comment = &CommentStmt{
		Lines: []string{
			"test",
		},
	}

	err = cg.VisitCommentStmt(comment)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "# test\n")

	cg.sb.Reset()

	comment = &CommentStmt{
		Lines: []string{
			"test",
			"test",
		},
	}

	err = cg.VisitCommentStmt(comment)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "\"\"\"\ntest\ntest\n\"\"\"\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitFunctionCallStmt(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &FunctionCallStmt{}
	err := cg.VisitFunctionCallStmt(stmt)

	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	stmt = &FunctionCallStmt{
		FunctionCall: FunctionCall{
			Name: "",
		},
	}

	err = cg.VisitFunctionCallStmt(stmt)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	stmt = &FunctionCallStmt{
		FunctionCall{
			Name: "test",
			Params: []FunctionCallParameter{
				{
					Name:  "a",
					Value: "1",
				},
			},
		},
	}

	err = cg.VisitFunctionCallStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "test(\n    a = 1\n)\n")

	cg.sb.Reset()

	stmt = &FunctionCallStmt{
		FunctionCall{
			Name: "test",
			Params: []FunctionCallParameter{
				{
					Name:  "a",
					Value: "1",
				},
				{
					Name:  "b",
					Value: "2",
				},
			},
		},
	}

	err = cg.VisitFunctionCallStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "test(\n    a = 1,\n    b = 2\n)\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitFunctionCallParameter(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	param := &FunctionCallParameter{}
	err := cg.VisitFunctionCallParameter(param)

	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	param = &FunctionCallParameter{
		Value: "",
	}

	err = cg.VisitFunctionCallParameter(param)
	if err == nil {
		t.Error("expected error")
		t.FailNow()
	}

	param = &FunctionCallParameter{
		Value: "1",
	}

	err = cg.VisitFunctionCallParameter(param)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "1")

	cg.sb.Reset()

	param = &FunctionCallParameter{
		Name:  "a",
		Value: "1",
	}

	err = cg.VisitFunctionCallParameter(param)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "a = 1")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitFunctionCall_ParamValueError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	functionCall := &FunctionCallStmt{
		FunctionCall: FunctionCall{
			Name: "test",
			Params: []FunctionCallParameter{
				{
					Name:  "a",
					Value: "",
				},
			},
		},
	}

	err := cg.VisitFunctionCallStmt(functionCall)

	test.AssertNotEqual(t, err, nil, "error expected")
}

func TestPythonCodeGenerator_VisitFunctionCall_ParamOrderError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)

	functionCall := &FunctionCallStmt{
		FunctionCall: FunctionCall{
			Name: "test",
			Params: []FunctionCallParameter{
				{
					Name:  "a",
					Value: "1",
				},
				{
					Name:  "",
					Value: "2",
				},
			},
		},
	}

	err := cg.VisitFunctionCallStmt(functionCall)

	test.AssertNotEqual(t, err, nil, "error expected")
}

func TestPythonCodeGenerator_VisitReturnStmt(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &ReturnStmt{}
	err := cg.VisitReturnStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "return\n")

	cg.sb.Reset()

	stmt = &ReturnStmt{
		Value: "",
	}

	err = cg.VisitReturnStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "return\n")

	cg.sb.Reset()

	stmt = &ReturnStmt{
		Value: "1",
	}

	err = cg.VisitReturnStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "return 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitIfStmt_WithEmptyCondition(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &IfStmt{}
	err := cg.VisitIfStmt(stmt)
	test.AssertNotEqual(t, err, nil, "expected error")
}

func TestPythonCodeGenerator_VisitIfStmt(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &IfStmt{}
	err := cg.VisitIfStmt(stmt)

	test.AssertNotEqual(t, err, nil, "expected error")
	cg.sb.Reset()

	stmt = &IfStmt{
		Condition: "a == 1",
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "a",
				StringValue: "1",
			},
		},
	}

	err = cg.VisitIfStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "if a == 1:\n    a = 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitIfStmt_WithBodyError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &IfStmt{
		Condition: "a == 1",
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "",
				StringValue: "1",
			},
		},
	}

	err := cg.VisitIfStmt(stmt)
	test.AssertNotEqual(t, err, nil, "expected error")
}

func TestPythonCodeGenerator_VisitIfStmt_WithEmptyBody(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &IfStmt{
		Condition: "a == 1",
	}

	err := cg.VisitIfStmt(stmt)
	test.AssertEqual(t, err, nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "if a == 1:\n    pass\n")
}

func TestPythonCodeGenerator_VisitIfStmt_WithElifsError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &IfStmt{
		Condition: "a == 1",
		Elifs: []*ElifStmt{
			{
				Condition: "a == 1",
				Body: []Statement{
					&AssignmentStmt{
						Variable:    "",
						StringValue: "1",
					},
				},
			},
		},
	}

	err := cg.VisitIfStmt(stmt)
	test.AssertNotEqual(t, err, nil, "expected error")
}

func TestPythonCodeGenerator_VisitIfStmt_WithElse(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &IfStmt{
		Condition: "a == 1",
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "a",
				StringValue: "1",
			},
		},
		Else: &ElseStmt{
			Body: []Statement{
				&AssignmentStmt{
					StringValue: "2",
				},
			},
		},
	}

	test.AssertNotEqual(t, cg.VisitIfStmt(stmt), nil, "expected error")

	cg.reset()

	// correct else statement
	stmt.Else.Body[0] = &AssignmentStmt{
		Variable:    "a",
		StringValue: "2",
	}

	err := cg.VisitIfStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "if a == 1:\n    a = 1\nelse:\n    a = 2\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitElifStmt(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &ElifStmt{}
	err := cg.VisitElifStmt(stmt)

	test.AssertNotEqual(t, err, nil, "expected error")
	cg.reset()

	stmt = &ElifStmt{
		Condition: "a == 1",
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "a",
				StringValue: "1",
			},
		},
	}

	err = cg.VisitElifStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "elif a == 1:\n    a = 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitElifStmt_WithEmptyBody(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &ElifStmt{
		Condition: "a == 1",
	}

	err := cg.VisitElifStmt(stmt)
	test.AssertEqual(t, err, nil, "no error expected")
	test.AssertEqual(t, cg.sb.String(), "elif a == 1:\n    pass\n")
}

func TestPythonCodeGenerator_VisitElifStmt_WithBodyError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &ElifStmt{
		Condition: "a == 1",
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "",
				StringValue: "1",
			},
		},
	}

	err := cg.VisitElifStmt(stmt)
	test.AssertNotEqual(t, err, nil, "expected error")
}

func TestPythonCodeGenerator_VisitElseStmt(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &ElseStmt{}
	err := cg.VisitElseStmt(stmt)

	test.AssertEqual(t, cg.indentLevel, 0)

	test.AssertEqual(t, err, nil, "no error expected")

	test.AssertEqual(t, cg.sb.String(), "else:\n    pass\n")

	cg.sb.Reset()

	stmt = &ElseStmt{
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "a",
				StringValue: "1",
			},
		},
	}

	err = cg.VisitElseStmt(stmt)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	test.AssertEqual(t, cg.sb.String(), "else:\n    a = 1\n")

	t.Logf("\n%s", cg.sb.String())
}

func TestPythonCodeGenerator_VisitElseStmt_WithBodyError(t *testing.T) {
	cg := NewPythonCodeGenerator(true)
	stmt := &ElseStmt{
		Body: []Statement{
			&AssignmentStmt{
				Variable:    "",
				StringValue: "1",
			},
		},
	}

	err := cg.VisitElseStmt(stmt)
	test.AssertNotEqual(t, err, nil, "expected error")
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
					&AssignmentStmt{
						Variable:    "a",
						StringValue: "1",
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
							&AssignmentStmt{
								Variable:    "self.a",
								StringValue: "1",
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
