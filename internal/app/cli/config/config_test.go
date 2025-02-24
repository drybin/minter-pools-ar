package config

import (
    "testing"
)

func TestConfig_Validate(t *testing.T) {
    t.Parallel()
    
    type fields struct {
        ServiceName string
        PostgreeDsn string
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        {
            name: "All config is setted, no error",
            fields: fields{
                ServiceName: "name",
            },
            wantErr: false,
        },
        {
            name: "Service name is empty, error",
            fields: fields{
                ServiceName: "",
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        tt := tt
        
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            c := Config{
                ServiceName: tt.fields.ServiceName,
            }
            if err := c.Validate(); (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
