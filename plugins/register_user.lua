local redis = require 'redis_module'
local password = require 'password'

function run_plugin(topic, payload)
  err, ret = password.generate_hash("password")
  err, ret = redis.execute("set", 2, "teste", ret)
  if err ~= nil then
    return err, 1
  end
  return nil, 0
end
