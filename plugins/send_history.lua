local pm = require 'persistence_module'

function run_plugin(topic, payload)
  pm.query_messages("chat/tanks/clans/17feaa05-cc54-4d9d-947b-096a69d19896", 50, 0)
  return nil, 0
end
