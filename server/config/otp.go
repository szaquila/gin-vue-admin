package config

type Otp struct {
	Enable bool   `mapstructure:"enable" json:"enable" yaml:"enable"` // 防爆破验证码开启此数，0代表每次登录都需要验证码，其他数字代表错误密码此数，如3代表错误三次后出现验证码
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Length int    `mapstructure:"open-captcha-timeout" json:"open-captcha-timeout" yaml:"open-captcha-timeout"` // 防爆破验证码超时时间，单位：s(秒)
	Ssh    string `mapstructure:"ssh" json:"ssh" yaml:"ssh"`
}
