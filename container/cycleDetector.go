package container

//CycleDetector used for identification of possible dependency cycles
type CycleDetector struct {
	recStack       map[string]bool
	recStackSorted []string
	visited        map[string]bool
	cycleDetected  bool
	cycle          []string
	isEnabled      bool
}

//NewCycleDetector constructor
func NewCycleDetector() *CycleDetector {
	cd := &CycleDetector{isEnabled: true}
	cd.Reset()
	return cd
}

//Reset removes all info about previous cycle detection
func (cd *CycleDetector) Reset() {
	if !cd.isEnabled {
		return
	}
	cd.cycle = []string{}
	cd.cycleDetected = false
	cd.recStack = make(map[string]bool)
	cd.visited = make(map[string]bool)
	cd.recStackSorted = []string{}
	cd.isEnabled = true
}

//VisitBeforeRecursion starts DFS
func (cd *CycleDetector) VisitBeforeRecursion(dep string) {
	if !cd.isEnabled || cd.cycleDetected {
		return
	}

	if isInRecStack, ok := cd.recStack[dep]; ok && isInRecStack {
		cd.registerCycle(dep)
		return
	}

	if isInVisited, ok := cd.visited[dep]; ok && isInVisited {
		return
	}

	cd.visited[dep] = true
	cd.recStack[dep] = true
	cd.recStackSorted = append(cd.recStackSorted, dep)
}

//VisitAfterRecursion ends current DFS
func (cd *CycleDetector) VisitAfterRecursion(dep string) {
	if !cd.isEnabled {
		return
	}
	cd.recStack[dep] = false
}

//DisableCycleDetection can be used to stop the cycle detection at runtime
func (cd *CycleDetector) DisableCycleDetection() {
	cd.isEnabled = false
}

//EnableCycleDetection can be used to start the cycle detection at runtime
func (cd *CycleDetector) EnableCycleDetection() {
	cd.isEnabled = true
}

//IsEnabled checks if cycle detection is enabled
func (cd *CycleDetector) IsEnabled() bool {
	return cd.isEnabled
}

func (cd *CycleDetector) registerCycle(dep string) {
	cd.recStackSorted = append(cd.recStackSorted, dep)
	for _, cyclicDep := range cd.recStackSorted {
		if isTrue, ok := cd.recStack[cyclicDep]; ok && isTrue {
			cd.cycle = append(cd.cycle, cyclicDep)
		}
	}
	cd.cycleDetected = true
}

//GetCycle get detected cycle info
func (cd *CycleDetector) GetCycle() []string {
	return cd.cycle
}

//HasCycle gives info if a cycle exists
func (cd *CycleDetector) HasCycle() bool {
	return cd.cycleDetected
}
