package test

type Command struct {
	Shell string

	Prebuilt string
	Params   map[string]interface{}
}
