package agents

import "time"

type Agent struct {
    Name string
    Run  func()
}

func (a *Agent) Start() {
    go func() {
        for {
            a.Run()
            time.Sleep(10 * time.Millisecond)
        }
    }()
}
