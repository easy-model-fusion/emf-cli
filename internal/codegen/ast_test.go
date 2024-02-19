package codegen

import "testing"

func TestAssignment_Accept(t *testing.T) {
	a := &Assignment{}
	v := newTestVisitor()
	err := a.Accept(v)
	if err != nil {
		t.Error(err)
	}
	if !v.visits["assignment"] {
		t.Error("VisitAssignment not called")
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

func (v *testVisitor) VisitAssignment(*Assignment) error {
	v.visits["assignment"] = true
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
