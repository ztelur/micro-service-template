package bootstrap

var (
	HttpConfig  HttpConf
	MySQlConfig MySQLConf
)

//Http配置
type HttpConf struct {
	Host string
	Port string
}

type MySQLConf struct {
	DriverName string
	UrlAddress string
}

func NewMySQLConf() *MySQLConf {
	return &MySQlConfig
}
