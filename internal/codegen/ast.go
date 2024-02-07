package codegen

type Node interface {
	Accept(visitor PythonVisitor) error // Accept method to allow visitors
}

type File struct {
	Name      string
	Imports   []Import
	Functions []*Function
	Classes   []*Class
}

type Function struct {
	Name       string
	ReturnType string
	Imports    []Import
	Params     []Parameter
	Body       []Statement
}

type Class struct {
	Name    string
	Extend  string
	Fields  []Field
	Methods []Function
}

type Statement interface {
}

type Expression interface {
}

type Parameter struct {
	Name string
	Type string
}

type Field struct {
	Name string
	Type string
}

type Import struct {
	What string
	From string
}

// Accept method for File
func (m *File) Accept(visitor PythonVisitor) error {
	return visitor.VisitFile(m)
}

// Accept method for Class
func (c *Class) Accept(visitor PythonVisitor) error {
	return visitor.VisitClass(c)
}

// Accept method for Function
func (f *Function) Accept(visitor PythonVisitor) error {
	return visitor.VisitFunction(f)
}

// Accept method for Field
func (f *Field) Accept(visitor PythonVisitor) error {
	return visitor.VisitField(f)
}

// Accept method for Parameter
func (p *Parameter) Accept(visitor PythonVisitor) error {
	return visitor.VisitParameter(p)
}

// Accept method for Import
func (i *Import) Accept(visitor PythonVisitor) error {
	return visitor.VisitImport(i)
}
