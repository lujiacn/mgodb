# Include MgoDB initiation and MgoSession initiation for controller

## Import 

- gopkg.in/mgo.v2
```go
import gopkg.in/lujiacn/mgodb.v0
```
- github.com/globalsign/mgo
```go
import gopkg.in/lujiacn/mgodb.v1
```


## Using method
in app/init.go
```
revel.OnAppStart(mgodb.MgoDBInit)
```

in controller/init.go
```
mgodb.MgoControllerInit()
```


