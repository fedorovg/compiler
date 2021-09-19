package ir

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
)

type Context struct {
	*ir.Block
	parent      *Context
	breakTarget *Context
	function    *Function
	symbols     map[string]value.Value
}

func (c *Context) lookup(name string) value.Value {
	if v, ok := c.symbols[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.lookup(name)
	}
	return nil
}

func (c *Context) newChildContext(name string) *Context {
	return &Context{
		Block:    c.function.NewBlock(name),
		parent:   c,
		breakTarget: c.breakTarget,
		function: c.function,
		symbols:  map[string]value.Value{},
	}
}

func (c *Context) newChildScope(name string) *Context {
	return &Context{
		Block:    c.Block,
		parent:   c,
		breakTarget: c.breakTarget,
		function: c.function,
		symbols:  map[string]value.Value{},
	}
}

type Function struct {
	*ir.Func
	context   *Context
	functions map[string]*Function
}

func (f *Function) newContext(name string) *Context {
	return &Context{
		Block:    f.NewBlock(name),
		parent:   nil,
		function: f,
		symbols:  map[string]value.Value{},
	}
}

func (f *Function) emit(functionTree *ast.Function) {
	f.context = f.newContext("entry")
	f.emitVariableDeclaration(&ast.VariableDeclaration{Name: f.Name()})
	for _, param := range f.Params {
		stackParam := f.context.NewAlloca(types.I32)
		f.context.NewStore(param, stackParam)
		f.context.symbols[param.Name()] = stackParam
	}
	f.emitStatement(functionTree.Body)
	if functionTree.Signature.Return == ast.VOID {
		f.context.NewRet(constant.NewInt(types.I32, 0))
	} else {
		// This lookup is safe, since return value is guaranteed to be allocated
		l := f.context.NewLoad(types.I32, f.context.lookup(functionTree.Signature.Name))
		f.context.NewRet(l)
	}
}

func (f *Function) emitStatement(node ast.Statement) {
	switch n := node.(type) {
	case *ast.Block:
		f.emitBlock(n)
	case *ast.VariableDeclaration:
		f.emitVariableDeclaration(n)
	case *ast.ConstantDeclaration:
		f.emitConstantDeclaration(n)
	case *ast.ProcedureCall:
		f.emitProcedureCall(n)
	case *ast.If:
		f.emitIf(n)
	case *ast.While:
		f.emitWhile(n)
	case *ast.Assignment:
		f.emitAssignment(n)
	case *ast.For:
		f.emitFor(n)
	case *ast.Break:
		if f.context.breakTarget != nil { // Break outside of a loop is noop
			f.context.NewBr(f.context.breakTarget.Block)
		}
	case *ast.Exit:
		if ret := f.context.lookup(f.Name()); ret != nil {
			f.context.NewRet(f.context.NewLoad(types.I32, ret))
		} else {
			f.context.NewRet(constant.NewInt(types.I32, 0))
		}
	default:
		panic("Unknown statement type!")
	}
}

func (f *Function) emitBlock(b *ast.Block) {
	old := f.context
	newContext := f.context.newChildScope("")
	f.context = newContext
	for _, s := range b.Statements {
		f.emitStatement(s)
	}
	if f.context == newContext {
		f.context = old
	}
}

func (f *Function) emitVariableDeclaration(declaration *ast.VariableDeclaration) value.Value {
	a := f.context.NewAlloca(types.I32)
	f.context.symbols[declaration.Name] = a
	return a
}

func (f *Function) emitConstantDeclaration(declaration *ast.ConstantDeclaration) {
	f.context.symbols[declaration.Name] = constant.NewInt(types.I32, declaration.Literal.Value)
}

func (f *Function) emitProcedureCall(pc *ast.ProcedureCall) {
	pointerFunctions := map[string]bool{
		"readln": true,
		"inc":    true,
		"dec":    true,
	}
	callee := f.functions[pc.Name]
	var args []value.Value
	// Check if calle is one of the 3 built-in functions, that accept params by reference
	if pointerFunctions[callee.Name()] {
		if ptr, ok := pc.Args[0].(*ast.Variable); ok {
			args = append(args, f.emitVariablePointer(ptr))
		} else {
			panic("Syntax error inside built-in reference function call.")
		}
	} else if callee.Name() == "write" {
		tok, ok := pc.Args[0].(ast.StringLiteral)
		if !ok {
			panic("Trying to print not a string.")
		}
		constStr := constant.NewCharArrayFromString(tok.Value)
		str := f.Parent.NewGlobalDef("", constStr)
		startPtr := f.context.NewGetElementPtr(constStr.Type(), str, constant.NewInt(types.I3, 0), constant.NewInt(types.I3, 0))
		f.context.NewCall(callee, startPtr)
		return
	} else {
		for _, a := range pc.Args {
			args = append(args, f.emitExpression(a))
		}
	}
	f.context.NewCall(callee, args...)
}

