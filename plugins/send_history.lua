local pm = require 'persistence_module'
local json = require 'json'

function run_plugin(topic, payload)
  local jsonMessage = json.decode(payload)
  --message: {"payload": {"from": "someone", "message": "history", "limit": 100, "start":0}}
  print(jsonMessage["payload"]["from"])
  pm.query_messages("chat/tanks/clans/17feaa05-cc54-4d9d-947b-096a69d19896", 50, 0)
  return nil, 0
end
