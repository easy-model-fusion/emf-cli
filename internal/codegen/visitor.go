package codegen

// Generator uses the visitor pattern to generate code from an AST

type PythonVisitor interface {
	VisitFile(*File) error
	VisitFunction(*Function) error
	VisitClass(*Class) error
	VisitField(*Field) error
	VisitParameter(*Parameter) error
	VisitImport(*Import) error
	VisitImportWhat(*ImportWhat) error

	// Statements
	VisitAssignmentStmt(*AssignmentStmt) error
	VisitFunctionCallStmt(*FunctionCallStmt) error
	VisitFunctionCallStmtParameter(*FunctionCallStmtParameter) error
	VisitCommentStmt(*CommentStmt) error
	VisitReturnStmt(*ReturnStmt) error
}
