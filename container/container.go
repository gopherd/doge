package container

// TreeFormatter contains options for formatting Tree
type TreeFormatter struct {
	Prefix string
	Debug  bool
	Color  bool

	IconParent     string // default "│  "
	IconSpace      string // default "   "
	IconBranch     string // default "├──"
	IconLastBranch string // default "└──"
}

func (formatter *TreeFormatter) Fix() {
	if formatter.IconParent == "" {
		formatter.IconParent = "│  "
	}
	if formatter.IconBranch == "" {
		formatter.IconBranch = "├──"
		if formatter.IconLastBranch == "" {
			formatter.IconLastBranch = "└──"
		}
	} else if formatter.IconLastBranch == "" {
		formatter.IconLastBranch = formatter.IconBranch
	}
	if formatter.IconSpace == "" {
		formatter.IconSpace = "   "
	}
	// append spaces
	formatter.IconParent += " "
	formatter.IconBranch += " "
	formatter.IconLastBranch += " "
	formatter.IconSpace += " "
}
