package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yourusername/go-mockserver/internal/config"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, s.buildUI())
			return
		}
		for _, route := range s.cfg.Routes {
			if route.Path == r.URL.Path && strings.EqualFold(route.Method, r.Method) {
				s.makeHandler(route)(w, r)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "route not found"}`, http.StatusNotFound)
	})

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	log.Printf("\nサーバーを起動しました → http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) makeHandler(route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if route.Response.Delay > 0 {
			time.Sleep(time.Duration(route.Response.Delay) * time.Millisecond)
		}
		for key, value := range route.Response.Headers {
			w.Header().Set(key, value)
		}
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(route.Response.Status)
		fmt.Fprint(w, route.Response.Body)
		log.Printf("  → %d  %s %s", route.Response.Status, r.Method, r.URL.Path)
	}
}

type routeJSON struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Status int    `json:"status"`
	Body   string `json:"body"`
	Delay  int    `json:"delay"`
}

func (s *Server) buildUI() string {
	var routes []routeJSON
	for _, r := range s.cfg.Routes {
		routes = append(routes, routeJSON{
			Method: r.Method,
			Path:   r.Path,
			Status: r.Response.Status,
			Body:   r.Response.Body,
			Delay:  r.Response.Delay,
		})
	}
	routesJSON, _ := json.Marshal(routes)

	html := `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>go-mockserver</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: system-ui, sans-serif; background: #f5f5f4; color: #1c1917; min-height: 100vh; }
.header { background: white; border-bottom: 1px solid #e7e5e4; padding: 1rem 1.5rem; display: flex; align-items: center; gap: 10px; }
.dot { width: 10px; height: 10px; border-radius: 50%; background: #4ade80; }
.title { font-size: 18px; font-weight: 600; }
.port { margin-left: auto; font-size: 13px; color: #78716c; background: #f5f5f4; padding: 4px 10px; border-radius: 6px; }
.container { max-width: 720px; margin: 0 auto; padding: 1.5rem; }
.card { background: white; border: 1px solid #e7e5e4; border-radius: 12px; padding: 1rem 1.25rem; margin-bottom: 10px; }
.route-header { display: flex; align-items: center; gap: 10px; }
.method { font-size: 12px; font-weight: 600; padding: 3px 9px; border-radius: 5px; }
.GET    { background: #dbeafe; color: #1e40af; }
.POST   { background: #dcfce7; color: #166534; }
.PUT    { background: #fef9c3; color: #854d0e; }
.DELETE { background: #fee2e2; color: #991b1b; }
.path { font-family: monospace; font-size: 14px; font-weight: 500; }
.btn { margin-left: auto; font-size: 13px; padding: 5px 14px; border: 1px solid #d4d4d4; border-radius: 7px; background: white; cursor: pointer; }
.btn:hover { background: #f5f5f4; }
.res { margin-top: 10px; background: #fafaf9; border-radius: 8px; padding: 10px 12px; font-family: monospace; font-size: 13px; display: none; white-space: pre-wrap; word-break: break-all; }
.res.show { display: block; }
.badge { display: inline-block; font-size: 11px; padding: 2px 8px; border-radius: 5px; margin-bottom: 6px; }
.ok  { background: #dcfce7; color: #166534; }
.err { background: #fee2e2; color: #991b1b; }
</style>
</head>
<body>
<div class="header">
  <div class="dot"></div>
  <span class="title">go-mockserver</span>
  <span class="port">localhost:` + fmt.Sprintf("%d", s.cfg.Port) + `</span>
</div>
<div class="container" id="app"></div>
<script>
const routes = ` + string(routesJSON) + `;
const app = document.getElementById('app');
routes.forEach(function(r, i) {
  const card = document.createElement('div');
  card.className = 'card';
  card.innerHTML =
    '<div class="route-header">' +
      '<span class="method ' + r.method + '">' + r.method + '</span>' +
      '<span class="path">' + r.path + '</span>' +
      '<button class="btn" onclick="send(' + i + ')">送信 →</button>' +
    '</div>' +
    '<div class="res" id="res' + i + '"></div>';
  app.appendChild(card);
});
function send(i) {
  var r = routes[i];
  var box = document.getElementById('res' + i);
  box.className = 'res show';
  box.textContent = '送信中...';
  fetch(r.path, { method: r.method })
    .then(function(res) {
      var status = res.status;
      return res.text().then(function(body) {
        var ok = status < 400;
        box.innerHTML = '<span class="badge ' + (ok ? 'ok' : 'err') + '">' + status + '</span>\n' + body;
      });
    })
    .catch(function(e) { box.textContent = 'エラー: ' + e.message; });
}
</script>
</body>
</html>`
	return html
}