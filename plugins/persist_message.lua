local pm = require 'persistence_module'

function run_plugin(topic, payload)
  err, ret = pm.index_message(topic, payload)
  return err, ret
end
