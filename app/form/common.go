package form

import ()

type SignIn struct {
	Username string `binding:"Required;MaxSize(254)"`
	Password string `binding:"Required;MaxSize(255)"`
	Remember bool
}

type PluginInstall struct {
	Path string `binding:"Required;MaxSize(254)"`
}

type PluginUninstall struct {
	Path string `binding:"Required;MaxSize(254)"`
}