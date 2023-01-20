# kol-tradeable-items-jsoner

> Processes tradeable item data to JSON format.

This tool uses [Colly](https://github.com/gocolly/colly).

###  Example
```shell
go build .
.\kol-tradeable-items-jsoner.exe (-u)
```

## TODO
Possibly use [goquery](https://github.com/PuerkitoBio/goquery) instead.

## The main problem:
If you wanted the item name, item number, and tradeability, you would have to visit each item's wiki url to get that data.

This means crawling and scraping through ~12000 items URLs. (The number of items constantly increase.)

The only issue with this is that the site is already pretty laggy. Adding more traffic would be bad manners to the site and the players.

## The solution:
There is a dropdown for tradeable items on a different site. 

I can just get the HTML file and then parse the item name and item number from it.

The tradability of the item is already guaranteed by the fact that it has trade data for it.

After doing some data processing, the item name and item number are encoded into a JSON file.

## Reasoning:
I didn't want to unnecessarily get a new HTML file every time the program is run.

I chose to scrape data from a local HTML file instead of the website because the data is static.

## Note:
First time users need to run with "-u" flag to generate the HTML file from the website locally.

After that users can just run without the "-u" flag since HTML file exists locally.

Users can rerun with "-u" flag if they need to re-update the HTML file (if new tradeable items are created).

Users can change filenames by modifying the const variables.
