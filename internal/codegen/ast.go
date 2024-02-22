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
	Statements []Statement
	Fields     []Field
	Methods    []*Function
}

type Expression interface {
}

type Parameter struct {
	Name    string
	Type    string
	Default string
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

type AssignmentStmt struct {
	Variable          string
	Type              string
	StringValue       string
	FunctionCallValue *FunctionCall
}

type FunctionCall struct {
	Name   string
	Params []FunctionCallParameter
}

type FunctionCallParameter struct {
	Name  string
	Value string
}

type FunctionCallStmt struct {
	FunctionCall
}

// CommentStmt One line is using # and multiple lines are using """
type CommentStmt struct {
	Lines []string
}

type ReturnStmt struct {
	Value string
}

type IfStmt struct {
	Condition string
	Body      []Statement
	Elifs     []*ElifStmt
	Else      *ElseStmt
}

type ElifStmt struct {
	Condition string
	Body      []Statement
}

type ElseStmt struct {
	Body []Statement
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

// Accept method for AssignmentStmt
func (s *AssignmentStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitAssignmentStmt(s)
}

// Accept method for FunctionComment
func (s *CommentStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitCommentStmt(s)
}

// Accept method for FunctionCallStmt
func (s *FunctionCallStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitFunctionCallStmt(s)
}

// Accept method for FunctionCall
func (s *FunctionCall) Accept(visitor PythonVisitor) error {
	return visitor.VisitFunctionCall(s)
}

// Accept method for FunctionCallParameter
func (s *FunctionCallParameter) Accept(visitor PythonVisitor) error {
	return visitor.VisitFunctionCallParameter(s)
}

// Accept method for ReturnStmt
func (s *ReturnStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitReturnStmt(s)
}

// Accept method for IfStmt
func (s *IfStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitIfStmt(s)
}

// Accept method for ElseStmt
func (s *ElseStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitElseStmt(s)
}

// Accept method for ElifStmt
func (s *ElifStmt) Accept(visitor PythonVisitor) error {
	return visitor.VisitElifStmt(s)
}
