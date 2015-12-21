package caculate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func get(x string) (ret string) {
	switch {
	case strings.Index(x, "true") != -1:
		return "bool"
	case strings.Index(x, "false") != -1:
		return "bool"
	case strings.Index(x, `.`) != -1:
		return "float"
	case strings.Index(x, `'`) != -1:
		return "string"
	default:
		return "int"
	}
}
func gettype(a, b string) (ret string) {
	ta, tb := get(a), get(b)
	if ta == tb {
		return ta
	}
	return ""
}
func comp(a, b string) bool {
	dic := map[string]int{
		`*`: 4,
		`/`: 4,
		`+`: 3,
		`-`: 3,
		`=`: 2,
		`!`: 2,
		`<`: 2,
		`>`: 2,
		`|`: 1,
		`&`: 1,
	}
	return dic[a] > dic[b]
}

type stack []string

func (s *stack) Put(i string) {
	*s = append(*s, i)
}
func (s *stack) Pop() (i string) {
	if len(*s) == 0 {
		return
	}
	i = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return
}
func (s *stack) Top() (i string) {
	if len(*s) == 0 {
		return
	}
	return (*s)[len(*s)-1]
}

func filter(s string) (n string) {
	s = strings.Replace(s, ` `, ``, -1)
	s = strings.Replace(s, `!=`, `!`, -1)
	s = strings.Replace(s, `==`, `=`, -1)
	s = strings.Replace(s, `&&`, `&`, -1)
	s = strings.Replace(s, `||`, `|`, -1)
	s = strings.Replace(s, `"`, `'`, -1)
	return s
}

func fparse(s string, d map[string]string) (ns string, err error) {
	re := regexp.MustCompile(`\${([a-zA-Z]+)}`)
	matches := re.FindAllStringSubmatchIndex(s, -1)
	if matches == nil {
		return s, nil
	}
	prev := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		if start > prev {
			ns += s[prev:start]
		}
		name := s[m[2]:m[3]]
		if _, ok := d[name]; !ok {
			return "", fmt.Errorf("${%v} not found", name)
		}
		ns += d[name]
		prev = end
	}
	end := s[prev:]
	if end != "" {
		ns += end
	}
	return ns, nil
}

func sparse(s string) (term []string, err error) {
	re := regexp.MustCompile(`[0-9\.]+|'[^']+'|[+\-*/\\(\)\|&=><!]`)
	matches := re.FindAllStringSubmatchIndex(s, -1)
	if matches == nil {
		return term, fmt.Errorf("parse fail")
	}
	prev := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		if start != prev {
			return term, fmt.Errorf("unknow op: %v", s[prev:start])
		}
		name := s[start:end]
		term = append(term, name)
		prev = end
	}
	return term, nil
}
func isOperate(s string) bool {
	re := regexp.MustCompile(`[+\-*/\\(\)\|&=><!]`)
	return re.MatchString(s)
}
func pparse(pn []string) []string {
	var vs, vp stack
	for _, t := range pn {
		if isOperate(t) {
			switch {
			case t == "(":
				vp.Put(t)
			case vp.Top() == "":
				vp.Put(t)
			case vp.Top() == "(":
				vp.Put(t)
			case t != ")" && comp(t, vp.Top()):
				vp.Put(t)
			case t == ")":
				for {
					if pop := vp.Pop(); pop == "(" {
						break
					} else {
						vs = append(vs, pop)
					}
				}
			default:
				for {
					top := vp.Top()
					if top != "" && comp(top, t) {
						vs = append(vs, vp.Pop())
					} else {
						vp = append(vp, t)
						break
					}
				}
			}
		} else {
			vs.Put(t)
		}
	}
	for {
		if pop := vp.Pop(); pop != "" {
			vs = append(vs, pop)
		} else {
			break
		}
	}
	return []string(vs)
}

func rparse(ss []string) (ret bool, err error) {
	var c stack
	for _, s := range ss {
		if isOperate(s) {
			b, a := c.Pop(), c.Pop()
			t := gettype(a, b)
			switch t {
			case "int":
				A, err := strconv.ParseInt(a, 10, 64)
				if err != nil {
					return ret, err
				}
				B, err := strconv.ParseInt(b, 10, 64)
				if err != nil {
					return ret, err
				}
				var tmp bool
				switch s {
				case `>`:
					tmp = A > B
				case `<`:
					tmp = A < B
				case `=`:
					tmp = A == B
				case `!`:
					tmp = A != B
				}
				c.Put(fmt.Sprintf("%v", tmp))
			case "bool":
				A, err := strconv.ParseBool(a)
				if err != nil {
					return ret, err
				}
				B, err := strconv.ParseBool(b)
				if err != nil {
					return ret, err
				}
				var tmp bool
				switch s {
				case `|`:
					tmp = A || B
				case `&`:
					tmp = A && B
				}
				c.Put(fmt.Sprintf("%v", tmp))
			case "string":
				var tmp bool
				switch s {
				case `=`:
					tmp = a == b
				case `!`:
					tmp = a != b
				}
				c.Put(fmt.Sprintf("%v", tmp))
			}
		} else {
			c.Put(s)
		}
	}
	if strings.Index(c.Pop(), "true") != -1 {
		return true, nil
	}
	return false, nil
}
func Caculate(dic map[string]string, str string) (ret bool, err error) {
	str, err = fparse(filter(str), dic)
	if err != nil {
		return ret, err
	}
	pn, err := sparse(str)
	if err != nil {
		return ret, err
	}
	rpn := pparse(pn)
	ret, err = rparse(rpn)
	return ret, err
}
