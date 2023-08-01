[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

# forex
A simple tool to quickly scrape and format current foreign exchange rates from Google. Provides an optional REST API to more easily integrate with tools such as the Microsoft Excel `QueryTables` function.

Usage: `forex [options...]`

Argument                  | Description
--------------------------|-----------------------------------------------------------------------------------------------------
 `[-base] <base>`         | Base currency
 `[-quote] <quote>`       | Quote currency
 `-decimal <separator>`   | Decimal separator
 `-thousands <separator>` | Thousands seperator
 `-rest <address:port>`   | Start REST API on given socket

**This tool scrapes from a standard Google web search, and as such will get flagged as "unusual traffic" if it is used cyclically in rapid succession. The intent of this tool is for personal use only, to simply to aid with certain one-off daily tasks, and should not be used in place of non-free Google API access.**

If `decimal` and `thousands` separators are not defined, a dot (".") will be used for decimal separators and no thousands separator will be used.

If `rest` API is defined, all other arguments are ignored.

# More About ScriptTiger

For more ScriptTiger scripts and goodies, check out ScriptTiger's GitHub Pages website:  
https://scripttiger.github.io/

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=MZ4FH4G5XHGZ4)
