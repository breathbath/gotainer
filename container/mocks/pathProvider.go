package mocks

import "strings"

type PathProvider interface {
	GetPath() string
}

type SimplePathProvider struct {
	Path string
}

func (pp SimplePathProvider) GetPath() string {
	return pp.Path
}

type PathsCollector struct {
	Paths []string
}

func (pc *PathsCollector) AddPath(path string) {
	pc.Paths = append(pc.Paths, path)
}

func (pc *PathsCollector) GetAllPaths() string {
	return strings.Join(pc.Paths, ",")
}

type LongestPathProvider struct {
	longestPath string
}

func (lpp *LongestPathProvider) GetPath() string {
	return lpp.longestPath
}

func (lpp *LongestPathProvider) EvaluatePath(path string) {
	if len(path) > len(lpp.longestPath) {
		lpp.longestPath = path
	}
}
