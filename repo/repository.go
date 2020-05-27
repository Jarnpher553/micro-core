package repo

import (
	"fmt"
	"github.com/Jarnpher553/micro-core/log"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)
import _ "github.com/jinzhu/gorm/dialects/mysql"

// Repository 仓储类
type Repository struct {
	*gorm.DB
	*log.LogrusEntry
	userName string
	password string
	addr     string
	host     string
	port     string
	dbName   string
	logMode  bool
}

// FieldName 字段类
type FieldName struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

// Option 配置函数
type Option func(*Repository)

// UserName 用户名配置
func UserName(userName string) Option {
	return func(repository *Repository) {
		repository.userName = userName
	}
}

func LogMode(mode bool) Option {
	return func(repository *Repository) {
		repository.logMode = mode
	}
}

// Pwd 密码配置
func Pwd(password string) Option {
	return func(repository *Repository) {
		repository.password = password
	}
}

func Addr(addr string) Option {
	return func(repository *Repository) {
		repository.addr = addr
	}
}

// Host 服务器配置
func Host(host string) Option {
	return func(repository *Repository) {
		repository.host = host
	}
}

// Port 端口配置
func Port(port string) Option {
	return func(repository *Repository) {
		repository.port = port
	}
}

// DbName 数据库名配置
func DbName(dbName string) Option {
	return func(repository *Repository) {
		repository.dbName = dbName
	}
}

var entry = log.Logger.Mark("Repo")

// New 构造函数
func New(options ...Option) *Repository {
	if len(options) != 4 {
		return nil
	}

	repo := &Repository{
		LogrusEntry: entry,
	}

	for i := range options {
		options[i](repo)
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", repo.userName, repo.password, repo.addr, repo.dbName))

	if err != nil {
		entry.Fatalln(err)
	}

	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)
	db.SetLogger(repo)
	db.LogMode(repo.logMode)
	repo.DB = db

	return repo
}

// Deprecated: NewFromConfigFile 通过配置文件实例化repo
/*func NewFromConfigFile(file *config.Config, fn *FieldName) *Repository {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", file.GetString(fn.Username), file.GetString(fn.Password), file.GetString(fn.Host), file.GetString(fn.Port), file.GetString(fn.DbName)))

	if err != nil {
		entry.Fatalln(err)
	}

	return &Repository{
		DB: db,
	}
}*/

// ReadAll 查询单条
func (repo *Repository) ReadAll(out interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Find(out, where...).Error
}

// Read 查询单条记录
func (repo *Repository) Read(out interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.First(out, where...).Error
}

func (repo *Repository) Exist(out interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	sql := repo.DB.First(out, where...)
	not := sql.RecordNotFound()
	if not || sql.Error != nil {
		return sql.Error
	} else {
		return nil
	}
}

// Remove 删除
func (repo *Repository) Remove(val interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	return repo.DB.Delete(val, where...).Error
}

// Insert 新增
func (repo *Repository) Insert(val interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Create(val).Error
}

func (repo *Repository) SoftRemove(value interface{}, where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	columns := make(map[string]interface{})
	columns["is_active"] = false

	if where != nil {
		return repo.DB.Model(value).Where(where[0], where[1:]...).UpdateColumns(columns).Error
	}
	return repo.DB.Model(value).UpdateColumns(columns).Error
}

// Modify 修改
func (repo *Repository) Modify(val interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	return repo.DB.Save(val).Error
}

// 更改单个字段
func (repo *Repository) ModifyColumn(val interface{}, attr string, upValue interface{}, where ...interface{}) (affects int64, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	kind := reflect.TypeOf(val).Kind()

	if where != nil {
		if kind == reflect.String {
			db := repo.DB.Table(val.(string)).Where(where[0], where[1:]...).Update(attr, upValue)
			return db.RowsAffected, db.Error
		} else {
			db := repo.DB.Model(val).Where(where[0], where[1:]...).Update(attr, upValue)
			return db.RowsAffected, db.Error
		}
	}
	if kind == reflect.String {
		db := repo.DB.Table(val.(string)).Update(attr, upValue)
		return db.RowsAffected, db.Error
	} else {
		db := repo.DB.Model(val).Update(attr, upValue)
		return db.RowsAffected, db.Error
	}
}

