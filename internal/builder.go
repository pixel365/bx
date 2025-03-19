package internal

type Builder interface {
	Build() error
	Prepare() error
	Rollback() error
	Collect() error
	Cleanup()
}
