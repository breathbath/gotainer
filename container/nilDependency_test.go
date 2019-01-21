package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestNilDependencyForPointerProvided(t *testing.T) {
	cont := NewRuntimeContainer()
	cont.AddConstructor("nilDb", func(c Container) (i interface{}, e error) {
		return nil, nil
	})
	err := RegisterParameters(cont, map[string]interface{}{"nilCache": nil})
	if err != nil {
		t.Error(err)
	}

	cont.AddNewMethod("bookStorage", mocks.NewBookStorage, "nilDb")
	cont.AddNewMethod("cacheManager", mocks.NewCacheManager, "nilCache")

	cont.Get("bookStorage", true)
	cont.Get("cacheManager", true)
}