// 更改多个字段
func (repo *Repository) ModifyColumns(val interface{}, columns interface{}, where ...interface{}) (affects int64, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()
	kind := reflect.TypeOf(val).Kind()

	if where != nil {
		if kind == reflect.String {
			db := repo.DB.Table(val.(string)).Where(where[0], where[1:]...).Updates(columns)
			return db.RowsAffected, db.Error
		} else {
			db := repo.DB.Model(val).Where(where[0], where[1:]...).Updates(columns)
			return db.RowsAffected, db.Error
		}
	}
	if kind == reflect.String {
		db := repo.DB.Table(val.(string)).Updates(columns)
		return db.RowsAffected, db.Error
	} else {
		db := repo.DB.Model(val).Updates(columns)
		return db.RowsAffected, db.Error
	}
}

// ModifyFunc 使用函数更新
func (repo *Repository) ModifyFunc(val interface{}, modifier func(interface{}), where ...interface{}) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)

	valNew := reflect.New(t.Elem())
	valNew.Elem().Set(v.Elem())

	i := valNew.Interface()

	modifier(i)

	if where != nil {
		return repo.DB.Model(val).Where(where[0], where[1:]...).Updates(i).Error
	}
	return repo.DB.Model(val).Updates(i).Error
}

// Transaction 执行包装事务
func (repo *Repository) Transaction(f func(*Repository) error) (e error) {
	repoTx := repo.begin()

	defer func() {
		if r := recover(); r != nil {
			repoTx.Rollback()
			e = fmt.Errorf("%v", r)
		}
	}()

	if err := repoTx.Error; err != nil {
		return err
	}

	if err := f(repoTx); err != nil {
		repoTx.Rollback()
		return err
	}

	if err := repoTx.Commit().Error; err != nil {
		repoTx.Rollback()
		return err
	}
	return nil
}

// begin 开始一个事务
func (repo *Repository) begin() *Repository {
	//开始一个事务
	tx := repo.DB.Begin()

	//返回一个包含事务的repo
	return &Repository{
		userName:    repo.userName,
		password:    repo.password,
		host:        repo.host,
		addr:        repo.addr,
		port:        repo.port,
		dbName:      repo.dbName,
		DB:          tx,
		LogrusEntry: repo.LogrusEntry,
	}
}

type Expression func(db *gorm.DB) *gorm.DB

func Page(pageNum int, perCount int) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((pageNum - 1) * perCount).Limit(perCount)
	}
}

func Model(value interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(value)
	}
}

func Table(name string) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(name)
	}
}

func Select(query interface{}, args ...interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(query, args...)
	}
}

func Order(value interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}

func Join(query string, args ...interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

func Where(query interface{}, args ...interface{}) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func Group(query string) Expression {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(query)
	}
}

// Query 查询列表（可分页）
func (repo *Repository) Query(out interface{}, count bool, exps ...Expression) (c int, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	db := repo.DB

	for _, exp := range exps {
		db = exp(db)
	}

	db = db.Scan(out)
	if err := db.Error; err != nil {
		e = err
		return
	}

	if count {
		if err := db.
			Offset(-1).
			Limit(-1).
			Count(&c).Error; err != nil {
			e = err
			return
		}
	}

	return
}

func (repo *Repository) Print(args ...interface{}) {
	formatter := gorm.LogFormatter(args...)
	repo.LogrusEntry.Info(formatter[2], formatter[3], strings.Replace(formatter[4].(string), "\n", "", -1))
}

func (repo *Repository) Migrate(initial func(*Repository), values ...interface{}) {
	repo.DB.AutoMigrate(values...)

	if initial != nil {
		initial(repo)
	}
}
