box.cfg{
    listen = 3301
}

box.schema.user.grant('guest', 'read,write,execute', 'universe')


s = box.schema.space.create('polls', {if_not_exists = true})

s:format({
      {name = 'id', type = 'string'},
      {name = 'owner_id', type = 'string'},
      {name = 'question', type = 'string'},
      {name = 'options', type = 'array'},
      {name = 'is_active', type = 'boolean'}
})

s:create_index('primary', {parts = {'id'}, if_not_exists = true})
s:create_index('owner', {parts = {'owner_id'}, if_not_exists = true})

v = box.schema.space.create('votes', {if_not_exists = true})
v:format({
    {name = 'poll_id', type = 'string'},
    {name = 'user_id', type = 'string'},
    {name = 'choice', type = 'unsigned'}
})
v:create_index('primary', {
    parts = {'poll_id', 'user_id'}, if_not_exists = true
})
