package user

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/penguinn/album/models"
	"github.com/penguinn/penguin/component/log"
	"github.com/penguinn/penguin/component/session"
	"sync"
	"time"
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
	sess.Set("userID", id)

	session.BatchUpdateByUser(id, sess.Token(), map[string]interface{}{"userID": id, generateKey("__kicked__"): true})

	updateMap := map[string]interface{}{
		"access_time": time.Now().Unix(),
	}
	err := models.User{}.Update(id, updateMap)
	if err != nil {
		log.Error(err)
	}

	return err
}

func (this *User) IsGuest() bool {
	return this.GetID() == 0
}

func (this *User) IsKicked() bool {
	sess, ok := session.SessionFromGin(this.ctx)
	if ok {
		return sess.GetBool(generateKey("__kicked__"))
	}
	return false
}

func (this *User) GetID() int {
	sess, ok := session.SessionFromGin(this.ctx)
	if ok {
		return sess.GetInt(generateKey("userID"))
	}
	return 0
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
