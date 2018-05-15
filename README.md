# token-manager
Web-server

Developed on Andrey's laptop)

curl  -X POST http://localhost:2803/reg-email -d'{"email":"natakr@i.ua", "nick":"marinakr"}'
SUCCESS:
{  
   "code":0,
   "mess":""
}

ERRORS:
{  
   "code":1,
   "mess":"Invalid email"
}
{  
   "code":2,
   "mess":"Invalid nickname"
}
//both email and nickname invalid
{  
   "code":3,
   "mess":"Invalid data"
}
{  
   "code":4,
   "mess":"Email  already in use"
}

{  
   "code":5,
   "mess":"Nick already in use"
}


curl  -X POST http://localhost:2803/confirm-email -d'{"email":"natakr@i.ua", "code":8452}'

SUCCESS:
{  
   "code":0,
   "mess":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTg1MjIwMTIsImlvdCI6MTUyNjM4MTIxMiwidXNlciI6Im1hcmluYWtyIn0.TuyKuzd3u8fRZk1yx1Gg0lf7Sbf38Ze_rOOYLUiWc-k"
}

ERRORS:
{  
   "code":7,
   "mess":"Confirmation code is not match"
}

{  
   "code":6,
   "mess":"Confirmation time expired / Email not found"
}


