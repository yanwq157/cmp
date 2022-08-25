package conf

type Mysql struct {
	Path     string `mapstructure:"path" json:"path" yaml:"path"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname   string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
type Zap struct {
	Level string `mapstructure:"level" json:"level" yaml:"level"`
}
type Server struct {
	Mysql Mysql `mapstructure:"mysql"  json:"mysql" yaml:"mysql"`
	Zap   Zap   `mapstructure:"zap"  json:"zap" yaml:"zap"`
}
