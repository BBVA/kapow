Remote tcpdump sniffer with source filtering
============================================

1. Add any filter you want to the `tcpdump` command inside `tcpdump.pow` to filter
   any traffic you don't want to be sniffed!
2. For the sake of simplicity, run `sudo -E kapow server tcpdump.pow`. In a
   production environment, `tcpdump` should be run with the appropiate permissions,
   but kapow can (and should) run as an unprivileged user.
3. In your local machine run:
   ```bash
   curl http://localhost:8080/sniff/<network-interface> | sudo -E wireshark -k -i -
   ```
   Again, for the sake of simplicity, `Wireshark` is running as root. If you don't want
   to run it this way, follow this guide:
   https://gist.github.com/MinaMikhailcom/0825906230cbbe478faf4d08abe9d11a
4. Profit!
