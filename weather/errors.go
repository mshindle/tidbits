package weather

type constError string

const (
	ErrProviderFailure = constError("provider unable to fulfill request")
)

func (err constError) Error() string {
	return string(err)
}
