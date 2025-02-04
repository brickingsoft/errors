# errors
增强错误。

使用于错误定位、序列化、传输或存储等。

## 安装
```shell
go get -u github.com/brickingsoft/errors
```

## 使用
使用`github.com/brickingsoft/errors`代替`errors`。
```go
e1 := errors.Define("e1")
e2 := errors.New("e2", errors.WithWrap(e1))
e3 := errors.New("e3", errors.WithWrap(e2),
    errors.WithMeta("s", "s"),
    errors.WithMeta("i", 1),
    errors.WithMeta("i32", int32(32)),
    errors.WithMeta("i64", int64(-64)),
    errors.WithMeta("u", uint(1)),
    errors.WithMeta("u64", uint64(64)),
    errors.WithMeta("f32", float32(32.32)),
    errors.WithMeta("f64", 64.640),
    errors.WithMeta("b", true),
    errors.WithMeta("any", struct{}{}),
    errors.WithMeta("byte", 'b'),
    errors.WithMeta("bytes", []byte("hello world")),
    errors.WithMeta("time", time.Now()),
    errors.WithMeta("ss", []string{"a a", "b"}),
    errors.WithDescription("desc"),
    errors.WithOccur(),
)
ok := errors.Is(e3, e1)
// ...
```
标准输出：
```text
EnhancedError:
>>>>>>>>>>>>>
ERRO      = message
DESC      = desc
META      = [s: s] [i: 1] [i32: 32] [i64: -64] [u: 1] [u64: 64] [f32: 32.32] [f64: 64.64] [b: true] [any: {}] [byte: 98] [bytes: hello world] [time: 2025-02-04T22:32:36+08:00] [ss: [a a b]]
OCCU      = 2025-02-04T22:32:36+08:00
FUNC      = github.com/brickingsoft/errors_test.TestErr
SEEK      = github.com/brickingsoft/errors/error_test.go:18
<<<<<<<<<<<<<
>>>>>>>>>>>>>
ERRO      = wrapped
<<<<<<<<<<<<<
```
Json:
```json
{
    "Message": "e3",
    "Description": "desc",
    "Stacktrace": {
        "Fn": "github.com/brickingsoft/errors_test.TestJson",
        "File": "github.com/brickingsoft/errors/error_test.go",
        "Line": 46
    },
    "Occur": "2025-02-04T22:34:36.2413536+08:00",
    "Meta": [
        {
            "Key": "s",
            "Value": "s"
        }
    ],
    "Wrapped": {
        "Message": "e2",
        "Description": "",
        "Stacktrace": {
            "Fn": "github.com/brickingsoft/errors_test.TestJson",
            "File": "github.com/brickingsoft/errors/error_test.go",
            "Line": 44
        },
        "Occur": "0001-01-01T00:00:00Z",
        "Meta": null,
        "Wrapped": {
            "Message": "e1",
            "Description": "",
            "Stacktrace": {
                "Fn": "",
                "File": "",
                "Line": 0
            },
            "Occur": "0001-01-01T00:00:00Z",
            "Meta": null,
            "Wrapped": null
        }
    }
}
```
## 函数说明

`errors.New`  创建一个增强错误。

`errors.From` 从一个错误中创建一个增强错误。

`errors.Define` 定义一个标准错误，一般用于错误判断时的 `target`，以及配合 `errors.From` 使用。

`errors.WithWrap` 在错误中添加包裹。

`errors.WithMeta` 为错误设置元数据，用于标注。

`errors.WithDescription` 为错误设置描述，一般用于 `message` 是枚举值。

`errors.WithOccur` 为错误设置发生时间。