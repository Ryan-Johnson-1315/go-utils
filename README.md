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
Inter-process Communication logger. Allows several processes/languages to log to the same place. Accepts JSON and prints to console

## GoLang
```
func main() {
    ip, port := "127.0.0.1", 50000
    remoteAddr := &net.UDPAddr{
        Port: port,
        IP:   net.ParseIP(ip),
    }
    if err := sockLog.StartLogger(ip, port); err == nil {
        log.Println("Unable to start logger!", err)    
    }

    rand.Seed(time.Now().UnixNano())

    for c := 0; c < 1500; c++ {
        sl := rand.Intn(500)

        msg := sockLog.SocketMessage{
            Caller:      "Test Object",
            MessageType: sockLog.MessageLevel(c % 5),
            Message:     "This is a test message!",
            Function:    "main()",
        }
        l.SendSocketMessage(msg)
        time.Sleep(time.Millisecond * time.Duration(sl))
    }

Output:
///////////////////////////////////////// \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  _________              __           __   ____                                      
 /   _____/ ____   ____ |  | __ _____/  |_|    |    ____   ____   ____   ___________ 
 \_____  \ /  _ \_/ ___\|  |/ // __ \   __\    |   /  _ \ / ___\ / ___\_/ __ \_  __ \
 /        (  <_> )  \___|    <\  ___/|  | |    |__(  <_> ) /_/  > /_/  >  ___/|  | \/
/_______  /\____/ \___  >__|_ \\___  >__| |_______ \____/\___  /\___  / \___  >__|   
        \/            \/     \/    \/             \/    /_____//_____/      \/       
\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ ///////////////////////////////////////////

```
```diff
+2021/04/05 10:19:39 Socket UDP Logger is starting! Listening @ 127.0.0.1:50000
-2021/04/05 10:19:39.722860 <UDP> Test Object::main() This is a test message!
+2021/04/05 10:19:39.834560 <UDP> Test Object::main() This is a test message!
+2021/04/05 10:19:40.279351 <UDP> Test Object::main() This is a test message!
-2021/04/05 10:19:40.748368 <UDP> Test Object::main() This is a test message!
-2021/04/05 10:19:40.859076 <UDP> Test Object::main() This is a test message!
+2021/04/05 10:19:41.184797 <UDP> Test Object::main() This is a test message!
-2021/04/05 10:19:41.450491 <UDP> Test Object::main() This is a test message!
+2021/04/05 10:19:41.540022 <UDP> Test Object::main() This is a test message!

```

## Python3.8

```
import socket
import json

if __name__ == "__main__":
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((ip, port))
    for i in range(0, 3):
        msg = {}
        msg["caller"] = "main.py"
        msg["level"] = i
        msg["message"] = "Hello from python"
        msg["function"] = "__main__"

        s.send(bytes(json.dumps(msg), encoding="utf-8"))
        time.sleep(1)
Output:
```
```diff
+2021/04/05 10:26:20 [127.0.0.1:35640] <TCP> main.py :: __main__ Hello from python
+2021/04/05 10:26:20 [127.0.0.1:35638] <TCP> main.py :: __main__ Hello from python
-2021/04/05 10:26:20 [127.0.0.1:35640] <TCP> main.py :: __main__ Hello from python
-2021/04/05 10:26:20 [127.0.0.1:35638] <TCP> main.py :: __main__ Hello from python
+2021/04/05 10:26:20 [127.0.0.1:35640] <TCP> main.py :: __main__ Hello from python

```

