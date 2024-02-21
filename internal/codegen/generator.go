package codegen

import (
	"errors"
	"fmt"
	"strings"
)

type PythonCodeGenerator struct {
	sb               *strings.Builder
	indentLevel      int  // current indentation level
	indentFourSpaces bool // if true, use 4 spaces for indentation, otherwise use 8 spaces
	currentLine      int  // current line number
	currentColumn    int  // current column number
}

// NewPythonCodeGenerator creates a new PythonCodeGenerator
// indentFourSpaces: if true, use 4 spaces for indentation, otherwise use 8 spaces
func NewPythonCodeGenerator(indentFourSpaces bool) *PythonCodeGenerator {
	return &PythonCodeGenerator{
		sb:               &strings.Builder{},
		indentLevel:      0,
		currentLine:      1,
		indentFourSpaces: indentFourSpaces,
	}
}

// reset is setting all values back to default
func (cg *PythonCodeGenerator) reset() {
	cg.sb.Reset()
	cg.indentLevel = 0
	cg.currentLine = 1
	cg.currentColumn = 0
}

// up increases the indentation level
func (cg *PythonCodeGenerator) up() {
	cg.indentLevel++
}

// down decreases the indentation level
func (cg *PythonCodeGenerator) down() {
	cg.indentLevel--
	if cg.indentLevel < 0 {
		cg.indentLevel = 0
	}
}

// appendIndented appends a line to the generated code with the correct indentation
func (cg *PythonCodeGenerator) appendIndented(line string) {
	var tab = "    "
	if !cg.indentFourSpaces {
		tab = strings.Repeat(tab, 2)
	}
	indentation := strings.Repeat(tab, cg.indentLevel)
	cg.sb.WriteString(indentation + line)
	if strings.HasSuffix(line, "\n") {
		cg.currentLine++
		cg.currentColumn = 0
		return
	}
	cg.currentColumn += len(indentation) + len(line)
}

// append appends a line to the generated code
func (cg *PythonCodeGenerator) append(line string) {
	cg.sb.WriteString(line)
	if strings.HasSuffix(line, "\n") {
		cg.currentLine++
		cg.currentColumn = 0
		return
	}
	cg.currentColumn += len(line)
}

// newLine adds a new line to the generated code
func (cg *PythonCodeGenerator) newLine() {
	cg.currentLine++
	cg.currentColumn = 0
	cg.sb.WriteString("\n")
}

// VisitFile visits a File node
func (cg *PythonCodeGenerator) VisitFile(file *File) error {
	for _, comment := range file.HeaderComments {
		cg.append("# " + comment + "\n")
	}

	if len(file.HeaderComments) > 0 {
		cg.newLine()
	}

	for _, importStmt := range file.Imports {
		err := importStmt.Accept(cg)
		if err != nil {
			return err
		}
	}

	if len(file.Imports) > 0 {
		cg.newLine()
	}

	for _, class := range file.Classes {
		err := class.Accept(cg)
		if err != nil {
			return err
		}
	}

	if len(file.Classes) > 0 {
		cg.newLine()
	}

	for _, function := range file.Functions {
		err := function.Accept(cg)
		if err != nil {
			return err
		}
	}

	if len(file.Functions) > 0 {
		cg.newLine()
	}

	return nil
}

// VisitFunction visits a Function node
func (cg *PythonCodeGenerator) VisitFunction(function *Function) error {
	cg.appendIndented("def ")

	if function.Name == "" {
		return errors.New("function name cannot be empty")
	}

	cg.append(function.Name + "(")

	defaultFound := false

	for i, param := range function.Params {
		if param.Default != "" {
			defaultFound = true
		}

		err := param.Accept(cg)
		if err != nil {
			return err
		}

		if param.Default == "" && defaultFound {
			return errors.New("non-default argument follows default argument")
		}

		if i < len(function.Params)-1 {
			cg.append(", ")
		}

	}

	if function.ReturnType != "" {
		cg.append(") -> " + function.ReturnType + ":\n")
	} else {
		cg.append("):\n")
	}

	cg.up()

	for _, imp := range function.Imports {
		err := imp.Accept(cg)
		if err != nil {
			return err
		}
	}

	for _, stmt := range function.Body {
		err := stmt.Accept(cg)
		if err != nil {
			return err
		}
	}

	// If the function has no body and no imports, add a pass statement
	if len(function.Body) == 0 && len(function.Imports) == 0 {
		cg.appendIndented("pass\n")
	}

	cg.down()

	return nil
}

