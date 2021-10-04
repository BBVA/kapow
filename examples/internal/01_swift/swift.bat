
kapow route add -X POST /auth -c "kapow get /request/body | switftcli_1_7.exe /A"
kapow route add -X POST /receipt -c "kapow get /request/body | switftcli_1_7.exe /R"
kapow route add -X POST /letter -c "kapow get /request/body | switftcli_1_7.exe /L"