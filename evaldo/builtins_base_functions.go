package evaldo

import (
	"github.com/refaktor/rye/env"
	"github.com/refaktor/rye/util"
	// JM 20230825	"github.com/refaktor/rye/term"
)

var builtins_functions = map[string]*env.Builtin{

	//
	// ##### Functions ##### "functions that create functions"
	//
	// Tests:
	// equal { var 'x 10 x:: 20 x } 20
	// Args:
	// * word: Tagword representing the variable name
	// * value: Initial value for the variable
	// Returns:
	// * The initial value
	"var": {
		Argsn: 2,
		Doc:   "Declares a word as a variable with the given value, allowing it to be modified. Can only be used once per word in a context.",
		Pure:  false,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch word := arg0.(type) {
			case env.Tagword:
				// Convert tagword to index
				idx := word.Index
				// Check if word already exists in context
				if _, exists := ps.Ctx.Get(idx); exists {
					ps.FailureFlag = true
					return env.NewError("Cannot redefine existing word '" + ps.Idx.GetWord(idx) + "' with var")
				}
				// Set the value
				ps.Ctx.SetNew(idx, arg1, ps.Idx)
				// Mark as variable
				ps.Ctx.MarkAsVariable(idx)
				return arg1
			case env.Word:
				// Use word index directly
				idx := word.Index
				// Check if word already exists in context
				if _, exists := ps.Ctx.Get(idx); exists {
					ps.FailureFlag = true
					return env.NewError("Cannot redefine existing word '" + ps.Idx.GetWord(idx) + "' with var")
				}
				ps.Ctx.SetNew(idx, arg1, ps.Idx)
				ps.Ctx.MarkAsVariable(idx)
				return arg1
			default:
				return MakeArgError(ps, 1, []env.Type{env.TagwordType, env.WordType}, "var")
			}
		},
	},
	// Tests:
	// equal { does { 123 } |type? } 'function
	// equal { x: does { 123 } x } 123
	// equal { x: does { 1 + 2 } x } 3
	// Args:
	// * body: Block containing the function body code
	// Returns:
	// * function object with no parameters
	"does": { // **
		Argsn: 1,
		Doc:   "Creates a function with no arguments that executes the given block when called.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch body := arg0.(type) {
			case env.Block:
				//spec := []env.Object{*env.NewWord(aaaidx)}
				//body := []env.Object{*env.NewWord(printidx), *env.NewWord(aaaidx), *env.NewWord(recuridx), *env.NewWord(greateridx), *env.NewInteger(99), *env.NewWord(aaaidx), *env.NewWord(incidx), *env.NewWord(aaaidx)}
				return *env.NewFunction(*env.NewBlock(*env.NewTSeries(make([]env.Object, 0))), body, false)
			default:
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "does")
			}
		},
	},

	// Tests:
	// equal { fn1 { .pass { } } |type? } 'function
	// equal { x: fn1 { } , x 123 } 123
	// equal { x: fn1 { .pass { } } , x 123 } 123
	// equal { x: fn1 { + 1 } , x 123 } 124
	// Args:
	// * body: Block containing the function body code
	// Returns:
	// * function object that accepts one anonymous argument
	"fn1": { // **
		Argsn: 1,
		Doc:   "Creates a function that accepts one anonymous argument and executes the given block with that argument.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch body := arg0.(type) {
			case env.Block:
				spec := []env.Object{*env.NewWord(1)}
				//body := []env.Object{*env.NewWord(printidx), *env.NewWord(aaaidx), *env.NewWord(recuridx), *env.NewWord(greateridx), *env.NewInteger(99), *env.NewWord(aaaidx), *env.NewWord(incidx), *env.NewWord(aaaidx)}
				return *env.NewFunction(*env.NewBlock(*env.NewTSeries(spec)), body, false)
			default:
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "fn1")
			}
		},
	},

	// Tests:
	// equal { fn { } { } |type? } 'function
	// equal { x: fn { } { 234 } , x } 234
	// equal { x: fn { x } { x } , x 123 } 123
	// equal { x: fn { x } { + 123 } , x 123 } 246
	// Args:
	// * spec: Block containing parameter specifications
	// * body: Block containing the function body code
	// Returns:
	// * function object with the specified parameters
	"fn": {
		Argsn: 2,
		Doc:   "Creates a function with named parameters specified in the first block and code in the second block.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch args := arg0.(type) {
			case env.Block:
				ok, doc := util.ProcessFunctionSpec(args)
				if !ok {
					return MakeBuiltinError(ps, doc, "fn")
				}
				switch body := arg1.(type) {
				case env.Block:
					//spec := []env.Object{*env.NewWord(aaaidx)}
					//body := []env.Object{*env.NewWord(printidx), *env.NewWord(aaaidx), *env.NewWord(recuridx), *env.NewWord(greateridx), *env.NewInteger(99), *env.NewWord(aaaidx), *env.NewWord(incidx), *env.NewWord(aaaidx)}
					// fmt.Println(doc)
					return *env.NewFunctionDoc(args, body, false, doc)
				default:
					return MakeArgError(ps, 2, []env.Type{env.BlockType}, "fn")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "fn")
			}
		},
	},

	// Tests:
	// equal { pfn { } { } |type? } 'function
	// equal { x: pfn { x } { + 123 } , x 123 } 246
	// ; TODO -- it seems pure namespace not also has print and append! error { x: pfn { } { ?append! } , x 123 }
	// ; TODO -- it seems pure namespace not also has print and append! error { x: pfn { x } { .print } , x 123 }
	// Args:
	// * spec: Block containing parameter specifications
	// * body: Block containing the function body code
	// Returns:
	// * pure function object with the specified parameters
	"pfn": {
		Argsn: 2,
		Doc:   "Creates a pure function (no side effects allowed) with named parameters and code body.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch args := arg0.(type) {
			case env.Block:
				ok, doc := util.ProcessFunctionSpec(args)
				if !ok {
					return MakeBuiltinError(ps, doc, "fn")
				}
				switch body := arg1.(type) {
				case env.Block:
					return *env.NewFunctionDoc(args, body, true, doc)
				default:
					return MakeArgError(ps, 2, []env.Type{env.BlockType}, "pfn")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "pfn")
			}
		},
	},

	/*
		"fnc": { // TODO -- fnc will maybe become fn\par context is set as parrent, fn\in will be executed directly in context
			// a function with context	 bb: 10 add10 [ a ] context [ b: bb ] [ add a b ]
			// 							add10 [ a ] this [ add a b ]
			// later maybe			   add10 [ a ] [ b: b ] [ add a b ]
			//  						   add10 [ a ] [ 'b ] [ add a b ]
			Argsn: 3,
			Doc:   "Creates a function with specific context.",
			Pure:  true,
			Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
				switch args := arg0.(type) {
				case env.Block:
					ok, doc := util.ProcessFunctionSpec(args)
					if !ok {
						return MakeBuiltinError(ps, doc, "fn")
					}
					switch ctx := arg1.(type) {
					case env.RyeCtx:
						switch body := arg2.(type) {
						case env.Block:
							return *env.NewFunctionC(args, body, &ctx, false, false, doc)
						default:
							ps.ErrorFlag = true
							return MakeArgError(ps, 3, []env.Type{env.BlockType}, "fnc")
						}
					default:
						ps.ErrorFlag = true
						return MakeArgError(ps, 2, []env.Type{env.CtxType}, "fnc")
					}
				default:
					ps.ErrorFlag = true
					return MakeArgError(ps, 1, []env.Type{env.BlockType}, "fnc")
				}
			},
		}, */

	// Tests:
	// equal { fn\cc { x } { x + y } |type? } 'function
	// equal { y: 5 , f: fn\cc { x } { x + y } , f 3 } 8
	// Args:
	// * spec: Block containing parameter specifications
	// * body: Block containing the function body code
	// Returns:
	// * function object with the current context captured
	"fn\\cc": {
		// a function with context	 bb: 10 add10 [ a ] context [ b: bb ] [ add a b ]
		// 							add10 [ a ] this [ add a b ]
		// later maybe			   add10 [ a ] [ b: b ] [ add a b ]
		//  						   add10 [ a ] [ 'b ] [ add a b ]
		Argsn: 2,
		Doc:   "Creates a function that captures the current context, allowing access to variables from the enclosing scope.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch args := arg0.(type) {
			case env.Block:
				ok, doc := util.ProcessFunctionSpec(args)
				if !ok {
					return MakeBuiltinError(ps, doc, "fn")
				}
				switch body := arg1.(type) {
				case env.Block:
					return *env.NewFunctionC(args, body, ps.Ctx, false, false, doc)
				default:
					ps.ErrorFlag = true
					return MakeArgError(ps, 2, []env.Type{env.BlockType}, "fn\\cc")
				}
			default:
				ps.ErrorFlag = true
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "fn\\cc")
			}
		},
	},

	// Tests:
	// equal { ctx: context { y: 5 } , f: fn\par { x } ctx { x + y } , f 3 } 8
	// Args:
	// * spec: Block containing parameter specifications
	// * context: Context object to use as parent context
	// * body: Block containing the function body code
	// Returns:
	// * function object with the specified parent context
	"fn\\par": {
		// a function with context	 bb: 10 add10 [ a ] context [ b: bb ] [ add a b ]
		// 							add10 [ a ] this [ add a b ]
		// later maybe			   add10 [ a ] [ b: b ] [ add a b ]
		//  						   add10 [ a ] [ 'b ] [ add a b ]
		Argsn: 3,
		Doc:   "Creates a function with a specified parent context, allowing access to variables from that context.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch args := arg0.(type) {
			case env.Block:
				ok, doc := util.ProcessFunctionSpec(args)
				if !ok {
					return MakeBuiltinError(ps, doc, "fn")
				}
				switch ctx := arg1.(type) {
				case env.RyeCtx:
					switch body := arg2.(type) {
					case env.Block:
						return *env.NewFunctionC(args, body, &ctx, false, false, doc)
					default:
						ps.ErrorFlag = true
						return MakeArgError(ps, 3, []env.Type{env.BlockType}, "fnc")
					}
				default:
					ps.ErrorFlag = true
					return MakeArgError(ps, 2, []env.Type{env.CtxType}, "fn\\par")
				}
			default:
				ps.ErrorFlag = true
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "fn\\par")
			}
		},
	},

	// Tests:
	// equal { ctx: context { y: 5 } , f: fn\in { x } ctx { x + y } , f 3 } 8
	// Args:
	// * spec: Block containing parameter specifications
	// * context: Context object to execute the function in
	// * body: Block containing the function body code
	// Returns:
	// * function object that executes directly in the specified context
	"fn\\in": {
		// a function with context	 bb: 10 add10 [ a ] context [ b: bb ] [ add a b ]
		// 							add10 [ a ] this [ add a b ]
		// later maybe			   add10 [ a ] [ b: b ] [ add a b ]
		//  						   add10 [ a ] [ 'b ] [ add a b ]
		Argsn: 3,
		Doc:   "Creates a function that executes directly in the specified context rather than creating a new execution context.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch args := arg0.(type) {
			case env.Block:
				ok, doc := util.ProcessFunctionSpec(args)
				if !ok {
					return MakeBuiltinError(ps, doc, "fn\\in")
				}
				switch ctx := arg1.(type) {
				case *env.RyeCtx:
					switch body := arg2.(type) {
					case env.Block:
						return *env.NewFunctionC(args, body, ctx, false, true, doc)
					default:
						ps.ErrorFlag = true
						return MakeArgError(ps, 3, []env.Type{env.BlockType}, "fn\\in")
					}
				case env.RyeCtx:
					switch body := arg2.(type) {
					case env.Block:
						return *env.NewFunctionC(args, body, &ctx, false, true, doc)
					default:
						ps.ErrorFlag = true
						return MakeArgError(ps, 3, []env.Type{env.BlockType}, "fn\\in")
					}
				default:
					ps.ErrorFlag = true
					return MakeArgError(ps, 2, []env.Type{env.CtxType}, "fn\\in")
				}
			default:
				ps.ErrorFlag = true
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "fn\\in")
			}
		},
	},

	// Tests:
	// equal { y: 5 , f: closure { x } { x + y } , f 3 } 8
	// Args:
	// * spec: Block containing parameter specifications
	// * body: Block containing the function body code
	// Returns:
	// * function object that captures the current context at creation time
	"closure": {
		// a function with context	 bb: 10 add10 [ a ] context [ b: bb ] [ add a b ]
		// 							add10 [ a ] this [ add a b ]
		// later maybe			   add10 [ a ] [ b: b ] [ add a b ]
		//  						   add10 [ a ] [ 'b ] [ add a b ]
		Argsn: 2,
		Doc:   "Creates a closure that captures the current context at creation time, preserving access to variables in that scope.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			ctx := *ps.Ctx
			switch args := arg0.(type) {
			case env.Block:
				ok, doc := util.ProcessFunctionSpec(args)
				if !ok {
					return MakeBuiltinError(ps, doc, "closure")
				}
				switch body := arg1.(type) {
				case env.Block:
					return *env.NewFunctionC(args, body, &ctx, false, false, doc)
				default:
					ps.ErrorFlag = true
					return MakeArgError(ps, 2, []env.Type{env.BlockType}, "closure")
				}
			default:
				ps.ErrorFlag = true
				return MakeArgError(ps, 1, []env.Type{env.BlockType}, "closure")
			}
		},
	},

	// Tests:
	// equal { prepend-star: partial ?concat [ "* " _ ] , prepend-star "hello" } "* hello"
	// equal { add-5: partial ?_+ [ _ 5 ] , add-5 10 } 15
	// equal { fn-add: fn { x y } { x + y } , add-5: partial ?fn-add [ _ 5 ] , add-5 10 } 15
	// Args:
	// * func: Function or builtin to partially apply
	// * args: Block of arguments, with _ (void) for arguments to be filled later
	// Returns:
	// * CurriedCaller object that can be called with the remaining arguments
	"partial": {
		Argsn: 2,
		Doc:   "Creates a partially applied function with specified arguments, using _ (void) for arguments to be filled later.",
		Pure:  true,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			// Check if first argument is a function or builtin
			var callerType int

			switch arg0.(type) {
			case env.Builtin:
				callerType = 0
			case env.Function:
				callerType = 1
			default:
				return MakeArgError(ps, 1, []env.Type{env.BuiltinType, env.FunctionType}, "partial")
			}

			// Check if second argument is a block
			var args []env.Object
			switch block := arg1.(type) {
			case env.Block:
				args = block.Series.GetAll()
				if len(args) > 5 {
					ps.FailureFlag = true
					return env.NewError("partial currently supports up to 5 arguments")
				}
			default:
				return MakeArgError(ps, 2, []env.Type{env.BlockType}, "partial")
			}

			// Extract arguments from the block
			var cur0, cur1, cur2, cur3, cur4 env.Object

			// Set arguments based on the block contents
			if len(args) > 0 {
				if args[0].Type() == env.VoidType {
					cur0 = nil
				} else {
					cur0 = args[0]
				}
			}

			if len(args) > 1 {
				if args[1].Type() == env.VoidType {
					cur1 = nil
				} else {
					cur1 = args[1]
				}
			}

			if len(args) > 2 {
				if args[2].Type() == env.VoidType {
					cur2 = nil
				} else {
					cur2 = args[2]
				}
			}

			if len(args) > 3 {
				if args[3].Type() == env.VoidType {
					cur3 = nil
				} else {
					cur3 = args[3]
				}
			}

			if len(args) > 4 {
				if args[4].Type() == env.VoidType {
					cur4 = nil
				} else {
					cur4 = args[4]
				}
			}

			// Create the CurriedCaller based on the function type
			if callerType == 0 {
				// Builtin
				return *env.NewCurriedCallerFromBuiltin(arg0.(env.Builtin), cur0, cur1, cur2, cur3, cur4)
			} else {
				// Function
				return *env.NewCurriedCallerFromFunction(arg0.(env.Function), cur0, cur1, cur2, cur3, cur4)
			}
		},
	},
}
