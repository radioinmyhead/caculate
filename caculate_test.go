package caculate

import (
	"testing"
)

func TestCal(t *testing.T) {
	dic := map[string]string{"a5": "10", "b": "100", "c": "ok"}
	cases := map[string]bool{
		`${a5}>1&& ${a5}>50||${b}>10&&"${c}"!="ok"`: false,
		`${a5}>1&&( ${a5}>50||${b}>10)`:             true,
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
