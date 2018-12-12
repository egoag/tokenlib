## Tokenlib

Golang port of [tokenlib](https://github.com/mozilla-services/tokenlib).

### Examples

```go
import (
	"fmt"
	"time"

	"github.com/youyaochi/tokenlib"
)

func main() {
	secret := "I_LIKE_UNICORNS"
	data := map[string]interface{}{
		"Id":   1,
		"Name": "One",
	}
	fmt.Println("Raw data:", data)

	token, err := tokenlib.MakeToken(data, secret, 1)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println("Token:", token)

	newData, err := tokenlib.ParseToken(token, secret)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println("Parsed Data:", newData)

	time.Sleep(2 * time.Second)
	expiredData, err := tokenlib.ParseToken(token, secret)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println("Expired data:", expiredData)
}
```
