local redis = require 'redis_module'
local password = require 'password'
local json = require 'json'

function register_subscriptions(username, subscriptions)
  for i = 1, #subscriptions do
    err, ret = redis.execute("set", 2, username.."-"..subscriptions[i], 2)
    if err ~= nil then
      return err, 1
    end
  end
  return nil, 0
end

function run_plugin(topic, payload)
  local json_message = json.decode(payload)
  --message: {"payload": {"message": "register", "username": "username", "password": "userpass", subscriptions: ["subs1", "subs2"]}}
  print(payload, json_message)
  local username = json_message["payload"]["username"]
  local pass = json_message["payload"]["password"]
  local subscriptions = json_message["payload"]["subscriptions"]
  err, ret = password.generate_hash(pass)
  if err ~= nil then
    return err, 1
  end
  err, ret = redis.execute("set", 2, username, ret)
  if err ~= nil then
    return err, 1
  end
  if subscriptions ~= nil then
    err, ret = register_subscriptions(username, subscriptions)
    if err ~= nil then
      return err, 1
    end
  end
  return nil, 0
end
