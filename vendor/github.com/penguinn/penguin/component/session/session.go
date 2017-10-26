package session

import (
	"encoding/json"
	"fmt"
	"net/http"

	"errors"

	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	cookieName        = "sess-token"
	contextSessionKey = "_session"
	sessionTableName  = "session_store"
)

var (
	rander = rand.New(rand.NewSource(time.Now().Unix()))

	codec Codec
	store Store

	defaultMaxAge = 3600 * 24 * 30 * 6

	codecNotFoundErr    = errors.New("Codec Not Found")
	storeNotFoundErr    = errors.New("Store Not Found")
	dbGetterNotFoundErr = errors.New("dbGetter Not Found")

	requireChecker = func() error {
		if codec == nil {
			return codecNotFoundErr
		}
		if store == nil {
			return storeNotFoundErr
		}
		return nil
	}

	newToken = func() string {
		r := [32]byte{}
		for i := 0; i < 32; i++ {
			r[i] = byte(rander.Int31n(16))
		}
		return fmt.Sprintf("%X%X%X%X%X%X%X%X-%X%X%X%X-%X%X%X%X-%X%X%X%X-%X%X%X%X%X%X%X%X%X%X%X%X",
			r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7],
			r[8], r[9], r[10], r[11],
			r[12], r[13], r[14], r[15],
			r[16], r[17], r[18], r[19],
			r[20], r[21], r[22], r[23], r[24], r[25], r[26], r[27], r[28], r[29], r[30], r[31])
	}
)

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type jsonCodec struct{}

func NewJsonCodec() *jsonCodec {
	return new(jsonCodec)
}

func (*jsonCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (*jsonCodec) Decode(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}

type Store interface {
	Get(string) (*StoreData, bool)
	Save(StoreData) error
	Del(string) error
	Each(func(StoreData))
	CleanUp()
	Users(int) ([]StoreData, error)
	BatchUpdateByUser(int, string, string) error
}

type dbStore struct {
	DBGetter func() *gorm.DB
}

func NewDBStore(f func() *gorm.DB) (*dbStore, error) {
	if f == nil {
		return nil, dbGetterNotFoundErr
	}
	if !f().HasTable(sessionTableName) {
		f().Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&StoreData{})
	}
	return &dbStore{DBGetter: f}, nil
}

func (d *dbStore) Get(t string) (*StoreData, bool) {
	s := &StoreData{}

	if err := d.DBGetter().Where("token = ? AND expire_at > ?", t, time.Now().Unix()).Find(s).Error; err != nil {

		fmt.Println(err)
		return nil, false
	}
	return s, true
}

func (d *dbStore) Save(s StoreData) error {

	return d.DBGetter().Save(&s).Error

}

func (d *dbStore) Del(t string) error {
	return d.DBGetter().Where("token = ?", t).Delete(&StoreData{}).Error

}

func (d *dbStore) Each(func(StoreData)) {}

func (d *dbStore) CleanUp() {
	d.DBGetter().Where("expire_at < ?", time.Now().Unix()).Delete(&StoreData{})
}

func (d *dbStore) Users(userID int) ([]StoreData, error) {
	r := []StoreData{}
	if err := d.DBGetter().Where("user_id = ?", userID).Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil

}

func (d *dbStore) BatchUpdateByUser(userID int, not string, data string) error {
	return d.DBGetter().Model(&StoreData{}).Where("user_id = ? AND token != ?", userID, not).Update("data", data).Error
}

type StoreData struct {
	Token      string `gorm:"column:token;type:varchar(36);primary_key;not null;default:''"`
	Data       string `gorm:"column:data;type:text;not null"`
	UserID     int    `gorm:"column:user_id;type:int(10);unsigned;not null;default:0"`
	CreateTime int64  `gorm:"column:create_time;type:int(10);unsigned;not null;default:0"`
	UpdateTime int64  `gorm:"column:update_time;type:int(10);unsigned;not null;default:0"`
	ExpireAt   int64  `gorm:"column:expire_at;type:int(10);unsigned;not null;default:0"`
}

func (StoreData) TableName() string {

	return sessionTableName
}

type options struct {
	Path     string
	Domain   string
	MaxAge   int64
	Secure   bool
	HttpOnly bool
}

type Option struct {
	o func(*options)
}

func HttpOnlyOption(httpOnly bool) Option {
	return Option{func(o *options) {
		o.HttpOnly = httpOnly
	}}
}

func SecureOption(secure bool) Option {
	return Option{func(o *options) {
		o.Secure = secure
	}}
}

func PathOption(path string) Option {
	return Option{func(o *options) {
		o.Path = path
	}}
}

func DomainOption(domain string) Option {
	return Option{func(o *options) {
		o.Domain = domain
	}}
}

func MaxAgeOption(maxAge int64) Option {
	return Option{func(o *options) {
		o.MaxAge = maxAge
	}}
}

type session struct {
	token    string
	values   map[string]interface{}
	options  *options
	modified bool
	created  int64
	//optionModified bool
	isNew    bool
	knockout bool
}

