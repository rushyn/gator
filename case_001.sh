go build -o out
./out reset
./out register TesetUser
./out login TesetUser
./out addfeed "TechCrunch" "https://techcrunch.com/feed"
./out addfeed "Hacker News" "https://news.ycombinator.com/rss"
./out addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
./out browse 5
./out browse