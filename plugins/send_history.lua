local pm = require 'persistence_module'
local json = require 'json'

function run_plugin(topic, payload)
  local jsonMessage = json.decode(payload)
  --message: {"payload": {"from": "someone", "message": "history", "limit": 100, "start":0}}
  err, messages = pm.query_messages(jsonMessage["payload"]["topic"], jsonMessage["payload"]["limit"], jsonMessage["payload"]["start"])
  return nil, 0
end
