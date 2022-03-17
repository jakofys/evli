# evli

Evli is a event listener implementation for go application.

## Usage

### Event

An `Event` is a simple structure with public scope field as:

```go
type RandomEvent struct {
    Name string
    Age int
    place string // will be ignore when transmitting payload
}
```

### Listener

For create a `Listener` simple implement `Listener` interface:

```go
func ListenerExample(*evli.Source){}
```

Declare it using `Register` to bind multiple `Listener` to a single `Event`:

```go
_ = evli.Register(RandomeEvent{},ListenerExample, SeconListenerExemple)
```

### Source


Make sure field data as same name as event payload field.
You can also use tag `field` to bind field data structure to event payload field:

```go
type RandomData struct {
    UserName string `field:"Name"`
    Age int
}
```

For retrieving data, use a pointer structure to rehydrate from the `Event` payload using the `Source` object:

```go
func ListenerExample(s *evli.Source){
    data := &RandomData{}
    _ = s.Listen(data)
}
```