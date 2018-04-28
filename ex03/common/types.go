package common

type Gender int
const (
	GENDER_NULL   Gender   = 0
	GENDER_MALE   Gender   = 1
	GENDER_FEMALE Gender   = 2
)

type OnlineState int
const (
	ONLINESTATE_UNKNOWN OnlineState = 0
	ONLINESTATE_ON		OnlineState = 1
	ONLINESTATE_OFF		OnlineState = 2
)

type AssetsType int
const (
	ASSETSTYPE_UNKNOWN AssetsType = 0
)


type BindingType int
const (
	BINDINGTYPE_UNKNOWN BindingType = 0
)
