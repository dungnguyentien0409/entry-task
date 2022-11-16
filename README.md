# ENTRY TASK
Owner: chris.nguyen@shopee.com

Mentor: tung.do@shopee.com

# STEPS TO RUN THE SYSTEM

1. RUN REDIS SERVER

    redis-server
2. RUN NGINX SERVER

   sudo brew services restart nginx
3. RUN TCP SERVER 

    entrytask/tcp-server: go main.go

    or run the binary file
4. RUN HTTP SERVER

    entrytask/http-server: go main.go

    or run the binary file
5. START TO TEST USE POSTMAN OR WRK

# PERFORMANCE TEST
    Run Performance Test Usint WRK

1. POINT TO NGINX

   wrk -t 12 -c 300 -d 30s http://entrytask.com/
2. POINT TO HTTP SERVER

    wrk -t 12 -c 300 -d 30s http://localhost:49/ping

    wrk -t 12 -c 300 -d 30s http://entrytask.com/api/ping
3. POINT TO HTTP SERVER THEN CALL TO TCP SERVER

    wrk -t 12 -c 300 -d 30s -s bench.lua http://entrytask.com/api/login

    wrk -t 12 -c 300 -d 30s -s bench.lua http://localhost:80/api/login

    wrk -t 12 -c 300 -d 30s -s bench.lua http://localhost:49/login

    wrk -t 12 -c 300 -d 30s -s bench.lua http://entrytask.com/api/login
