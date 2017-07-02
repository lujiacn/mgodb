# Include MgoDB initiation and MgoSession initiation for controller

## Using method

in app/init.go
```
revel.OnAppStart(mgodb.MgoDBInit)
```

in controller/init.go
```
mgodb.MgoControllerInit()
```


