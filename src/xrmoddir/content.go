package xrmoddir

type XRModDirContent struct {
	Content string
	Data    map[string]interface{}
}

func NewContent() *XRModDirContent {
	m := make(map[string]interface{})
	c := &XRModDirContent{
		Data: m,
	}
	c.Data["title"] = "X Rebirth Mod Directory"
	return c
}
