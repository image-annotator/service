package app

import (
	"fmt"
	"net/http"
	"strings"
)

func (rs *UserResource) basicAuthFactory(userrole []string) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			if len(auth) != 2 || auth[0] != "Basic" {
				http.Error(w, "authorization failed", http.StatusUnauthorized)
				return
			}

			if !validate(auth[1], rs, userrole) {
				http.Error(w, "authorization failed", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
	return
}

func validate(cookie string, rs *UserResource, userrole []string) bool {

	user, err := rs.Store.GetByCookie(cookie)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if stringInSlice(user.UserRole, userrole) {
		return true
	}
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
