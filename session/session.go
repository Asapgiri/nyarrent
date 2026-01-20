package session

import (
	"net/http"
	"nyarrent/logic"
	"time"

	"github.com/gorilla/sessions"
)

type Sessioner struct {
    Auth        logic.Auth
    Main        string
    Path        string
    Dto         any
    UserLocal   time.Location
    Meta        map[string]string
}
//FIXME: Handle fully separately in every function/session!!
//var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var store = sessions.NewCookieStore([]byte("slfjlksdjflsdfjlkfjs"))
var sessionName = "dunakeke"

func (session *Sessioner) Authenticate(r *http.Request) {
    // TODO: Add request aut header
    real_session, _ := store.Get(r, sessionName)
    uname, _ := real_session.Values[sessionName].(string)
	// if err != nil {
	// 	log.Println(err)
	// }

    session.Auth.Username = uname
    logic.Authenticate(&session.Auth)
}

func (session *Sessioner) New(w http.ResponseWriter, r *http.Request, uname string) {
    // FIXME: Store auth headers in database with associated user
    store.MaxAge(86400)
    rsess, _ := store.New(r, sessionName)

    rsess.Values[sessionName] = uname
    rsess.Save(r, w)
    session.Auth.Username = uname
}

func (session *Sessioner) Delete(w http.ResponseWriter, r *http.Request) {
    real_session, _ := store.Get(r, sessionName)
    real_session.Options.MaxAge = -1
    real_session.Save(r, w)
}

func (session *Sessioner) SetError(msg string) {
    session.Auth.Username = ""
    session.Auth.Error = msg
}

func GetCurrentSession(r *http.Request) Sessioner {
    session := Sessioner{}
    session.Authenticate(r)
    // FIXME: Get from user or set somewhere.
    session.UserLocal = *time.Local
    session.Meta["title"] = "Nyarrent"
    return session
}
