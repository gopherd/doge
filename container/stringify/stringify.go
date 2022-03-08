package stringify

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopherd/doge/container"
	"github.com/gopherd/doge/operator"
)

// Options represents a options for stringify Node
type Options struct {
	Prefix     string
	Parent     string // default "│  "
	Space      string // default "   "
	Branch     string // default "├──"
	LastBranch string // default "└──"
}

var defaultOptions = &Options{
	Parent:     "│   ",
	Space:      "    ",
	Branch:     "├── ",
	LastBranch: "└── ",
}

func (options *Options) fix() {
	options.Parent = operator.Or(options.Parent, defaultOptions.Parent)
	if options.Branch == "" {
		options.Branch = defaultOptions.Branch
		options.LastBranch = operator.Or(options.LastBranch, defaultOptions.LastBranch)
	} else if options.LastBranch == "" {
		options.LastBranch = options.Branch
	}
	options.Space = operator.Or(options.Space, defaultOptions.Space)
}

// Stringify converts node to string
func Stringify[T comparable](node container.Node[T], options *Options) string {
	if options == nil {
		options = defaultOptions
	} else {
		options.fix()
	}
	if stringer, ok := node.(interface {
		Stringify(*Options) string
	}); ok {
		return stringer.Stringify(options)
	}
	var (
		buf   bytes.Buffer
		stack bytes.Buffer
	)
	if options.Prefix != "" {
		stack.WriteString(options.Prefix)
	}
	recursivelyPrintNode[T](node, &buf, &stack, "", 0, false, options)
	return buf.String()
}

func recursivelyPrintNode[T comparable](
	x interface{},
	w io.Writer,
	stack *bytes.Buffer,
	prefix string,
	depth int,
	isLast bool,
	options *Options,
) {
	var node, ok = x.(container.Node[T])
	if !ok {
		return
	}
	var nprefix = stack.Len()
	var value = node.String()
	var parent = node.Parent()
	fmt.Fprintf(w, "%s%s%s\n", stack.String(), prefix, value)
	var zero T

	if parent != zero {
		if isLast {
			stack.WriteString(options.Space)
		} else {
			stack.WriteString(options.Parent)
		}
	}
	var n = node.NumChild()
	for i := 0; i < n; i++ {
		isLast = i+1 == n
		var appended = operator.If(isLast, options.LastBranch, options.Branch)
		child := node.GetChildByIndex(i)
		if child == zero {
			continue
		}
		recursivelyPrintNode[T](child, w, stack, appended, depth+1, isLast, options)
	}

	if nprefix != stack.Len() {
		stack.Truncate(nprefix)
	}
}