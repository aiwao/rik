# rik
 golang net/http helper
 
# Usage
## Do request
```go
res, err := rik.Get("URL").
    Client(&http.Client{}).
    Do()
```
## Read response
```go
res, err := rik.Get("URL").
    Client(&http.Client{}).
    Do()
if err != nil ...

//No error check
res := rik.Get
    ...
    MustDo()

//Read res.Body
//return []byte
b, err := rik.ReadByte(res)

//return string
s, err := rik.ReadString(res)

//No error check
b := rik.MustReadByte(res)
s := rik.MustString(res)
```
## Do and read response
```go
b, res, err := rik.Get("URL").
    Client(&http.Client{}).
    DoReadByte()
    //DoReadString()
if err != nil ...

//No error check
b, res := rik.Get
    ...
    MustDoReadByte()
    //MustDoReadString()
```
## With context
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
rik.Get("URL").
    Context(ctx)
```
## Request data
rik automatically sets the content-type
```go
//JSON
rik.Post("URL").
    JSON(NewJSON().
        Set("Key", "Value").
        Set("Keys", []any{"Value1"}).
        Add("Keys", "Value2").
        Build(),
    )
)

//Query
rik.Post(requestTester).
    Query("Key", "Value1").
    Query("Key", "Value2").
    Query("Key2", "Key2Value1", "Key2Value2").
    Queries(url.Values{
        "Key3": {"Key3Value1", "Key3Value2"},
        "Key4": {"Key4Value1", "Key4Value2"},
    })

//File + Header
file, err := os.Open("hello.txt")
if err != nil ...
defer file.Close()
rik.Post(requestTester).
    Header("Key", "Value1").
    Header("Key", "Value2").
    Header("Key2", "Value1", "Value2").
    Headers(http.Header{
        "Key3": {"Key3Value1", "Key3Value2"},
        "Key4": {"Key4Value1", "Key4Value2"},
    }).
    File(file)

//Form
rik.Put(requestTester).
    Form("Key", "Value1").
    Form("Key", "Value2").
    Form("Key2", "Key2Value1", "Key2Value2").
    Forms(url.Values{
        "Key3": {"Key3Value1", "Key3Value2"},
        "Key4": {"Key4Value1", "Key4Value2"},
    })

//Multipart
file, err := os.Open("hello.txt")
if err != nil ...
defer file.Close()
rik.Post(requestTester).
    Multipart(NewMultipart().
        File("hello.txt", file).
        Field("content", strings.NewReader("data")).
        Boundary("END_OF_PART").
        //No error check
        MustBuild(),
    )
```
## Set client
```go
//In builder
rik.Get("URL").
    Client(&http.Client{})
//In do request
rik.Get("URL").
    DoClient(&http.Client{})
    //DoReadByte(&http.Client{})
    //DoReadString(&http.Client{})
    //MustDoReadByte(&http.Client{})
    //MustDoReadString(&http.Client{})
```
## Set default client
```go
rik.DefaultClient = &http.Client{}
```