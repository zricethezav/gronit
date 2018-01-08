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
* `/create` generate new job id
* `/run/{id}` update jobs's status to be 'running'
* `/complete/{id}` update jobs's status to be 'complete'
* `/status/{id}` status of job
* `/summary/{id}` job statistics for a job
* `/history/{id}` full list of status updates for a job

