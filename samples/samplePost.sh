curl -H 'Content-Type:application/json' -d \
"[{\"name\":\"jobpost\",
 \"second\":\"*\",
 \"hour\":\"*\",
 \"minute\":\"*\",
 \"day\":\"*\",
 \"month\":\"*\",
 \"command\":\"echo post\"}]" http://localhost:3231/add
