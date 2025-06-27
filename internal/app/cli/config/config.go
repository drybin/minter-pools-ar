package config

import (
    "errors"
    "time"
    
    "github.com/drybin/minter-pools-ar/pkg/env"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
    validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
    ServiceName          string
    MinterApiUrl         string
    MinterApiGateUrl     string
    MinterApiExplorerUrl string
    PassPhrase           string
    TgConfig             TgConfig
}

type TgConfig struct {
    BotToken string
    ChatId   string
    Timeout  time.Duration
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
        ServiceName:          env.GetString("APP_NAME", "minter-pools-ar"),
        MinterApiUrl:         env.GetString("MINTER_API_URL", "https://api.minter.one/v2/"),
        MinterApiGateUrl:     env.GetString("MINTER_API_GATE_URL", "https://gate-api.minter.network/api/v2/"),
        MinterApiExplorerUrl: env.GetString("MINTER_API_GATE_URL", "https://explorer-api.minter.network/api/v2/"),
        PassPhrase:           env.GetString("PASS_PHRASE", ""),
        TgConfig:             initTgConfig(),
    }
    
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func initTgConfig() TgConfig {
    return TgConfig{
        BotToken: env.GetString("TG_BOT_TOKEN", ""),
        ChatId:   env.GetString("TG_CHAT_ID", ""),
        Timeout:  10 * time.Second,
    }
}
