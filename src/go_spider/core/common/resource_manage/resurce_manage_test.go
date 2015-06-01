package resource_manage

import (
	"fmt"
	"testing"
)

func TestResourceManage(t *testing.T) {
	var mc *ResourceManageChan
	mc = NewResourceManageChan(1)
	mc.GetOne()
	fmt.Println("increase")
	mc.FreeOne()
	fmt.Println("decrease")
	mc.GetOne()
	fmt.Println("increase")
}
