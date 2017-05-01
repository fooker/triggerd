package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    ini "gopkg.in/ini.v1"
    "path"
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

func LoadConfig(dir string) ([]Config, error) {
    log.Printf("Loading configs from: %s", dir)

    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, fmt.Errorf("Failed to open directory: %s: %v", dir, err)
    }

    var result []Config = make([]Config, 0)
    for _, f := range files {
        name := path.Join(dir, f.Name())

        if f.IsDir() {
            continue
        }

        log.Printf("Loading config file: %s", name)

        cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, name)
        if err != nil {
            return nil, fmt.Errorf("Failed to parse config: %s: %v", name, err)
        }

        // Enable expansion of env vars
        cfg.ValueMapper = os.ExpandEnv

        // Map to configPath structure
        var config Config
        if err := cfg.MapTo(&config); err != nil {
            return nil, fmt.Errorf("Failed to load config: %s: %v", name, err)
        }

        log.Printf("Loaded config: %v", config)

        result = append(result, config)
    }

    return result, nil
}
