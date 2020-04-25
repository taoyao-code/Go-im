package config

type QiNiuYun struct {
	QINIU_DOMAIN      string
	QINIU_ACCESS_KEY  string
	QINIU_SECRET_KEY  string
	QINIU_TEST_BUCKET string
}

type MySQL struct {
	Address  string
	Port     string
	Username string
	Password string
	Database string
}
type Redis struct {
	Address  string
	Port     string
	Password string
}

type Configuration struct {
	QNY   QiNiuYun
	MySQL MySQL
	Redis Redis
}
