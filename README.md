# Gronit
#### A Cron monitor and interface written in Go.

##### Features:
 * Remotely update crontabs
 * View active crontabs in a browser
 * Add cron jobs to crontab locally via json or yaml
 * View run logs of jobs in a browser 

##### Server:
```
$ ./gronit start
```
* view your crontab: ```curl localhost:3231/list``` 
* add jobs to crontab via curl: 
	```
    curl -H 'Content-Type:application/json' -d \
	"[{\"name\":\"example_post\",
 	\"second\":\"*\",
 	\"hour\":\"*\",
 	\"minute\":\"*\",
 	\"day\":\"*\",
 	\"month\":\"*\",
 	\"command\":\"echo post\"}]" http://localhost:3231/add
    ```  
 	* `/add` accepts POST requests with a JSON array containing job objects. See example request [here](https://github.com/zricethezav/gronit/blob/master/samples/sample.json)

##### Local Usage:
You can use gronit to add cron jobs to your crontab in more readable forms.

<b>Yaml</b>: ```./gronit --loadyaml {yamlfile}```

<b>JSON</b>: ```./gronit --loadjson {jsonfile}```

	   
### TODO: 
* Assign keys to individual jobs so they can be remotely opperated on
* Include hook like cronitor has to update status of the job
* Finish handlers: remove, logs, update

