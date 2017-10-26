package user

import (
	"sync"
	"bytes"
	"github.com/penguinn/penguin/component/session"
	"strconv"
	"github.com/penguinn/penguin/component/log"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/penguinn/album/models"
)

const (
	contextKey = "__user__"
	sessPrefix = "user#"
)

func generateKey(key string) string {
	b := bytes.Buffer{}
	b.WriteString(sessPrefix)
	b.WriteString(key)
	return b.String()
}

func FromContext(ctx *gin.Context) *User {
	var (
		ok   bool
		user *User
		data interface{}
	)
	if data, ok = ctx.Get(contextKey); ok {
		if user, ok = data.(*User); ok {
			return user
		}
	}
	user = NewUser(ctx)
	ctx.Set(contextKey, user)
	return user
}

type User struct {
	mu   sync.RWMutex
	id   string
	name string
	ctx  *gin.Context
}

func NewUser(ctx *gin.Context) *User {
	return &User{
		ctx: ctx,
	}
}

func (this *User) Login(id int, name string) error {
	sess, _ := session.SessionFromGin(this.ctx)
	sess.Set(generateKey("userID"), id)
	sess.Set(generateKey("userName"), name)
	sess.Set("user", id)

	session.BatchUpdateByUser(id, sess.Token(), map[string]interface{}{ "userID": id, generateKey("__kicked__"): true})

	updateMap := map[string]interface{} {
		"access_time" : time.Now().Unix(),
	}
	err := models.User{}.Update(id, updateMap)
	if err != nil {
		log.Error(err)
	}

	return err
}

func (this *User) IsGuest() bool {
	return this.GetID() == ""
}

func (this *User) IsKicked() bool {
	sess, ok := session.SessionFromGin(this.ctx)
	if ok {
		return sess.GetBool(generateKey("__kicked__"))
	}
	return false
}

func (this *User) GetID() string {
	sess, ok := session.SessionFromGin(this.ctx)
	if ok {
		return sess.GetString(generateKey("userID"))
	}
	return ""
}

func (this *User) GetName() string {
	sess, ok := session.SessionFromGin(this.ctx)
	if ok {
		return sess.GetString(generateKey("userName"))
	}
	return ""
}

func (this *User) Logout() {
	sess, ok := session.SessionFromGin(this.ctx)
	if ok {
		sess.Knockout()
	}
}
