package types

type ChangelogSection struct {
	Version  string
	Date     string
	Features []string
	Fixes    []string
	Others   []string
	LastCommitHash string
	Hash     string
	Append   bool
}