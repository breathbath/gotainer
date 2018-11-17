package container

//CycleDetector creates Services at runtime with registered callbacks
type CycleDetector struct {
	recStack       map[string]bool
	recStackSorted []string
	visited        map[string]bool
	cycleDetected  bool
	cycle          [] string
	isEnabled      bool
}

func NewCycleDetector() *CycleDetector {
	cd := &CycleDetector{isEnabled:true}
	cd.Reset()
	return cd
}

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

func (cd *CycleDetector) VisitAfterRecursion(dep string) {
	if !cd.isEnabled {
		return
	}
	cd.recStack[dep] = false
}

func (cd *CycleDetector) DisableCycleDetection() {
	cd.isEnabled = false
}

func (cd *CycleDetector) EnableCycleDetection() {
	cd.isEnabled = true
}

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

func (cd *CycleDetector) GetCycle() []string {
	return cd.cycle
}

func (cd *CycleDetector) HasCycle() bool {
	return cd.cycleDetected
}
