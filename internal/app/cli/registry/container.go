package registry

import (
    "github.com/drybin/minter-pools-ar/internal/app/cli/config"
    "github.com/drybin/minter-pools-ar/internal/app/cli/usecase"
    "github.com/drybin/minter-pools-ar/pkg/logger"
)

type Container struct {
    Logger   logger.ILogger
    Usecases *Usecases
    Clean    func()
}

type Usecases struct {
    HelloWorld *usecase.HelloWorld
}

func NewContainer(
    config *config.Config,
) (*Container, error) {
    log := logger.NewLogger()
    
    //httpClient := resty.New()
    
    container := Container{
        Logger: log,
        Usecases: &Usecases{
            HelloWorld: usecase.NewHelloWorldUsecase(),
        },
        Clean: func() {
        },
    }
    
    return &container, nil
}
