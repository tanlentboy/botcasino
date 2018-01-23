package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zhangpanyi/basebot/logger"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// Service 服务配置
type Service struct {
	Address string `yaml:"address"` // 地址
	Port    int    `yaml:"port"`    // 端口号
}

// MySQLCfg MySQL配置
type MySQLCfg struct {
	User     string            `yaml:"user"`     // 访问用户
	Password string            `yaml:"password"` // 访问密码
	Database string            `yaml:"db"`       // 数据库名
	Host     string            `yaml:"host"`     // 服务地址
	Conns    int               `yaml:"conns"`    // 连接数量
	Options  map[string]string `yaml:"options"`  // 附加选项
}

// Serve 服务配置
type Serve struct {
	BindAddress       string   `yaml:"bind-address"`        // 绑定地址
	Port              uint16   `yaml:"port"`                // 端口号
	Domain            string   `yaml:"domain"`              // 服务域名
	APIWebsite        string   `yaml:"api_website"`         // API服务站点
	Token             string   `yaml:"token"`               // 机器人token
	BucketNum         uint32   `yaml:"bucket_num"`          // 记录桶数量
	Account           string   `yaml:"account"`             // 账户名称
	MySQL             MySQLCfg `yaml:"mysql"`               // 数据库配置
	WalletService     Service  `yaml:"wallet_service"`      // 钱包服务配置
	Dynamic           string   `yaml:"dynamic"`             // 动态文件配置
	BolTDBPath        string   `yaml:"boltdb_path"`         // BoltDB路径
	Languages         string   `yaml:"languages"`           // 语言配置路径
	RedEnvelopesCover []string `yaml:"red_envelopes_cover"` // 红包封面图
}

// parser 配置解析器
type parser interface {
	parse([]byte) error
}

// Manager 配置管理器
type Manager struct {
	serve      *Serve
	languges   *Languges
	dynamic    *Dynamic
	watcher    *fsnotify.Watcher
	fileparser map[string]parser
}

// GetServe 获取服务配置
func GetServe() Serve {
	return *globalManager.serve
}

// GetLanguge 获取语言配置
func GetLanguge() *Languges {
	return globalManager.languges
}

// GetDynamic 获取动态配置
func GetDynamic() DynamicCfg {
	return globalManager.dynamic.DynamicCfg
}

// LoadConfig 加载配置文件
func LoadConfig(path string) {
	once.Do(func() {
		// 创建观察器
		fileparser := make(map[string]parser)
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}

		// 加载主配置
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		serve := Serve{}
		err = yaml.Unmarshal(data, &serve)
		if err != nil {
			panic(err)
		}

		// 加载动态配置
		dynamic := readDynamic(serve.Dynamic)
		watcher.Add(serve.Dynamic)
		fileparser[serve.Dynamic] = dynamic

		// 加载语言包配置
		languages, files := readLanguages(serve.Languages)
		for _, filename := range files {
			watcher.Add(filename)
			fileparser[filename] = languages
		}

		// 初始化全局配置
		globalManager = &Manager{
			serve:      &serve,
			languges:   languages,
			dynamic:    dynamic,
			fileparser: fileparser,
			watcher:    watcher,
		}
		go globalManager.watch()
	})
}

// 全局配置管理器
var once sync.Once
var globalManager *Manager

// 观察文件变更
func (m *Manager) watch() {
	for {
		select {
		case evt := <-m.watcher.Events:
			if evt.Op&fsnotify.Write == fsnotify.Write {
				handler, ok := m.fileparser[evt.Name]
				if ok {
					data, err := ioutil.ReadFile(evt.Name)
					if err == nil {
						err = handler.parse(data)
						if err != nil {
							logger.Warnf("File notify: parse failed, %v, %v", evt.Name, err)
						} else {
							logger.Infof("File notify: realod file finished, %v", evt.Name)
						}
					} else {
						logger.Warnf("File notify: handle failed, %v, %v", evt.Name, err)
					}
				} else {
					logger.Warnf("File notify: ignore write event, %v", evt.Name)
				}
			}
		case err := <-m.watcher.Errors:
			logger.Warnf("File notify: recv error event, %v", err)
		}
	}
}

// 读取动态配置
func readDynamic(filename string) *Dynamic {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	dynamic := Dynamic{}
	err = dynamic.parse(data)
	if err != nil {
		panic(err)
	}
	return &dynamic
}

// 读取语言包配置
func readLanguages(dir string) (*Languges, []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	languges := NewLanguges()
	paths := make([]string, 0)
	for _, v := range files {
		if !v.IsDir() {
			ext := strings.ToLower(filepath.Ext(v.Name()))
			if ext == ".lang" {
				fullname := dir + string(filepath.Separator) + v.Name()
				data, err := ioutil.ReadFile(fullname)
				if err != nil {
					panic(err)
				}
				err = languges.parse(data)
				if err != nil {
					panic(err)
				}
				paths = append(paths, fullname)
			}
		}
	}
	return languges, paths
}
