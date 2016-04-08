# Go Fly ~

On-the-fly image resizing using Go and VIPS.

It is deployable to Heroku :D

# Why VIPS ?

VIPS is the fastest image resizing library. Check the benchmark here:

https://github.com/fawick/speedtest-resize#benchmark

# Deploy

The current implementation is Heroku-ready(TM) with custom buildpacks.

If you want to deploy to your own infrastructure, please have libvips installed properly.

See: https://github.com/DAddYE/vips

# Demo

https://go-image-fly.herokuapp.com/resize?src=https://pbs.twimg.com/media/BnSJTo9CcAIzy0B.jpg&size=320x240&mode=crop

# Enjoy

[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)
