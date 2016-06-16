local redis = require 'redis_module'

function run_plugin(topic, payload)
  err, ret = redis.execute("set", 2, "teste", "teste")
  if err ~= nil then
    return err, 1
  end
  return nil, 1
end
