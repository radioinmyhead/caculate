# caculate
caculate

## install 
```shell
go get github.com/radioinmyhead/caculate
```

## usage
```golang
dic := map[string]string{"a": "10", "b": "100", "c": "ok"}
exp := `${a} > 1 && ( ${a}>50||${b}>10)&&  "${c}" != "ok"` 
ret, err := caculate.Caculate(dic,exp)
if err != nil {
  panic(err)
}
fmt.Println(ret) // false
```
  
