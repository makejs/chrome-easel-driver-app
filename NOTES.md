### Difficulties

First thing we need is an http server. I didn't look to hard, but it doesn't seem
super popular to write an http server *inside* the browser. Node's own libraries
are not meant to run there, and even aside from the native extensions they aren't
written in a very composable way for this purpose.

Since Go's runtime is 100% Go and the language lends itself to this type of re-use
it made an interesting solution. `gopherjs` will compile Go code into javascript,
and the *Hello, World!* http server running in chrome took no time at all.

Another pain point is socket.io. While it works well, it is definately not portable. The
server component requires native extensions, and a whole handful of node libraries.
This isn't necessarily bad, but it meant using socket.io wasn't possible. It also seems
to be all but impossible to find documentation on the actual on-wire protocol.

Luckily someone wrote a compatible Go library. Most of the code around it is handling
types/serialization to-from javascript.

And the last pain point, javascript is a dynamic language. Trying to reason about
code you know little about, digging through reference by reference throughout the
code, for hints to what is going on can be painful. The code isn't that messy, but
there's no documentation as to what properties the *config* object supports. Or
what format is *timestamp* anyway... (I'm really hoping it's just a string)
