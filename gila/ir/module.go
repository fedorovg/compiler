package ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
	"os"
)

type Module struct {
	*ir.Module
	functions map[string]*Function
}

func NewModule(program *ast.Program) *Module {
	module := &Module{ir.NewModule(), make(map[string]*Function)}
	module.SourceFilename = program.Name
	module.declareStl()
	for _, f := range program.Functions {
		module.functions[f.Signature.Name] = module.newFunction(f)
	}
	return module
}

func (m *Module) newFunction(function *ast.Function) *Function {
	var f *Function
	// If function has already been declared, but hasn't been implemented yet.
	if val, ok := m.functions[function.Signature.Name]; ok {
		f = val
	} else {
		f = &Function{
			Func:      m.createFuncFromSignature(function.Signature),
			functions: m.functions,
			context:   nil,
		}
		m.functions[f.Name()] = f
	}
	if function.Body != nil {
		f.emit(function)
	}
	return f
}

// declareStl simply emits ir for pre-defined functions
func (m *Module) declareStl() {
	i32 := types.I32

	wl := m.NewFunc("writeln", i32, ir.NewParam("x", i32))
	m.functions["writeln"] = &Function{Func: wl}

	w := m.NewFunc("write", i32, ir.NewParam("x", types.I8Ptr))
	m.functions["write"] = &Function{Func: w}

	rl := m.NewFunc("readln", i32, ir.NewParam("x", types.I32Ptr))
	m.functions["readln"] = &Function{Func: rl}

	// Manual implementations of increment and decrement functions
	inc := m.NewFunc("inc", i32, ir.NewParam("x", types.I32Ptr))
	incBody := inc.NewBlock("entry")
	incBody.NewStore(incBody.NewAdd(incBody.NewLoad(i32, inc.Params[0]), constant.NewInt(i32, 1)), inc.Params[0])
	incBody.NewRet(constant.NewInt(i32, 0))
	m.functions["inc"] = &Function{Func: inc}

	dec := m.NewFunc("dec", i32, ir.NewParam("x", types.I32Ptr))
	decBody := dec.NewBlock("entry")
	decBody.NewStore(decBody.NewSub(decBody.NewLoad(i32, dec.Params[0]), constant.NewInt(i32, 1)), dec.Params[0])
	decBody.NewRet(constant.NewInt(i32, 0))
	m.functions["dec"] = &Function{Func: dec}
}

func (m *Module) createFuncFromSignature(s *ast.Signature) *ir.Func {
	var params []*ir.Param
	for _, p := range s.Parameters {
		params = append(params, ir.NewParam(p.Name, types.I32))
	}
	return m.NewFunc(s.Name, types.I32, params...)
}

func (m Module) String() string {
	return m.Module.String()
}

func (m *Module) DumpToFile(path string) {
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := fo.Write([]byte(m.Module.String())); err != nil {
		panic(err)
	}
}
