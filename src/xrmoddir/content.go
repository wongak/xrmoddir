package xrmoddir

type XRModDirContent struct {
	Data map[string]interface{}
}

func NewContent() *XRModDirContent {
	m := make(map[string]interface{})
	return &XRModDirContent{
		Data: m,
	}
}

func (x XRModDirContent) Get() map[string]interface{} {
	return x.Data
}
