package xrmoddir

type XRModDirContent struct {
	Data map[string]interface{}
}

func NewContent() *XRModDirContent {
	m := make(map[string]interface{})
	c := &XRModDirContent{
		Data: m,
	}
	c.Data["title"] = "X Rebirth Mod Directory"
	return c
}
