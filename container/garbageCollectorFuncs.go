package container

type namedGarbageCollectorFunc struct {
	f    GarbageCollectorFunc
	name string
}

type GarbageCollectorFuncs struct {
	garbageCollectors []namedGarbageCollectorFunc
	namedMap          map[string]bool
}

func NewGarbageCollectorFuncs() *GarbageCollectorFuncs {
	return &GarbageCollectorFuncs{
		garbageCollectors: []namedGarbageCollectorFunc{},
		namedMap:          make(map[string]bool),
	}
}

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

func (gcf *GarbageCollectorFuncs) Range(iterFunc func(gcName string, f GarbageCollectorFunc) bool) {
	for _, namedGcFunc := range gcf.garbageCollectors {
		result := iterFunc(namedGcFunc.name, namedGcFunc.f)
		if !result {
			break
		}
	}
}
