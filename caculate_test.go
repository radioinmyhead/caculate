package caculate

import (
	"testing"
)

func TestCal(t *testing.T) {
	dic := map[string]string{"a": "10", "b": "100", "c": "ok"}
	cases := map[string]bool{
		`${a}>1&& ${a}>50||${b}>10&&"${c}"!="ok"`: false,
		`${a}>1&&( ${a}>50||${b}>10)`:             true,
	}
	for k, v := range cases {
		testv, err := Caculate(dic, k)
		if err != nil {
			t.Errorf("caculate error: %v", err)
			continue
		}
		if testv != v {
			t.Errorf("dic=%v,caculate `%v`: want: %v get: %v", dic, k, v, testv)
		}
	}
}
