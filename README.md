This project accepts events and pushes them into redis. If pushing to redis fails,
the api returns error

Invoke endpoint "/event" using POST and pass below struct in json

{
   "userId":"value",
   "payload":"value"
}