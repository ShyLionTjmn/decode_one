package main

import (
  "regexp"
  "errors"
  _ "strings"
  _ "strconv"
  _ "time"
  "fmt"
  "os"
  "github.com/fatih/color"
  "github.com/gomodule/redigo/redis"
  "encoding/json"
  . "github.com/ShyLionTjmn/mapaux"
  . "github.com/ShyLionTjmn/decode_dev"
)

const REDIS_SOCKET = "/tmp/redis.sock"
const REDIS_DB = "0"

const USAGE = "Usage: %s [-M] [-A] [-J] IP\n\t-M\t- do not process MACs\n\t-A\t- do not process ARP\n\t-J\t- do not output JSON\n"

var red_db string = REDIS_DB

func main() {
  var err error

  ip_regex := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)

  opt_j := true //print out resulting json
  opt_m := true //process MAC addresses
  opt_a := true //process ARP info

  defer func() {
    if err != nil {
      color.New(color.FgRed).Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }
  } ()

  args := os.Args[1:]

  for len(args) > 0 && len(args[0]) > 0 && args[0][0] == '-' {
    switch args[0] {
    case "-M":
      opt_m = false
    case "-J":
      opt_j = false
    case "-A":
      opt_a = false
    default:
      fmt.Printf(USAGE, os.Args[0])
      return
    }
    args = args[1:]
  }

  if len(args) != 1 {
    fmt.Printf(USAGE, os.Args[0])
    return
  }

  dev_ip := args[0]

  if !ip_regex.MatchString(dev_ip) {
    err = errors.New("Bad IP: "+dev_ip)
    return
  }


  var red redis.Conn

  red, err = redis.Dial("unix", REDIS_SOCKET)
  if err != nil { return }

  defer func() {
    red.Close()
  } ()

  _, err = red.Do("SELECT", red_db)
  if err != nil { return }

  var raw M

  raw, err = GetRawRed(red, dev_ip)
  if err != nil { return }

  var dev = Dev{Opt_m: opt_m, Opt_a: opt_a, Dev_ip: dev_ip}

  err = dev.Decode(raw)
  if err != nil { return }

  if opt_j {
    var j []byte
//  j, err = json.MarshalIndent(raw, "", "  ")
//  if err != nil { return }
//  fmt.Println(string(j))

    j, err = json.MarshalIndent(dev.Dev, "", "  ")
    if err != nil { return }
    fmt.Println(string(j))
  }
}
