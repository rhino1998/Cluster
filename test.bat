
for /L %%A in (1,1,100) do (
    PING 1.1.1.1 -n 1 -w 500 >NUL
    start curl -X POST -m 10000 -d @testtask.json http://108.56.251.125:2002/api/task
)