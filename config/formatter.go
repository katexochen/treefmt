package config

type Formatter struct {
	// Command is the command to invoke when applying this Formatter.
	Command string
	// Options are an optional list of args to be passed to Command.
	Options []string
	// Includes is a list of glob patterns used to determine whether this Formatter should be applied against a path.
	Includes []string
	// Excludes is an optional list of glob patterns used to exclude certain files from this Formatter.
	Excludes []string
	// Indicates the order of precedence when executing this Formatter in a sequence of Formatters.
	Priority int
}
