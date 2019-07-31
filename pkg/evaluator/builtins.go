package evaluator

import "go.smartmachine.io/cumulus/pkg/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: lenBuiltin,
	},
}

func lenBuiltin(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s, want STRING", arg.Type())
	}
}
