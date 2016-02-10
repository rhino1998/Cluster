package timelist

import (
	"testing"
	"time"
)

func TestTimeList(t *testing.T) {
	temp := New()

	temp.Insert([]byte("hello"), time.Now().Add(-35*time.Second))
	temp.Insert("wjede", time.Now().Add(65*time.Second))
	temp.Insert("wjee", time.Now().Add(5*time.Second))
	temp.Insert("wjee", time.Now().Add(5*time.Second))
	temp.Append(5)
	t.Log(temp.Length())
	for i, item := range temp.After(time.Now().Add(-10 * time.Second)).Items() {
		temp.Insert("wjee", time.Now().Add(time.Duration(i)*time.Second))
		t.Log(item.Value())
		t.Log(item.Time())
		t.Log(item)
	}
	temp, _ = temp.PopAfter(time.Now())
	temp, _ = temp.PopBefore(time.Now().Add(5 * time.Second))

}
