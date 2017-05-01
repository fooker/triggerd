package main

import (
    "flag"
    "log"
    "os"
    "os/exec"
    "os/signal"
    "strings"
    "time"
    "github.com/google/shlex"
    "periph.io/x/periph/host"
    "periph.io/x/periph/conn/gpio"
    "periph.io/x/periph/conn/gpio/gpioreg"
)

var configPath = flag.String("config", "/etc/triggerd", "The configPath directory")

type Trigger struct {
    name string

    pin     gpio.PinIO
    reverse bool

    command []string
}

func (trigger *Trigger) Handle() {
    counter := 0

    for {
        if (trigger.pin.Read() == gpio.High) != trigger.reverse {
            counter += 1
        } else {
            counter = 0
        }

        if counter == 3 {
            log.Printf("Triggered: %s", trigger.name)

            cmd := exec.Cmd{
                Path: trigger.command[0],
                Args: trigger.command,
            }
            if err := cmd.Start(); err != nil {
                log.Printf("Failed to execute command: %s: %v", strings.Join(trigger.command, " "), err)
            }

            time.Sleep(250 * time.Millisecond)
        }

        time.Sleep(10 * time.Millisecond)
    }
}

func main() {
    flag.Parse()

    configs, err := LoadConfig(*configPath)
    if err != nil {
        panic(err)
    }

    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

    for _, config := range configs {
        trigger := &Trigger{}

        trigger.name = config.Name

        // Parse command to execute an trigger
        trigger.command, err = shlex.Split(config.Command)
        if err != nil {
            panic(err)
        }

        // Determine pull up/down resistors
        var pull gpio.Pull = gpio.Float
        if config.Pull {
            if !config.Reverse {
                pull = gpio.PullDown
            } else {
                pull = gpio.PullUp
            }
        }

        // Configure pin
        trigger.pin = gpioreg.ByName(config.Pin)
        if err := trigger.pin.In(pull, gpio.NoEdge); err != nil {
            panic(err)
        }
        trigger.reverse = config.Reverse

        go trigger.Handle()
    }

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    for range signalChan {
        break
    }
}
