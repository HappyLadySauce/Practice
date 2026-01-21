# kitex 文档

```bash
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

```bash
kitex -module happyladysauce ./proto/math.proto
```

```bash
kitex -module happyladysauce -service math -use happyladysauce/kitex_gen ./proto/math.proto
```