// VisitClass visits a Class node
func (cg *PythonCodeGenerator) VisitClass(class *Class) error {
	cg.appendIndented("class ")

	if class.Name == "" {
		return errors.New("class name cannot be empty")
	}

	cg.append(class.Name)

	if class.Extend != "" {
		cg.append("(" + class.Extend + "):\n")
	} else {
		cg.append(":\n")
	}

	cg.up()

	for _, field := range class.Fields {
		err := field.Accept(cg)
		if err != nil {
			return err
		}
	}

	if len(class.Fields) > 0 {
		cg.newLine()
	}

	for _, stmt := range class.Statements {
		err := stmt.Accept(cg)
		if err != nil {
			return err
		}
	}

	if len(class.Statements) > 0 {
		cg.newLine()
	}

	for _, method := range class.Methods {
		err := method.Accept(cg)
		if err != nil {
			return err
		}
	}

	// If the class has no fields and no methods, add a pass statement
	if len(class.Fields) == 0 && len(class.Methods) == 0 && len(class.Statements) == 0 {
		cg.appendIndented("pass\n")
	}

	cg.down()

	return nil
}

// VisitField visits a Field node
func (cg *PythonCodeGenerator) VisitField(field *Field) error {
	if field.Name == "" {
		return errors.New("field name cannot be empty")
	}

	cg.appendIndented(field.Name + ": ")

	if field.Type == "" {
		return errors.New("field type cannot be empty")
	}

	cg.append(field.Type + "\n")

	return nil
}

// VisitParameter visits a Parameter node
func (cg *PythonCodeGenerator) VisitParameter(parameter *Parameter) error {
	if parameter.Name == "" {
		return errors.New("parameter name cannot be empty")
	}
	if parameter.Type == "" {
		cg.append(parameter.Name)
		return nil
	}

	if parameter.Default != "" {
		cg.append(parameter.Name + ": " + parameter.Type + " = " + parameter.Default)
		return nil
	}

	cg.append(parameter.Name + ": " + parameter.Type)
	return nil
}

// VisitImport visits an Import node
func (cg *PythonCodeGenerator) VisitImport(importStmt *Import) error {
	if importStmt.From != "" {
		cg.appendIndented("from " + importStmt.From + " import ")
	} else {
		cg.appendIndented("import ")
	}

	if len(importStmt.What) == 0 {
		return errors.New("import statement must have at least one item")
	}

	for i, what := range importStmt.What {
		if i > 0 && i < len(importStmt.What) {
			cg.append(", ")
		}
		err := what.Accept(cg)
		if err != nil {
			return err
		}
	}

	cg.newLine()

	return nil
}

// VisitImportWhat visits an ImportWhat node
func (cg *PythonCodeGenerator) VisitImportWhat(importWhat *ImportWhat) error {
	if importWhat.Name == "" {
		return errors.New("import what name cannot be empty")
	}

	cg.append(importWhat.Name)

	if importWhat.Alias != "" {
		cg.append(" as " + importWhat.Alias)
	}

	return nil
}

