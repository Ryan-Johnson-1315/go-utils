# GoLang Misc Utilities
Test each one of the packages

```
go run . [fps]
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
