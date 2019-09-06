Remote tcpdump sniffer with source filtering
============================================

1. Add any filter you want to the tcpdump command inside `tcpdump.pow`.
2. For the sake of simplicity run `sudo kapow server tcpdump.pow`. In a
   production environment tcpdump should be run with the appropiate permissions
   but kapow can (and should) run as an unprivilieged user.
3. In your local machine run `curl http://localhost:8080/sniff/<network-interface> | sudo
   wireshark -k -i -` if you don't want to run Wireshark as root follow this
   guide: https://gist.github.com/MinaMikhailcom/0825906230cbbe478faf4d08abe9d11a
4. Profit!
