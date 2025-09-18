### Example
```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/Yapcheekian/worker"
)

func main() {
    pool := worker.New(context.Background(), 5)

    for i := 0; i < 10; i++ {
        idx := i
        if err := pool.Add(func() {
            fmt.Printf("Processing task %d\n", idx)
            time.Sleep(500 * time.Millisecond)
        }); err != nil {
            fmt.Println("Failed to add task:", err)
        }
    }

    pool.Wait()

    fmt.Println("All tasks completed")
}
```
