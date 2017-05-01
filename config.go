package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    ini "gopkg.in/ini.v1"
)

type HardwareConfig struct {
    Pin     string
    Reverse bool
    Pull    bool
}

type ActionConfig struct {
    Command string
}

type Config struct {
    Name string

    HardwareConfig `ini:"GPIO"`
    ActionConfig   `ini:"Action"`
}

func LoadConfig(path string) ([]Config, error) {
    log.Printf("Loading configs from: %s", path)

    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, fmt.Errorf("Failed to open directory: %s: %v", path, err)
    }

    var result []Config = make([]Config, 0)
    for _, f := range files {
        if f.IsDir() {
            continue
        }

        log.Printf("Loading config file: %s", f.Name())

        cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, f.Name())
        if err != nil {
            return nil, fmt.Errorf("Failed to parse config: %s: %v", f.Name(), err)
        }

        // Enable expansion of env vars
        cfg.ValueMapper = os.ExpandEnv

        // Map to configPath structure
        var config Config
        if err := cfg.MapTo(&config); err != nil {
            return nil, fmt.Errorf("Failed to load config: %s: %v", f.Name(), err)
        }

        log.Printf("Loaded config: %v", config)

        result = append(result, config)
    }

    return result, nil
}
