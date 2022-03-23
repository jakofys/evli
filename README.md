# evli

Evli is a event Subscriber implementation for go application.

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

### Subscriber

For create a `Subscriber` simple implement `Subscriber` interface:

```go
func SubscriberExample(*evli.Source){}
```

Declare it using `Subscribe` to bind multiple `Subscriber` to a single `Event`:

```go
_ = evli.Subscribe(RandomeEvent{},SubscriberExample, SeconSubscriberExemple)
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
func SubscriberExample(s *evli.Source){
    data := &RandomData{}
    _ = s.Payload(data)
}
```