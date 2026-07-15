package frontend

// "io"
// "net/http"
// "net/http/httptest"
// "testing"

// "github.com/labstack/echo/v5"

type renderer struct {
	name string
	data map[string]any
}

// func (r *renderer) Render(_ *echo.Context, _ io.Writer, name string, data any) error {
// 	r.name = name
// 	r.data, _ = data.(map[string]any)
// 	return nil
// }
// func TestGetIndexRendersFrontendTemplate(t *testing.T) {
// 	e := echo.New()
// 	r := &renderer{}
// 	e.Renderer = r
// 	&FrontendHandler{}
// 	rec := httptest.NewRecorder()
// 	if err := GetIndex(e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)); err != nil || r.name != "frontend:index/index.html" || r.data["Name"] != "Yien" {
// 		t.Fatalf("name=%q data=%#v err=%v", r.name, r.data, err)
// 	}
// }
