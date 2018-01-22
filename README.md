# Gronit
#### A Cron monitor written in Go.

##### Features:
 * Update and check job status remotely
 * View statistics of jobs
 * Bare bones alternative to services like Cronitor

##### Installing:
```
$ go get github.com/zricethezav/gronit
```

##### Usage: 

```
$ ./gronit
```
This will start a server on default port `3231`. Change the port with `-p` or `--port` option.

Generate a job tracker (with a gronit server running):
```
$ curl 127.0.0.1:3231/create
{"id":"f7a324"}
```
`/create` generates a job tracking token. Now we can wrap any command with gronit.

```
# sample crontab
*/2 * * * * curl -s 127.0.0.1:3231/run/f7a324 && sleep `echo $((1 + RANDOM \% 5))` && curl -s 127.0.0.1:3231/complete/f7a324

```
#### API 
GET
* `/create` generate new job id
* `/run/{id}` update jobs's status to be 'running'
* `/complete/{id}` update jobs's status to be 'complete'
* `/status/{id}` status of job
  * sample response:```{"status":"complete","time":"2018-01-10T13:36:05.197347-06:00"}```
* `/summary/{id}` job statistics for a job
  * sample response: ```{"status_count":24,"run_count":12,"completion_count":12,"average_time_to_completion":2200,"created_at":"2018-01-10T13:08:11.675697-06:00"}```
* `/history/{id}` full list of status updates for a job
  * sample response:
  ```[{"status":"running","time":"2018-01-10T13:10:00.281757-06:00"},{"status":"complete","time":"2018-01-10T13:10:01.408825-06:00"}]```

