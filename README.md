# rssmail

This program literally just polls an RSS feed for updates and sends the
description via `/sbin/sendmail`.

Usage:
```sh
./rssmail rss_url to_address ts_file from_address default_title
```

Example:
```sh
./rssmail https://social.treehouse.systems/@AsahiLinux.rss '~runxiyu/asahi-announce@lists.sr.ht' ts.1 me@runxiyu.org 'Asahi Linux Update'
```
