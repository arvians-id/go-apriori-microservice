package config

import "github.com/spf13/viper"

type Config struct {
	AppName  string `mapstructure:"APP_NAME"`
	AppUrl   string `mapstructure:"APP_URL"`
	AppUrlFE string `mapstructure:"APP_URL_FE"`
	Port     string `mapstructure:"PORT"`
	XApiKey  string `mapstructure:"X_API_KEY"`

	AprioriSvcUrl       string `mapstructure:"APRIORI_SERVICE_URL"`
	CategorySvcUrl      string `mapstructure:"CATEGORY_SERVICE_URL"`
	CommentSvcUrl       string `mapstructure:"COMMENT_SERVICE_URL"`
	NotificationSvcUrl  string `mapstructure:"NOTIFICATION_SERVICE_URL"`
	PasswordResetSvcUrl string `mapstructure:"PASSWORD_RESET_SERVICE_URL"`
	PaymentSvcUrl       string `mapstructure:"PAYMENT_SERVICE_URL"`
	ProductSvcUrl       string `mapstructure:"PRODUCT_SERVICE_URL"`
	TransactionSvcUrl   string `mapstructure:"TRANSACTION_SERVICE_URL"`
	UserOrderSvcUrl     string `mapstructure:"USER_ORDER_SERVICE_URL"`
	UserSvcUrl          string `mapstructure:"USER_SERVICE_URL"`
	MessageBrokerUrl    string `mapstructure:"MESSAGE_BROKER_URL"`

	DBConnection string `mapstructure:"DB_CONNECTION"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUsername   string `mapstructure:"DB_USERNAME"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBDatabase   string `mapstructure:"DB_DATABASE"`
	DBSSLMode    string `mapstructure:"DB_SSL_MODE"`

	DBPoolMin           int `mapstructure:"DATABASE_POOL_MIN"`
	DBPoolMax           int `mapstructure:"DATABASE_POOL_MAX"`
	DBMaxIdleTimeSecond int `mapstructure:"DATABASE_MAX_IDLE_TIME_SECOND"`
	DBMaxLifeTimeSecond int `mapstructure:"DATABASE_MAX_LIFE_TIME_SECOND"`

	MailMailer      string `mapstructure:"MAIL_MAILER"`
	MailHost        string `mapstructure:"MAIL_HOST"`
	MailPort        string `mapstructure:"MAIL_PORT"`
	MailUsername    string `mapstructure:"MAIL_USERNAME"`
	MailPassword    string `mapstructure:"MAIL_PASSWORD"`
	MailEncryption  string `mapstructure:"MAIL_ENCRYPTION"`
	MailFromAddress string `mapstructure:"MAIL_FROM_ADDRESS"`

	JwtSecretAccessKey    string `mapstructure:"JWT_SECRET_ACCESS_KEY"`
	JwtSecretRefreshKey   string `mapstructure:"JWT_SECRET_REFRESH_KEY"`
	JwtAccessExpiredTime  int    `mapstructure:"JWT_ACCESS_EXPIRED_TIME"`
	JwtRefreshExpiredTime int    `mapstructure:"JWT_REFRESH_EXPIRED_TIME"`

	AwsAccessKeyId string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretKey   string `mapstructure:"AWS_SECRET_KEY"`
	AwsRegion      string `mapstructure:"AWS_REGION"`
	AwsBucket      string `mapstructure:"AWS_BUCKET"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	RedisExpired  int    `mapstructure:"REDIS_EXPIRED"`

	MidtransServerKey string `mapstructure:"MIDTRANS_SERVER_KEY"`
	MidtransClientKey string `mapstructure:"MIDTRANS_CLIENT_KEY"`

	RajaOngkirSecretKey string `mapstructure:"RAJA_ONGKIR_SECRET_KEY"`
}

func LoadConfig(filenames ...string) (c *Config, err error) {
	if filenames != nil {
		viper.AddConfigPath(filenames[0])
	} else {
		viper.AddConfigPath("./config/envs")
	}

	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
