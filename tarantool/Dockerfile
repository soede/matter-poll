FROM tarantool/tarantool:3.1.0
WORKDIR /opt/tarantool

COPY init.lua /opt/tarantool/init.lua

EXPOSE 3301

CMD ["tarantool", "/opt/tarantool/init.lua"]