// VisitAssignmentStmt visits an AssignmentStmt node
func (cg *PythonCodeGenerator) VisitAssignmentStmt(assignment *AssignmentStmt) error {
	if assignment.Variable == "" {
		return errors.New("assignment variable cannot be empty")
	}

	if assignment.Type != "" {
		cg.appendIndented(assignment.Variable + ": " + assignment.Type + " = ")
	} else {
		cg.appendIndented(assignment.Variable + " = ")
	}

	if assignment.FunctionCallValue != nil && assignment.StringValue != "" {
		return errors.New("assignment cannot have both function call and string value")
	}

	if assignment.FunctionCallValue == nil && assignment.StringValue == "" {
		return errors.New("assignment must have either function call or string value")
	}

	if assignment.FunctionCallValue != nil {
		err := assignment.FunctionCallValue.Accept(cg)
		if err != nil {
			return err
		}
		return nil
	}

	cg.append(assignment.StringValue + "\n")

	return nil
}

// VisitCommentStmt visits a CommentStmt node
func (cg *PythonCodeGenerator) VisitCommentStmt(comment *CommentStmt) error {
	if len(comment.Lines) == 0 {
		return errors.New("comment must have at least one line")
	}

	// Single line comment
	if len(comment.Lines) == 1 {
		cg.appendIndented("# " + comment.Lines[0] + "\n")
		return nil
	}

	// Multi line comment
	cg.appendIndented("\"\"\"\n")

	for _, line := range comment.Lines {
		cg.appendIndented(line + "\n")
	}

	cg.appendIndented("\"\"\"\n")

	return nil
}

// VisitFunctionCallStmt visits a FunctionCallStmt node
func (cg *PythonCodeGenerator) VisitFunctionCallStmt(functionCall *FunctionCallStmt) error {
	cg.appendIndented("")
	return functionCall.FunctionCall.Accept(cg)
}

// VisitFunctionCall visits a FunctionCall node
func (cg *PythonCodeGenerator) VisitFunctionCall(functionCall *FunctionCall) error {
	if functionCall.Name == "" {
		return errors.New("function call name cannot be empty")
	}

	cg.append(functionCall.Name + "(\n")
	cg.up() // cool looking code

	positionalFound := false

	for i, param := range functionCall.Params {
		if param.Name != "" {
			positionalFound = true
		}
		cg.appendIndented("")

		err := param.Accept(cg)
		if err != nil {
			return err
		}

		if param.Name == "" && positionalFound {
			return errors.New("positional argument follows keyword argument")
		}

		if i < len(functionCall.Params)-1 {
			cg.append(",\n")
		}
	}

	cg.down()
	cg.newLine()
	cg.append(")\n")

	return nil
}

// VisitFunctionCallParameter visits a FunctionCallParameter node
func (cg *PythonCodeGenerator) VisitFunctionCallParameter(functionCallParameter *FunctionCallParameter) error {
	if functionCallParameter.Value == "" {
		return errors.New("function call parameter value cannot be empty")
	}

	if functionCallParameter.Name == "" {
		cg.append(functionCallParameter.Value)
		return nil
	}

	cg.append(functionCallParameter.Name + " = " + functionCallParameter.Value)

	return nil
}

// VisitReturnStmt visits a ReturnStmt node
func (cg *PythonCodeGenerator) VisitReturnStmt(returnStmt *ReturnStmt) error {
	if returnStmt.Value == "" {
		cg.appendIndented("return\n")
		return nil
	}

	cg.appendIndented("return " + returnStmt.Value + "\n")
	return nil
}

// Generate generates the code from the AST
func (cg *PythonCodeGenerator) Generate(file *File) (string, error) {
	cg.reset()

	err := file.Accept(cg)
	if err != nil {
		cg.sb.WriteString("\n")

		if cg.currentColumn >= 1 {
			cg.sb.WriteString(strings.Repeat("~", cg.currentColumn-1))
		}

		cg.sb.WriteString("^^^\n")

		err = fmt.Errorf("error generating code (L%d, Col%d): %w", cg.currentLine, cg.currentColumn, err)
	}

	return cg.sb.String(), err
}
