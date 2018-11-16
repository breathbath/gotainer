package container

//CycleDetector creates Services at runtime with registered callbacks
type CycleDetector struct {
	recStack       map[string]bool
	recStackSorted []string
	visited        map[string]bool
	cycleDetected  bool
	cycle          [] string
}

func NewCycleDetector() *CycleDetector {
	cd := &CycleDetector{}
	cd.Reset()
	return cd
}

func (cd *CycleDetector) Reset() {
	cd.cycle = []string{}
	cd.cycleDetected = false
	cd.recStack = make(map[string]bool)
	cd.visited = make(map[string]bool)
}

func (cd *CycleDetector) VisitBeforeRecursion(dep string) bool {
	if cd.cycleDetected {
		return true
	}

	if isInRecStack, ok := cd.recStack[dep]; ok && isInRecStack {
		cd.registerCycle(dep)
		return true
	}

	if isInVisited, ok := cd.visited[dep]; ok && isInVisited {
		return false
	}

	cd.visited[dep] = true
	cd.recStack[dep] = true
	cd.recStackSorted = append(cd.recStackSorted, dep)

	return false
}

func (cd *CycleDetector) VisitAfterRecursion(dep string) {
	cd.recStack[dep] = false
}

func (cd *CycleDetector) registerCycle(dep string) {
	cd.recStackSorted = append(cd.recStackSorted, dep)
	if len(cd.recStackSorted) == 1 {
		cd.recStackSorted = append(cd.recStackSorted, cd.recStackSorted[0])
	}
	for _, cyclicDep := range cd.recStackSorted {
		if isTrue, ok := cd.recStack[cyclicDep]; ok && isTrue {
			cd.cycle = append(cd.cycle, cyclicDep)
		}
	}
	cd.cycleDetected = true
}

func (cd *CycleDetector) GetCycle() []string {
	return cd.cycle
}

func (cd *CycleDetector) HasCycle(dep string) bool {
	return cd.cycleDetected
}
