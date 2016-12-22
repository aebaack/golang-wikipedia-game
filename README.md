#What is the Wikipedia Game
The Wikipedia Game begins with a player choosing a starting Wikipedia page and an ending Wikipedia page.
From the starting page, the player then has to traverse internal links on the articles in order 
to arrive at the ending page. A player is scored based on the number of links that they had to click on
in order to get from start to end.

#Golang Wikipedia Game
This is a program that will find a path from the starting Wikipedia page to the ending page. The pages can be specified
as command line arguments like: 

`go run scraper.go https://en.wikipedia.org/wiki/Gengar https://en.wikipedia.org/wiki/Banana`

The program will find one possible path from the starting page to the ending page, so it will reveal a new path each time
that it is run.

##Work in Progress
This program is definitely still buggy, especially if it takes a long path. Try to pick articles that seem somewhat closely
linked for best results. If the path gets too long, the program starts to return a false route from the starting page to 
the ending page.
