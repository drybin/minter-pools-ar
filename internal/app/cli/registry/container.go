package registry

import (
    "github.com/MinterTeam/minter-go-sdk/v2/api/http_client"
    "github.com/drybin/minter-pools-ar/internal/adapter/webapi"
    "github.com/drybin/minter-pools-ar/internal/app/cli/config"
    "github.com/drybin/minter-pools-ar/internal/app/cli/usecase"
    "github.com/drybin/minter-pools-ar/pkg/logger"
    "github.com/drybin/minter-pools-ar/pkg/telegram"
    "github.com/drybin/minter-pools-ar/pkg/wrap"
    "github.com/go-resty/resty/v2"
)

type Container struct {
    Logger   logger.ILogger
    Usecases *Usecases
    Clean    func()
}

type Usecases struct {
    HelloWorld *usecase.HelloWorld
    Search     *usecase.Search
    SearchWeb  *usecase.SearchWeb
}

func NewContainer(
    config *config.Config,
) (*Container, error) {
    log := logger.NewLogger()
    
    httpClient := resty.New()
    
    minterClient, err := http_client.New(config.MinterApiUrl)
    if err != nil {
        return nil, wrap.Errorf("failed to create Minter client: %w", err)
    }
    
    //Minter console from other coin(not BIP) use this endpoints
    minterClientGate, err := http_client.New(config.MinterApiGateUrl)
    if err != nil {
        return nil, wrap.Errorf("failed to create Minter gate client: %w", err)
    }
    
    chainikApi := webapi.NewChainikWebapi(httpClient)
    
    container := Container{
        Logger: log,
        Usecases: &Usecases{
            HelloWorld: usecase.NewHelloWorldUsecase(),
            Search: usecase.NewSearchUsecase(
                chainikApi,
                webapi.NewMinterWebapi(minterClient, minterClientGate, config.PassPhrase),
            ),
            SearchWeb: usecase.NewSearchWebUsecase(
                webapi.NewMinterWeb(httpClient),
                webapi.NewMinterWebapi(minterClient, minterClientGate, config.PassPhrase),
                telegram.NewTelegramWebapi(httpClient, config.TgConfig.BotToken, config.TgConfig.ChatId),
            ),
        },
        Clean: func() {
        },
    }
    
    return &container, nil
}
