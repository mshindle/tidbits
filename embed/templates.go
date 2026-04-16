package embed

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	layoutPattern = "layouts/*html"
	pageHome      = "index.html"
)

var errMissingTemplate = echo.NewHTTPError(http.StatusInternalServerError, "unable to render page")

// TemplateRenderer is a custom html/template renderer for Echo
type TemplateRenderer struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
	tmplFS    fs.FS
	funcMap   template.FuncMap
}

func NewTemplateRenderer(tmplFS fs.FS) *TemplateRenderer {
	funcMap := template.FuncMap{
		"fmtDate": func(t time.Time) string {
			return t.Format("Jan 2, 2006")
		},
		"fmtDuration": formatDuration,
		"slug": func(s string) string {
			s = strings.ToLower(s)
			s = strings.ReplaceAll(s, " ", "_")
			s = strings.ReplaceAll(s, ".", "_")
			return s
		},
		"lower": strings.ToLower,
		"pct":   func(v float32) float64 { return float64(v) * 100.0 },
	}

	tr := &TemplateRenderer{
		templates: make(map[string]*template.Template),
		tmplFS:    tmplFS,
		funcMap:   funcMap,
	}
	tr.Add(pageHome, layoutPattern, "pages/index.html")

	return tr
}

func (t *TemplateRenderer) Add(name string, patterns ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.templates[name] = template.Must(
		template.New(name).
			Funcs(t.funcMap).
			ParseFS(t.tmplFS, patterns...),
	)
}

func (t *TemplateRenderer) Exists(name string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	_, ok := t.templates[name]
	return ok
}

func (t *TemplateRenderer) Get(name string) (*template.Template, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tmpl, ok := t.templates[name]
	if !ok {
		return nil, errMissingTemplate
	}
	return tmpl, nil
}

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	log := Logger(c)
	log.WithField("template", name).WithField("data", data).Info("rendering template")

	tmpl, ok := t.templates[name]
	if !ok {
		log.WithField("template", name).Error("template not found")
		return errMissingTemplate
	}

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.WithError(err).Error("error rendering template")
	}
	return err
}

func formatDuration(d time.Duration) string {
	// If it's less than an hour, show minutes
	if d <= time.Hour {
		return fmt.Sprintf("%.0f mins", d.Minutes())
	}

	// 2. For an hour or more, round to the nearest 15 minutes
	rounded := d.Round(15 * time.Minute)

	h := int(rounded.Hours())
	m := int(rounded.Minutes()) % 60

	// 3. Format as H:MM hours
	// %d for hours, %02d ensures the minutes always have two digits (e.g., 8:00, 8:15)
	return fmt.Sprintf("%d:%02d hours", h, m)
}
