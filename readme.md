REST API

000-user-service

Invoke-RestMethod -Uri "http://localhost:5432/register" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"username":"user1","password":"1234"}'
Invoke-RestMethod -Uri "http://localhost:5432/login" -Method Post -Headers @{"Content-Type"="application/json"} -Body '{"username":"user1","password":"1234"}'
Invoke-RestMethod -Uri "http://localhost:5432/check-token" -Method Get -Headers @{"Authorization"="Mk7UV0NjJrCU7/vT57OUVfdmg827YTtrZqilOlhdbqc="}



001-log-service

Invoke-RestMethod -Uri "http://localhost:5434/log" -Method Post -Body (@{message = "This is a log message from PowerShell"} | ConvertTo-Json) -ContentType "application/json"