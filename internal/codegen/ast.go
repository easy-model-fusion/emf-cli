package codegen

type Node interface {
	Accept(visitor PythonVisitor) error // Accept method to allow visitors
}

type File struct {
	Name           string
	HeaderComments []string
	Imports        []Import
	Functions      []*Function
	Classes        []*Class
}

type Function struct {
	Name       string
	ReturnType string
	Params     []Parameter
	Imports    []Import
	Body       []Statement
}

type Class struct {
	Name       string
	Extend     string
	Fields     []Field
	Statements []Statement
	Methods    []*Function
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
	What  []ImportWhat
	From  string
	Alias string
}

type ImportWhat struct {
	Name  string
	Alias string
}

// Statements

type Statement interface {
	Node
}

type Assignment struct {
	Variable string
	Type     string
	Value    string
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

// Accept method for ImportWhat
func (i *ImportWhat) Accept(visitor PythonVisitor) error {
	return visitor.VisitImportWhat(i)
}

// Accept method for Assignment
func (s *Assignment) Accept(visitor PythonVisitor) error {
	return visitor.VisitAssignment(s)
}
