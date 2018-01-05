curl -H 'Content-Type:application/json' -d \
"[{\"name\":\"jobpost\",
 \"second\":\"*\",
 \"hour\":\"*\",
 \"minute\":\"*\",
 \"day\":\"*\",
 \"month\":\"*\",
 \"command\":\"echo post\",
 \"monitor\":true}]" http://localhost:3231/add

curl -H 'Content-Type:application/json' -d \
"[{\"name\":\"nomonitor\",
 \"second\":\"*\",
 \"hour\":\"*\",
 \"minute\":\"*\",
 \"day\":\"*\",
 \"month\":\"*\",
 \"monitor\":true,
 \"command\":\"echo no monitor\"}]" http://localhost:3231/add
