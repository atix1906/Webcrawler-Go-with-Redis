package redis

import(
  "fmt"
  "errors"
  "github.com/mediocregopher/radix.v2/pool"
  // Import the Radix.v2 redis package (we need access to its Nil type).
  "github.com/mediocregopher/radix.v2/redis"
  "log"
)

var db *pool.Pool

func init(){
  var err error
  db,err = pool.New("tcp", "localhost:6379", 10)
  if err != nil{
    log.Panic(err)
  }
}

func saveInRedisDB(link string){
  
}
