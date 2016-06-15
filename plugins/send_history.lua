local pm = require 'persistence_module'
local json = require 'json'

function run_plugin(topic, payload)
  local jsonMessage = json.decode(payload)
  --message: {"payload": {"from": "someone", "message": "history", "limit": 100, "start":0}}
  err, payloads = pm.query_messages(jsonMessage["payload"]["topic"], jsonMessage["payload"]["limit"], jsonMessage["payload"]["start"])
  if err ~= nil then
    return err, 1
  end
  local payloadTable = {}
  for i = 1, #payloads do
    payloadTable[i] = {from = payloads[i]["from"], message = payloads[i]["message"], timestamp = payloads[i]["timestamp"]..""} 
  end
  return nil, 0
end
