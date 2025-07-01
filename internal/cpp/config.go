package cpp

type Config struct {
	RecursiveIncludePaths    []string `yaml:"recursiveIncludePaths"`
	NonRecursiveIncludePaths []string `yaml:"nonRecursiveIncludePaths"`
}
