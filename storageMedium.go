package manifest

type StorageMedium string

const (
	StorageMediumDefault   StorageMedium = ""          // use whatever the default is for the node, assume anything we don't explicitly handle is this
	StorageMediumMemory    StorageMedium = "Memory"    // use memory (e.g. tmpfs on linux)
	StorageMediumHugePages StorageMedium = "HugePages" // use hugepages
)
