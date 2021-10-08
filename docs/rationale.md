# Rationale

How do you make a static website? You use one of the many static site
generators to produce a bunch of html/css/js files. You put the files
behind some web server, Apache and Nginx spring to mind, but it can be
anything. And you have yourself a website. Easy to set up, easy to maintain
But...

What if you wanted to make the website a little bit more private? So that
only a few people could see it? This is where this little utility comes in!

All you need is a list of users with their passwords, point it to
a directory with your static website and you're done! It will start
serving the website and every request will have to be authenticated
against the list of users.

The next step would be to put everything behind a proxy with TLS
configured. Or maybe put it inside a Docker container? Who knows? Sky is
the limit!

