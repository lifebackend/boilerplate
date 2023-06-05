package handlers

func NewHandlers(
	imageTag string,
) *Handlers {
	return &Handlers{
		imageTag: imageTag,
	}
}

type Handlers struct {
	imageTag string
}
