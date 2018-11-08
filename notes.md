
### Moon API Info 

[API Source - US Naval Observatory](http://aa.usno.navy.mil/data/docs/api.php#rstt)

[ST = sidereal time](http://aa.usno.navy.mil/data/docs/api.php#rstt)

```
{
"error":false,
"apiversion":"2.1.0",
"year":2018,
"month":11,
"day":7,
"dayofweek":"Wednesday",
"datechanged":false,
"state":"NY",
"city":"New York",
"lon":-73.92,
"lat":40.73,
"county":"",
"tz":-5,
"isdst":"no",

"sundata":[
            {"phen":"BC", "time":"6:04 a.m. ST"},
            {"phen":"R", "time":"6:33 a.m. ST"},
            {"phen":"U", "time":"11:39 a.m. ST"},
            {"phen":"S", "time":"4:45 p.m. ST"},
            {"phen":"EC", "time":"5:14 p.m. ST"}],

"moondata":[
            {"phen":"R", "time":"6:15 a.m. ST"},
            {"phen":"U", "time":"11:47 a.m. ST"},
            {"phen":"S", "time":"5:12 p.m. ST"}],

"closestphase":{"phase":"New Moon","date":"November 7, 2018","time":"11:02 a.m. ST"}
}
```


Lat and Lon 

Lunar Phenomena:
- Rise (R)
- Upper Transit (U)
- Set (S)

Closest phase 

### Proto 

```

protoc -I proto/ proto/phases.proto --go_out=plugins=grpc:proto

```

. venv/bin/activate
FLASK_APP=app.py
flask run 


python -m grpc_tools.protoc -I../moonphases/proto --python_out=. --grpc_python_out=. ../moonphases/proto/phases.proto



### Docker run 

docker run -p 8001:8001 meganokeefe/moonphases:latest --city "Malvern, PA"

docker run -p 2345:2345  meganokeefe/moonfacts:latest

docker run -p 5000:5000 meganokeefe/moonsite:latest phases="127.0.0.1:8001"
facts="http://127.0.0.1:2345/" 