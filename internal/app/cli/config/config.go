package config

import (
    "errors"
    
    "github.com/drybin/minter-pools-ar/pkg/env"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
    validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
    ServiceName  string
    MinterApiUrl string
    PassPhrase   string
}

func (c Config) Validate() error {
    var errs []error
    
    err := validation.ValidateStruct(&c,
        validation.Field(&c.ServiceName, validation.Required),
        validation.Field(&c.PassPhrase, validation.Required),
    )
    if err != nil {
        return wrap.Errorf("failed to validate cli config: %w", err)
    }
    
    return errors.Join(errs...)
}

func InitConfig() (*Config, error) {
    config := Config{
        ServiceName:  env.GetString("APP_NAME", "minter-pools-ar"),
        MinterApiUrl: env.GetString("MINTER_API_URL", "https://api.minter.one/v2/"),
        PassPhrase:   env.GetString("PASS_PHRASE", ""),
    }
    
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return &config, nil
}