func (f *Function) emitIf(ifStatement *ast.If) {
	cond := f.emitExpression(ifStatement.Condition)
	thenLabel := f.context.newChildContext("")
	// This function is WAY longer than it is supposed to be, because I wanted to keep IR labels in the same order
	// in which they appear in the source code.
	if ifStatement.Else != nil {
		elseLabel := f.context.newChildContext("")
		contLabel := f.context.newChildContext("")
		f.context.NewCondBr(cond, thenLabel.Block, elseLabel.Block)
		f.context = thenLabel
		f.emitStatement(ifStatement.Then)
		if f.context.Term == nil {
			f.context.NewBr(contLabel.Block)
		}
		f.context = elseLabel
		f.emitStatement(ifStatement.Else)
		if f.context.Term == nil {
			f.context.NewBr(contLabel.Block)
		}
		f.context = contLabel
	} else {
		contLabel := f.context.newChildContext("")
		//thenLabel.breakTarget = contLabel
		f.context.NewCondBr(cond, thenLabel.Block, contLabel.Block)
		f.context = thenLabel
		f.emitStatement(ifStatement.Then)
		if f.context.Term == nil {
			f.context.NewBr(contLabel.Block)
		}
		f.context = contLabel
	}
}

func (f *Function) emitWhile(while *ast.While) {
	cond := f.emitExpression(while.Condition)
	loopLabel := f.context.newChildContext("")
	contLabel := f.context.newChildContext("")
	loopLabel.breakTarget = contLabel
	f.context.NewCondBr(cond, loopLabel.Block, contLabel.Block)
	f.context = loopLabel
	f.emitStatement(while.Body)
	if f.context.Term == nil {
		f.context.NewCondBr(f.emitExpression(while.Condition), loopLabel.Block, contLabel.Block)
	}
	f.context = contLabel
}

func (f *Function) emitFor(forLoop *ast.For) {
	f.emitAssignment(forLoop.Initial)
	cond := ast.Binary{
		Left:      &forLoop.Initial.Variable,
		Right:     forLoop.Target,
		Operation: ast.NOTEQUALS,
	}
	var updateFunctionName = "dec"
	if forLoop.Upto {
		updateFunctionName = "inc"
	}
	forLoop.Body.Statements = append(forLoop.Body.Statements, &ast.ProcedureCall{
		Name: updateFunctionName,
		Args: []ast.Expression{&forLoop.Initial.Variable},
	})
	f.emitWhile(&ast.While{
		Condition: &cond,
		Body:      forLoop.Body,
	})
}

func (f *Function) emitExpression(expression ast.Expression) value.Value {
	switch e := expression.(type) {
	case *ast.Literal:
		return f.emitLiteral(e)
	case *ast.Variable:
		return f.emitVariable(e)
	case *ast.Binary:
		return f.emitBinary(e)
	case *ast.Unary:
		return f.emitUnary(e)
	case *ast.FunctionCall:
		return f.emitFunctionCall(e)
	default:
		panic("Not all expressions are implemented yet!")
	}
}

func (f *Function) emitAssignment(a *ast.Assignment) {
	if variable := f.context.lookup(a.Variable.Name); variable != nil {
		val := f.emitExpression(a.Value)
		f.context.NewStore(val, variable)
	} else {
		panic("Undefined symbol in assignment")
	}
}

func (f *Function) emitLiteral(l *ast.Literal) value.Value {
	return constant.NewInt(types.I32, l.Value)
}

func (f *Function) emitVariable(variable *ast.Variable) value.Value {
	if v := f.context.lookup(variable.Name); v != nil {
		switch symbol := v.(type) {
		case constant.Constant:
			return symbol
		default:
			return f.context.NewLoad(types.I32, v)
		}
	} else {
		errorMessage := fmt.Sprintf("Undefined symbol %T.", variable.Name)
		panic(errorMessage)
	}
}

func (f *Function) emitVariablePointer(variable *ast.Variable) value.Value {
	if v := f.context.lookup(variable.Name); v != nil {
		return v
	} else {
		errorMessage := fmt.Sprintf("Undefined symbol %T.", variable.Name)
		panic(errorMessage)
	}
}

func (f *Function) emitBinary(e *ast.Binary) value.Value {
	switch e.Operation {
	case ast.PLUS:
		return f.context.NewAdd(f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.MINUS:
		return f.context.NewSub(f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.MULTIPLY:
		return f.context.NewMul(f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.EQUALS:
		return f.context.NewICmp(enum.IPredEQ, f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.NOTEQUALS:
		return f.context.NewICmp(enum.IPredNE, f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.LESS:
		return f.context.NewICmp(enum.IPredSLT, f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.LESSEQ:
		return f.context.NewICmp(enum.IPredSLE, f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.GREATER:
		return f.context.NewICmp(enum.IPredSGT, f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.GREATEREQ:
		return f.context.NewICmp(enum.IPredSGE, f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.AND:
		return f.context.NewAnd(f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.OR:
		return f.context.NewOr(f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.MOD:
		return f.context.NewSRem(f.emitExpression(e.Left), f.emitExpression(e.Right))
	case ast.DIV:
		return f.context.NewSDiv(f.emitExpression(e.Left), f.emitExpression(e.Right))
	default:
		panic("Invalid operation type inside Binary node.")
	}
}

func (f *Function) emitUnary(u *ast.Unary) value.Value {
	switch u.Operation {
	case ast.PLUS:
		return f.emitExpression(u.Operand)
	case ast.MINUS:
		return f.context.NewSub(constant.NewInt(types.I32, 0), f.emitExpression(u.Operand))
	default:
		panic("Invalid operation type inside Unary node.")
	}
}

func (f *Function) emitFunctionCall(pc *ast.FunctionCall) value.Value {
	callee := f.functions[pc.Name]
	var args []value.Value
	for _, a := range pc.Args {
		args = append(args, f.emitExpression(a))
	}
	return f.context.NewCall(callee, args...)
}
