package form

import ()

type SignIn struct {
	Username string `binding:"Required;MaxSize(254)"`
	Password string `binding:"Required;MaxSize(255)"`
	Remember bool
}

type ArgsPluginInstall struct {
	Path string `binding:"Required;MaxSize(254)"`
}

type ArgsPluginUninstall struct {
	Path string `binding:"Required;MaxSize(254)"`
}

type ArgsPluginMenu struct {
	Name string `binding:"Required;MaxSize(254)"`
	Tag  string `binding:"Required;MaxSize(254)"`
}

type ArgsPluginData struct {
	Name   string `binding:"Required;MaxSize(254)"`
	Action string `binding:"Required;MaxSize(254)"`
	Page   int64
	Limit  int64
	Extra  string
	Args   map[string]interface{}
}
