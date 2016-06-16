local redis = require 'redis_module'
local password = require 'password'
local json = require 'json'

function run_plugin(topic, payload)
  local json_message = json.decode(payload)
  --message: {"payload": {"message": "register", "username": "username",
  --          "password": "userpass"}}
  local username = json_message["payload"]["username"]
  local pass = json_message["payload"]["password"]
  err, ret = password.generate_hash(pass)
  if err ~= nil then
    return err, 1
  end
  err, ret = redis.execute("set", 2, username, ret)
  if err ~= nil then
    return err, 1
  end
  return nil, 0
end
