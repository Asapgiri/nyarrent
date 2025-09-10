package pages

import (
	"net/http"
	"nyarrent/dbase"
	"nyarrent/logic"
	"nyarrent/session"
	"slices"
)

func checkForAccess(session session.Sessioner, requested_roles []string) bool {
	if "" == session.Auth.Username {
		return false
	}

	for _, reqrole := range(requested_roles) {
		if slices.Contains(session.Auth.Roles, reqrole) {
			return true
		}
	}

	return len(requested_roles) == 0
}

func Login(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

	if "" != session.Auth.Username {
        http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

    // TODO: login with phone number...?
    uname := r.FormValue("form[userNameOrEmail]")
    upass := r.FormValue("form[userPass]")

    if "" != uname {
        user := logic.User{}
        err := user.Login(uname, upass)
        if nil != err {
            session.SetError(err.Error())
        } else {
            session.Delete(w, r)
            session.New(w, r, user.Username)
        }
		log.Println(user)
    } else {
        session.SetError("")
    }

	log.Println(session.Auth)
    if "" == session.Auth.Username {
        fil, _ := read_artifact("auth/login.html", w.Header())
        Render(session, w, fil, nil)
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Register(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

	if "" != session.Auth.Username {
        http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sha := dbase.RegisterSha{}
	sha.Find(r.PathValue("sha"))
	if "" == sha.Sha {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

    uuname := r.FormValue("form[userUsername]")
    uemail := r.FormValue("form[userEmail]")
    upassa := r.FormValue("form[userPassA]")
    upassb := r.FormValue("form[userPassB]")

    if "" != uuname {
        user := logic.User{
            Username: uuname,
            Email: uemail,
        }
        err := user.Register(upassa, upassb)
        if nil != err {
            session.SetError(err.Error())
            log.Println(err.Error())
        } else {
            session.Delete(w, r)
            session.New(w, r, user.Username)
        }
    } else {
        session.SetError("")
    }

    if "" == session.Auth.Username {
        fil, _ := read_artifact("auth/register.html", w.Header())
        Render(session, w, fil, nil)
    } else {
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    session.Delete(w, r)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
