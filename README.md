# GoLang Misc Utilities
Test each one of the packages

```
go run . [fps|logger]
```
# FPS

Clock code speed

```
import (
    "time"
    
    "github.com/Ryan-Johnson-1315/go-utils/fps"
)


func main() {
    fr := fps.NewFrameRaterWithDescription("Hello World")
    for {
        fr.Tick()
        time.Sleep(15 * time.Millisecond)
    }
}

Output:

2021/02/27 22:09:31 Hello World 66 fps
2021/02/27 22:09:32 Hello World 65 fps
2021/02/27 22:09:33 Hello World 65 fps
2021/02/27 22:09:34 Hello World 65 fps

```

# Socket Logger

## GoLang
```
func main() {
    ip, port := "127.0.0.1", 50000
    remoteAddr := &net.UDPAddr{
        Port: port,
        IP:   net.ParseIP(ip),
    }
    if err := sockLog.StartLogger(ip, port); err == nil {
        addr := &net.UDPAddr{
            Port: 50001,
            IP:   net.ParseIP(ip),
        }

        sock, _ := net.ListenUDP("udp", addr)
        count := 0
        for {
            count++
            msg := sockLog.SocketMessage{
                Caller:      "Test Object",
                MessageType: sockLog.MessageLevel(count % 5),
                Message:     "This is a test message!",
                Function:    "main()",
            }
            bts, _ := json.Marshal(msg)
            sock.WriteToUDP(bts, remoteAddr)
            time.Sleep(time.Second * 1)
        }
}
```

## Python3.8

