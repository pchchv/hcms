package handlers

import (
	"net/http"
	"strings"

	"github.com/pchchv/hcms/database"
	"golang.org/x/crypto/bcrypt"
)

// HandleLogin processes POST /admin/login.
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderer.Standalone(w, http.StatusBadRequest, "admin/login.html", map[string]any{
			"Title": "Вход в систему",
			"Error": "Ошибка обработки формы",
		})
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	settings, err := database.Get(s.db)
	if err != nil {
		s.renderer.Standalone(w, http.StatusInternalServerError, "admin/login.html", map[string]any{
			"Title": "Вход в систему",
			"Error": "Внутренняя ошибка сервера",
		})
		return
	}

	// case-insensitive email comparison
	if !strings.EqualFold(email, settings.AdminEmail) {
		s.renderer.Standalone(w, http.StatusUnauthorized, "admin/login.html", map[string]any{
			"Title": "Вход в систему",
			"Error": "Неверный email или пароль",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(settings.AdminPassword), []byte(password)); err != nil {
		s.renderer.Standalone(w, http.StatusUnauthorized, "admin/login.html", map[string]any{
			"Title": "Вход в систему",
			"Error": "Неверный email или пароль",
		})
		return
	}

	sessionID, err := database.CreateSession(s.db, settings.ID)
	if err != nil {
		s.renderer.Standalone(w, http.StatusInternalServerError, "admin/login.html", map[string]any{
			"Title": "Вход в систему",
			"Error": "Не удалось создать сессию",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "cms_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/admin/", http.StatusFound)
}

// HandleLoginPage renders the GET /admin/login page.
func (s *Server) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	s.renderer.Standalone(w, http.StatusOK, "admin/login.html", map[string]any{
		"Title": "Вход в систему",
	})
}

// HandleLogout processes POST /admin/logout.
func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cms_session")
	if err == nil && cookie.Value != "" {
		_ = database.DeleteSession(s.db, cookie.Value)
	}

	// clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "cms_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}
