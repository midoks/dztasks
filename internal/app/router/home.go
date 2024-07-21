package router

import (
	"github.com/midoks/dztasks/internal/app/context"
)

const (
	HOME = "mail/list"
)

func Home(c *context.Context) {
	c.Data["PageIsHome"] = true
	c.Data["PageIsMail"] = true

	// mail.RenderMailSearch(c, &mail.MailSearchOptions{
	// 	PageSize: 10,
	// 	OrderBy:  "id ASC",
	// 	TplName:  HOME,
	// 	Type:     db.MailSearchOptionsTypeInbox,
	// })
}
