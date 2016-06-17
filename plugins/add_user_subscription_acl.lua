local redis = require 'redis_module'
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
  --message: {"payload": {"message": "aclpermitsubscription", "username": "username", subscriptions: ["subs1", "subs2"]}}
  local username = json_message["payload"]["username"]
  local subscriptions = json_message["payload"]["subscriptions"]
  if subscriptions ~= nil then
    err, ret = register_subscriptions(username, subscriptions)
    if err ~= nil then
      return err, 1
    end
  end
  return nil, 0
end
