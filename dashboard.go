package guerrilla

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	sc "github.com/gorilla/securecookie"
)

const (
	dashboard      = "index.html"
	login          = "login.html"
	dashboardPath  = "dashboard/html/index.html"
	loginPath      = "dashboard/html/index.html"
	cookieName     = "guerrilla_dashboard"
	idHeader       = "X-Guerrilla-ID"
	sessionTimeout = time.Hour * 24
)

var (
	// Cache of HTML templates
	templates = template.Must(template.ParseFiles(dashboardPath, loginPath))
	// Analytics configuration
	config   *AnalyticsConfig = nil
	sessions sessionStore
	nextID   = 1
)

type Session struct {
	Start, Expires time.Time
	ID             int
	SecureCookie   *sc.SecureCookie
}

type sessionStore map[int]*Session

func (ss sessionStore) Clean() {
	now := time.Now()
	for id, sess := range ss {
		if sess.Expires.Before(now) {
			delete(ss, id)
		}
	}
}

func (ss sessionStore) cleaner(ticker *time.Ticker) {
	for {
		<-ticker.C
		ss.Clean()
	}
}

func Run(ac *AnalyticsConfig /*, ds *AnalyticsDataStore*/) {
	config = ac
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/login", loginHandler)

	sessions := make(sessionStore)
	tick := time.NewTicker(sessionTimeout)
	go sessions.cleaner(tick)

	http.ListenAndServe(ac.ListenInterface, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	if isLoggedIn(r) {
		w.WriteHeader(http.StatusAccepted)
		templates.ExecuteTemplate(w, dashboard, nil)
	} else {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	switch r.Method {
	case "GET":
		// Print login screen
		templates.ExecuteTemplate(w, login, nil)

	case "POST":
		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user == config.WebUsername && pass == config.WebPassword {
			err := startSession(w)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				// TODO Internal error
			}
			http.Redirect(w, r, "/", http.StatusAccepted)
		} else {
			templates.ExecuteTemplate(w, login, nil) // TODO info about failed login
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func startSession(w http.ResponseWriter) error {
	key := sc.GenerateRandomKey(64)
	s := sc.New(key, nil)
	contents := map[string]string{
		"connected": time.Now().Format(time.RFC3339),
		"expires":   time.Now().Add(time.Hour * 24).Format(time.RFC3339),
		"id":        strconv.Itoa(nextID),
	}

	encoded, err := s.Encode(cookieName, contents)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  "guerrilla_dashboard",
		Value: encoded,
		Path:  "/",
		// Secure: true,
	}

	sess := &Session{
		Start:        time.Now(),
		Expires:      time.Now().Add(sessionTimeout), // TODO config for this
		SecureCookie: s,
	}

	http.SetCookie(w, cookie)
	w.Header().Set(idHeader, contents["id"])
	sessions[nextID] = sess
	nextID++

	return nil
}

func isLoggedIn(r *http.Request) bool {
	id, err := strconv.Atoi(r.Header.Get(idHeader))
	if err != nil {
		return false
	}

	sess, ok := sessions[id]
	if !ok || sess == nil {
		return false
	}

	c, err := r.Cookie(cookieName)
	if err != nil || c == nil {
		return false
	}

	if sess.Expires.After(time.Now()) {
		return false
	}

	contents := make(map[string]string)
	err = sess.SecureCookie.Decode(cookieName, c.Value, &contents)
	if err != nil {
		return false
	}

	sid, _ := strconv.Atoi(contents["id"])
	if sid != id {
		return false
	}

	return true
}
