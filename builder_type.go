package manifest

type BuilderType string

const (
	BuilderTypeUnknown    BuilderType = ""
	BuilderTypeDocker     BuilderType = "docker"
	BuilderTypeKubernetes BuilderType = "kubernetes"
)
