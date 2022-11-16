-- example HTTP POST script which demonstrates setting the
-- HTTP method, body, and adding a header

math.randomseed(os.time())
math.random(); math.random(); math.random(); --uses this to make random number generator in the request function below functioning well

request = function()
    num = math.random(1,10000000)
    num = math.fmod(num, 200)
    wrk.method = "POST"
    wrk.body   = "{\"account\": \"user" .. num .. "\",\"password\": \"202cb962ac59075b964b07152d234b70\"}"
    wrk.headers["Content-Type"] = "text/plain; charset=utf-8"
    return wrk.format(wrk.method, nil, wrk.headers, wrk.body)
end