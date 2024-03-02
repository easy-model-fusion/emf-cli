package codegen

import "testing"

func TestImport_Equals(t *testing.T) {
	i1 := &Import{
		From:  "from",
		Alias: "alias",
	}
	i2 := &Import{
		From:  "from",
		Alias: "alias",
	}
	if !i1.Equals(i2) {
		t.Error("i1 should equal i2")
	}

	i1.What = []ImportWhat{
		{
			Name:  "name",
			Alias: "alias",
		},
		{
			Name:  "test",
			Alias: "testalias",
		},
	}

	if i1.Equals(i2) {
		t.Error("i1 should not equal i2")
	}

	i2.What = []ImportWhat{
		{
			Name:  "name",
			Alias: "alias",
		},
		{
			Name:  "test",
			Alias: "testalias",
		},
	}

	if !i1.Equals(i2) {
		t.Error("i1 should equal i2")
	}

	i1.What[1].Alias = "alias"

	if i1.Equals(i2) {
		t.Error("i1 should not equal i2")
	}

}

func TestAssignment_Accept(t *testing.T) {
	a := &AssignmentStmt{}
	v := newTestVisitor()
	err := a.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["assignment_stmt"] {
		t.Error("VisitAssignmentStmt not called")
	}
}

func TestClass_Accept(t *testing.T) {
	c := &Class{}
	v := newTestVisitor()
	err := c.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["class"] {
		t.Error("VisitClass not called")
	}
}

func TestFile_Accept(t *testing.T) {
	f := &File{}
	v := newTestVisitor()
	err := f.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["file"] {
		t.Error("VisitFile not called")
	}
}

func TestFunction_Accept(t *testing.T) {
	f := &Function{}
	v := newTestVisitor()
	err := f.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["function"] {
		t.Error("VisitFunction not called")
	}
}

func TestParameter_Accept(t *testing.T) {
	p := &Parameter{}
	v := newTestVisitor()
	err := p.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["parameter"] {
		t.Error("VisitParameter not called")
	}
}

func TestField_Accept(t *testing.T) {
	f := &Field{}
	v := newTestVisitor()
	err := f.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["field"] {
		t.Error("VisitField not called")
	}
}

func TestImport_Accept(t *testing.T) {
	i := &Import{}
	v := newTestVisitor()
	err := i.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["import"] {
		t.Error("VisitImport not called")
	}
}

func TestImportWhat_Accept(t *testing.T) {
	i := &ImportWhat{}
	v := newTestVisitor()
	err := i.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["import_what"] {
		t.Error("VisitImportWhat not called")
	}
}

func TestFunctionCallStmt_Accept(t *testing.T) {
	f := &FunctionCallStmt{}
	v := newTestVisitor()
	err := f.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["function_call_stmt"] {
		t.Error("VisitFunctionCallStmt not called")
	}
}

func TestFunctionCallParameter_Accept(t *testing.T) {
	f := &FunctionCallParameter{}
	v := newTestVisitor()
	err := f.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["function_call_parameter"] {
		t.Error("VisitFunctionCallParameter not called")
	}
}

func TestCommentStmt_Accept(t *testing.T) {
	c := &CommentStmt{}
	v := newTestVisitor()
	err := c.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["comment_stmt"] {
		t.Error("VisitCommentStmt not called")
	}
}

func TestReturnStmt_Accept(t *testing.T) {
	r := &ReturnStmt{}
	v := newTestVisitor()
	err := r.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["return_stmt"] {
		t.Error("VisitReturnStmt not called")
	}
}

func TestFunctionCall_Accept(t *testing.T) {
	f := &FunctionCall{}
	v := newTestVisitor()
	err := f.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["function_call"] {
		t.Error("VisitFunctionCall not called")
	}
}

func TestIfStmt_Accept(t *testing.T) {
	i := &IfStmt{}
	v := newTestVisitor()
	err := i.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["if_stmt"] {
		t.Error("VisitIfStmt not called")
	}
}

func TestElifStmt_Accept(t *testing.T) {
	e := &ElifStmt{}
	v := newTestVisitor()
	err := e.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["elif_stmt"] {
		t.Error("VisitElifStmt not called")
	}
}

func TestElseStmt_Accept(t *testing.T) {
	e := &ElseStmt{}
	v := newTestVisitor()
	err := e.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["else_stmt"] {
		t.Error("VisitElseStmt not called")
	}
}

type testVisitor struct {
	visits map[string]bool
}

func newTestVisitor() *testVisitor {
	return &testVisitor{
		visits: map[string]bool{},
	}
}

func (v *testVisitor) VisitFile(*File) error {
	v.visits["file"] = true
	return nil
}

func (v *testVisitor) VisitFunction(*Function) error {
	v.visits["function"] = true
	return nil
}

func (v *testVisitor) VisitClass(*Class) error {
	v.visits["class"] = true
	return nil
}

func (v *testVisitor) VisitField(*Field) error {
	v.visits["field"] = true
	return nil
}

func (v *testVisitor) VisitParameter(*Parameter) error {
	v.visits["parameter"] = true
	return nil
}

func (v *testVisitor) VisitImport(*Import) error {
	v.visits["import"] = true
	return nil
}

func (v *testVisitor) VisitImportWhat(*ImportWhat) error {
	v.visits["import_what"] = true
	return nil
}

func (v *testVisitor) VisitAssignmentStmt(*AssignmentStmt) error {
	v.visits["assignment_stmt"] = true
	return nil
}

func (v *testVisitor) VisitStatement(*Statement) error {
	v.visits["statement"] = true
	return nil
}

func (v *testVisitor) VisitExpression(*Expression) error {
	v.visits["expression"] = true
	return nil
}

func (v *testVisitor) VisitFunctionCallStmt(*FunctionCallStmt) error {
	v.visits["function_call_stmt"] = true
	return nil
}

func (v *testVisitor) VisitFunctionCallParameter(*FunctionCallParameter) error {
	v.visits["function_call_parameter"] = true
	return nil
}

func (v *testVisitor) VisitFunctionCall(*FunctionCall) error {
	v.visits["function_call"] = true
	return nil
}

func (v *testVisitor) VisitCommentStmt(*CommentStmt) error {
	v.visits["comment_stmt"] = true
	return nil
}

func (v *testVisitor) VisitReturnStmt(*ReturnStmt) error {
	v.visits["return_stmt"] = true
	return nil
}

func (v *testVisitor) VisitIfStmt(*IfStmt) error {
	v.visits["if_stmt"] = true
	return nil
}

func (v *testVisitor) VisitElifStmt(*ElifStmt) error {
	v.visits["elif_stmt"] = true
	return nil
}

func (v *testVisitor) VisitElseStmt(*ElseStmt) error {
	v.visits["else_stmt"] = true
	return nil
}