func (s *session) setOptions(os ...Option) {
	for _, o := range os {
		o.o(s.options)
	}
	//s.optionModified = true
}

func (s *session) Token() string {
	return s.token
}

func (s *session) Knockout() {
	s.modified = true
	s.Clean()
}

func (s *session) Get(k string) (interface{}, bool) {
	v, ok := s.values[k]
	return v, ok
}

func (s *session) Set(k string, v interface{}) {
	s.values[k] = v
	s.modified = true
}

func (s *session) Del(k string) {
	delete(s.values, k)
	s.modified = true
}

func (s *session) Clean() {
	s.values = map[string]interface{}{}
	s.modified = true
}

func (s *session) GetString(k string) (r string) {
	v, _ := s.values[k]
	r, _ = v.(string)
	return

}

func (s *session) GetInt(k string) (r int) {
	v, _ := s.values[k]
	r, _ = v.(int)
	return
}

func (s *session) GetBool(k string) (r bool) {
	v, _ := s.values[k]

	r, _ = v.(bool)
	return
}

func (s *session) GetStringSlice(k string) (r []string) {
	v, _ := s.values[k]

	r, _ = v.([]string)
	return
}

func (s *session) HttpCookie() *http.Cookie {
	return &http.Cookie{
		Name:   cookieName,
		Value:  s.token,
		Path:   s.options.Path,
		Domain: s.options.Domain,
		MaxAge: int(s.options.MaxAge),
	}
}

type SessionComponent struct {
}

func (SessionComponent) Init() error {
	startCleanUp()
	return nil
}

func SessionFromGin(context *gin.Context) (*session, bool) {
	t, ok := context.Get(contextSessionKey)
	if !ok {
		return nil, ok
	}
	s, ok := t.(*session)
	return s, ok
}

func startCleanUp() {

	go func() {
		for {
			time.Sleep(time.Second * 10)
			if store == nil {
				continue
			}
			store.CleanUp()

		}

	}()
}

func SetCodec(c Codec) {
	codec = c
}

func SetStore(s Store) {
	store = s
}

func Middleware(context *gin.Context) {

	if err := requireChecker(); err != nil {
		fmt.Println(err)
		context.AbortWithStatus(500)
		return
	}

	s := &session{
		token:    newToken(),
		values:   map[string]interface{}{},
		options:  &options{MaxAge: int64(defaultMaxAge)},
		modified: false,
		created:  time.Now().Unix(),
		//optionModified: false,
		isNew: true,
	}

	cookie, err := context.Request.Cookie(cookieName)
	if err == nil {
		s.token = cookie.Value
		d, ok := store.Get(s.token)

		if ok {
			s.created = d.CreateTime
			s.setOptions(MaxAgeOption(d.ExpireAt - d.CreateTime))
			s.isNew = false
			if err := codec.Decode([]byte(d.Data), &s.values); err != nil {
				fmt.Println(err)
				context.AbortWithStatus(500)
				return
			}
		}
	}

	if s.isNew {
		context.SetCookie(cookieName, s.token, int(s.options.MaxAge), s.options.Path, s.options.Domain, s.options.Secure, s.options.HttpOnly)
	}

	context.Set(contextSessionKey, s)

	context.Next()

	if !s.modified {
		return
	}

	if s.knockout {
		store.Del(s.token)
		return
	}

	if err := sessionSave(s); err != nil {
		context.AbortWithStatus(500)
	}
}

func sessionSave(s *session) error {
	d, err := codec.Encode(s.values)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sd := StoreData{
		Token:      s.token,
		Data:       string(d),
		UpdateTime: time.Now().Unix(),
		CreateTime: s.created,
		ExpireAt:   s.created + s.options.MaxAge,
	}

	if userID, ok := s.values["userID"]; ok {
		sd.UserID = userID.(int)
	}

	if err := store.Save(sd); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Update(s *session) error {
	if err := requireChecker(); err != nil {
		fmt.Println(err)
		return err
	}
	return sessionSave(s)

}

func BatchUpdateByUser(userID int, not string, val map[string]interface{}) error {

	if err := requireChecker(); err != nil {
		fmt.Println(err)
		return err
	}

	data, err := codec.Encode(val)
	if err != nil {
		return err
	}
	return store.BatchUpdateByUser(userID, not, string(data))
}

func UserSessions(userID int) ([]*session, error) {
	if err := requireChecker(); err != nil {
		return nil, err
	}

	sds, err := store.Users(userID)
	if err != nil {
		return nil, err
	}

	r := []*session{}
	for _, sd := range sds {

		t := &session{
			token:  sd.Token,
			values: map[string]interface{}{},
			options: &options{
				MaxAge: sd.ExpireAt - sd.CreateTime,
			},
			modified: false,
			created:  sd.CreateTime,
			isNew:    false,
			knockout: false,
		}

		if err := codec.Decode([]byte(sd.Data), &t.values); err != nil {
			return nil, err
		}

		r = append(r, t)

	}

	return r, nil
}

func Del(token string) error {
	if err := requireChecker(); err != nil {
		return err
	}
	return store.Del(token)
}
