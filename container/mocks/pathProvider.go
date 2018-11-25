package mocks

import "strings"

//PathProvider gives a path
type PathProvider interface {
	GetPath() string
}

//SimplePathProvider gust gives an
type SimplePathProvider struct {
	Path string
}

//GetPath gives internally saved path
func (pp SimplePathProvider) GetPath() string {
	return pp.Path
}

//PathsCollector contains all possible paths
type PathsCollector struct {
	Paths []string
}

//AddPath adder for new paths
func (pc *PathsCollector) AddPath(path string) {
	pc.Paths = append(pc.Paths, path)
}

//GetAllPaths getter for all paths
func (pc *PathsCollector) GetAllPaths() string {
	return strings.Join(pc.Paths, ",")
}

//LongestPathProvider finds the longest path
type LongestPathProvider struct {
	longestPath string
}

//GetPath implements the common interface
func (lpp *LongestPathProvider) GetPath() string {
	return lpp.longestPath
}

//EvaluatePath checks if the path is the longest one
func (lpp *LongestPathProvider) EvaluatePath(path string) {
	if len(path) > len(lpp.longestPath) {
		lpp.longestPath = path
	}
}
