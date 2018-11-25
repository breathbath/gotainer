package container

type namedGarbageCollectorFunc struct {
	f    GarbageCollectorFunc
	name string
}

//GarbageCollectorFuncs used to hold info about all garbage collectors
type GarbageCollectorFuncs struct {
	garbageCollectors []namedGarbageCollectorFunc
	namedMap          map[string]bool
}

//NewGarbageCollectorFuncs constructor
func NewGarbageCollectorFuncs() *GarbageCollectorFuncs {
	return &GarbageCollectorFuncs{
		garbageCollectors: []namedGarbageCollectorFunc{},
		namedMap:          make(map[string]bool),
	}
}

//Add a new garbage collector func
func (gcf *GarbageCollectorFuncs) Add(name string, gcFunc GarbageCollectorFunc) {
	if _, exists := gcf.namedMap[name]; exists {
		return
	}

	namedGcFunc := namedGarbageCollectorFunc{
		f:    gcFunc,
		name: name,
	}

	gcf.garbageCollectors = append(gcf.garbageCollectors, namedGcFunc)
	gcf.namedMap[name] = true
}

//Range iterates over garbage collectors
func (gcf *GarbageCollectorFuncs) Range(iterFunc func(gcName string, f GarbageCollectorFunc) bool) {
	for _, namedGcFunc := range gcf.garbageCollectors {
		result := iterFunc(namedGcFunc.name, namedGcFunc.f)
		if !result {
			break
		}
	}
}
