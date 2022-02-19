package container

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopherd/doge/operator"
)

// Pair contains two values like c++ std::pair
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// MakePair make a pair by values
func MakePair[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{First: first, Second: second}
}

// Node represents a generic printable node
type Node[T comparable] interface {
	ToString() string        // ToString returns node self information
	Parent() T               // Parent returns parent node or nil
	NumChild() int           // NumChild returns number of child
	GetChildByIndex(i int) T // GetChildByIndex gets child by index
}

// PrintOptions represents a options for printing Node
type PrintOptions struct {
	Prefix         string
	IconParent     string // default "│  "
	IconSpace      string // default "   "
	IconBranch     string // default "├──"
	IconLastBranch string // default "└──"
}

func (options *PrintOptions) Fix() {
	options.IconParent = operator.Or(options.IconParent, "│  ")
	if options.IconBranch == "" {
		options.IconBranch = "├──"
		options.IconLastBranch = operator.Or(options.IconLastBranch, "└──")
	} else if options.IconLastBranch == "" {
		options.IconLastBranch = options.IconBranch
	}
	options.IconSpace = operator.Or(options.IconSpace, "   ")
	// append spaces
	options.IconParent += " "
	options.IconBranch += " "
	options.IconLastBranch += " "
	options.IconSpace += " "
}

// PrintNode recursively prints node
func PrintNode[T comparable](node Node[T], options PrintOptions) string {
	options.Fix()
	if printer, ok := node.(interface{ Print(PrintOptions) string }); ok {
		return printer.Print(options)
	}
	var (
		buf         bytes.Buffer
		prefixstack bytes.Buffer
	)
	if options.Prefix != "" {
		prefixstack.WriteString(options.Prefix)
	}
	recursivelyPrintNode[T](node, &buf, &prefixstack, "", 0, false, options)
	return buf.String()
}

func recursivelyPrintNode[T comparable](
	x interface{},
	w io.Writer,
	pstack *bytes.Buffer,
	prefix string,
	depth int,
	isLast bool,
	options PrintOptions,
) {
	var node, ok = x.(Node[T])
	if !ok {
		return
	}
	var nprefix = pstack.Len()
	var value = node.ToString()
	var parent = node.Parent()
	fmt.Fprintf(w, "%s%s%s\n", pstack.String(), prefix, value)
	var zero T

	if parent != zero {
		if isLast {
			pstack.WriteString(options.IconSpace)
		} else {
			pstack.WriteString(options.IconParent)
		}
	}
	var n = node.NumChild()
	for i := 0; i < n; i++ {
		isLast = i+1 == n
		var appended = operator.If(isLast, options.IconLastBranch, options.IconBranch)
		child := node.GetChildByIndex(i)
		if child == zero {
			continue
		}
		recursivelyPrintNode[T](child, w, pstack, appended, depth+1, isLast, options)
	}

	if nprefix != pstack.Len() {
		pstack.Truncate(nprefix)
	}
}
