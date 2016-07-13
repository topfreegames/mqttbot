local pm = require 'persistence_module'
local client = require 'mqttclient_module'
local json = require 'json'

function run_plugin(topic, payload)
  local json_message = json.decode(payload)
  --message: {"payload": {"from": "someone", "topic": "xxx", "message": "history", "limit": 100, "start":0}}
  local history_topic = json_message["payload"]["topic"]
  local user_requesting = json_message["payload"]["from"]
  err, payloads = pm.query_messages(history_topic, json_message["payload"]["limit"], json_message["payload"]["start"])
  if err ~= nil then
    return err, 1
  end
  local payload_table = {}
  for i = 1, #payloads do
    payload_table[i] = {id = payloads[i]["id"], timestamp = payloads[i]["timestamp"].."", payload = payloads[i]["payload"]} 
  end
  client.send_message(history_topic.."/history/"..user_requesting, 2, false, json.encode(payload_table))
  return nil, 0
end
