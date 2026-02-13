// Package web
package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gin-alpine/src/internal/domain/auth"
	"gin-alpine/src/pkg/utils"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

var baseLayoutComponents = []string{
	"templates/partials/sidebar.html",
	"templates/partials/footer.html",
	"templates/partials/navbar.html",
}

//go:embed static/*
var StaticFiles embed.FS

//go:embed templates/**/*.html
var TemplateFS embed.FS

type Renderer struct {
	pages map[string]*template.Template
	mode  string
	fs    fs.FS
}

func NewRenderer(mode string) *Renderer {
	r := &Renderer{
		pages: make(map[string]*template.Template),
		mode:  mode,
	}
	path, err := utils.GetFilePath([]string{"src", "services", "web"})
	if err != nil {
		utils.FatalResult("error getting abs template path", err)
	}
	if mode == gin.ReleaseMode {
		r.fs = TemplateFS
	} else {
		r.fs = os.DirFS(path)
	}
	// no layouts
	r.pages["404"] = r.page("base", "404")
	r.pages["500"] = r.page("base", "500")
	r.pages["login"] = r.page("auth", "login")

	// base layouts
	r.pages["home"] = r.page("main", "home", r.setExtra()...)

	// extra layouts
	// r.pages["padroes"] = r.page("main", "padroes",
	// 	r.setExtra("templates/icons/error-logo.html")...,
	// )

	return r
}
func (r *Renderer) setExtra(extra ...string) []string {
	return append(baseLayoutComponents, extra...)
}

func (r *Renderer) page(layout, page string, extra ...string) *template.Template {
	files := []string{
		"templates/layouts/" + layout + ".html",
		"templates/pages/" + page + ".html",
	}
	files = append(files, extra...)

	return r.mustParse(files...)
}

func (r *Renderer) mustParse(files ...string) *template.Template {
	t := template.New("")
	t, err := t.ParseFS(r.fs, files...)
	if err != nil {
		panic(err)
	}
	return t
}

func (r *Renderer) Render(
	w http.ResponseWriter,
	layout,
	page string,
	data any,
) error {
	tmpl, ok := r.pages[page]
	if !ok {
		return fmt.Errorf("template not found: %s", page)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
		return err
	}

	return nil
}

// Page is indicate to handle simple 200 GET pages, expecting only possible
// template errors, to use in a more imperative way, in handlers for example,
// should use Render method directly.
func (r *Renderer) Page(
	layout string,
	page string,
	dataFn func(*gin.Context) gin.H,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := gin.H{
			"csrf": csrf.GetToken(c),
		}

		// always inject auth context
		for k, v := range ExtractAuthContext(c) {
			data[k] = v
		}

		if dataFn != nil {
			for k, v := range dataFn(c) {
				data[k] = v
			}
		}

		if err := r.Render(
			c.Writer,
			layout,
			page,
			data,
		); err != nil {
			// Centralized error handling
			log.Printf("error rendering page: %v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

// ExtractAuthContext example usage:
// use direct in the templates like this:
// {{ if .IsAuth }}
//
//	<p>Welcome, {{ .User.Name }}</p>
//
// {{ end }}
//
// {{ if .Can.Admin }}
//
//	<a href="/admin">Admin Panel</a>
//
// {{ end }}
//
// {{ if not .IsAuth }}
//
//	<a href="/login">Login</a>
//
// {{ end }}
func ExtractAuthContext(c *gin.Context) gin.H {
	v, exists := c.Get("auth_user")
	appDataCtx, existsAppData := c.Get("app_data")
	if !exists || !existsAppData {
		return gin.H{
			"AppData": nil,
			"User":    nil,
			"Can": gin.H{
				"Customer": false,
				"Manager":  false,
				"Admin":    false,
				"Dev":      false,
			},
			"IsAuth": false,
		}
	}

	user := v.(auth.UserAuth)
	appData := appDataCtx.(utils.AppData)

	return gin.H{
		"AppData": appData,
		"User":    user,
		"IsAuth":  true,
		"Can": gin.H{
			"Customer": user.Role >= auth.RoleCustomer,
			"Manager":  user.Role >= auth.RoleManager,
			"Admin":    user.Role >= auth.RoleAdmin,
			"Dev":      user.Role >= auth.RoleDev,
		},
	}
}

// func debugFS(fsys fs.FS) {
// 	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			log.Println("ERR:", err)
// 			return nil
// 		}
// 		log.Println("FS:", path)
// 		return nil
// 	})
// }
