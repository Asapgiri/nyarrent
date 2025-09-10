package logic

import (
	"nyarrent/dbase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (user *User) Map(duser dbase.User) {
    user.Id         = duser.Id.Hex()
    user.RegDate    = duser.RegDate
    user.EditDate   = duser.EditDate
    user.Username   = duser.Username
    user.Email      = duser.Email
    user.Roles      = duser.Roles
}

func (user *User) UnMap() dbase.User {
    duser := dbase.User{}

    duser.Id, _      = primitive.ObjectIDFromHex(user.Id)
    duser.RegDate    = user.RegDate
    duser.EditDate   = user.EditDate
    duser.Username   = user.Username
    duser.Email      = user.Email
    duser.Roles      = user.Roles

    return duser
}

func (user *User) Find(id primitive.ObjectID) {
    duser := dbase.User{}
    err := duser.Select(id)

    if nil != err {
        user.Id = ""
        return
    }

    user.Map(duser)
}

func (user *User) FindByUsername(username string) {
    duser := dbase.User{}
    err := duser.FindByUsername(username)

    if nil != err {
        user.Username = ""
        user.Id = ""
        return
    }

    user.Map(duser)
}
