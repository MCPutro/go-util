package snowflake

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestGenerateSnowflake(t *testing.T) {
	snow, err := NewGenerator(123)
	if err != nil {
		t.Error(err)
		return
	}

	group := sync.WaitGroup{}
	a := time.Now()
	for i := 0; i < 10000; i++ {
		go func() {
			group.Add(1)
			id, err := snow.NextID()
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("%d Generated ID: %d", i, id)
			group.Done()
		}()
	}
	group.Wait()
	b := time.Now()
	fmt.Println("All IDs generated successfully")
	sub := b.Sub(a)
	log.Println("start time: ", a)
	log.Println("end   time: ", b)
	log.Println("selisih : ", sub)
	log.Println("selisih dalam menit: ", sub.Minutes())
	log.Println("selisih dalam Seconds: ", sub.Seconds())
	log.Println("selisih dalam Milliseconds: ", sub.Milliseconds())

}
