package resource_manage

type ResourceManage interface {
	GetOne()
	FreeOne()
	Has() uint
	Left() uint
}
